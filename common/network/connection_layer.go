package network

import (
	"GameServer/gslog"
	"context"
	"errors"
	"sync"
	"time"
)

var (
	// TimeoutWaitInterval 读写超时等待时间
	TimeoutWaitInterval = 20 * time.Millisecond

	ErrOperationCancel       = errors.New("operation cancelled")
	ErrConnectionLayerClosed = errors.New("connection layer closed")
)

type ConnBrokerFactory func(ctx context.Context) ConnectionBroker

type TcpConnectionKeeper struct {
	connID string

	readChan  chan ControlPacket
	writeChan chan ControlPacket
	stopChan  chan struct{}

	isClosed bool

	ctx       context.Context
	ctxCancel context.CancelFunc

	lock sync.RWMutex
}

func NewTcpConnectionKeeper(ctx context.Context, factory ConnBrokerFactory) ConnectionLayer {
	instance := &TcpConnectionKeeper{
		readChan:  make(chan ControlPacket),
		writeChan: make(chan ControlPacket),
		stopChan:  make(chan struct{}),
	}
	instance.ctx, instance.ctxCancel = context.WithCancel(ctx)

	instance.loop(factory)

	return instance
}

func (gs *TcpConnectionKeeper) loop(factory ConnBrokerFactory) {
	broker := factory(gs.ctx)
	if broker == nil {
		gslog.Error("[TcpConnectionKeeper] create connection broker failed...")
		return
	}
	gs.connID = broker.ConnectionID()

	go func() {
		defer func() {
			gs.ctxCancel()
			close(gs.stopChan)
			close(gs.writeChan)
			close(gs.readChan)
		}()

		for broker != nil {
			gs.startBroker(broker)
			if gs.IsClosed() {
				break
			}
			broker = factory(gs.ctx)
		}
	}()
}

func (gs *TcpConnectionKeeper) IsClosed() bool {
	gs.lock.RLock()
	defer gs.lock.RUnlock()
	return gs.isClosed
}

func (gs *TcpConnectionKeeper) startBroker(broker ConnectionBroker) {
	defer broker.Close()

	cid := broker.ConnectionID()
	if cid == "" || cid != gs.connID {
		gslog.Error("[tcpConnectionKeeper] broker connection identify error", "cid", cid)
		return
	}

	heartbeatAckChan := make(chan ControlPacket)
	defer close(heartbeatAckChan)

	ctxTask, ctxTaskCancel := context.WithCancel(gs.ctx)
	defer ctxTaskCancel()

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer ctxTaskCancel()
		gs.readLoop(ctxTask, broker, heartbeatAckChan)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer ctxTaskCancel()
		gs.writeLoop(ctxTask, broker)
	}()

	heartbeat := broker.Keepalive()
	if heartbeat > 0 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer ctxTaskCancel()
			gs.keepaliveLoop(ctxTask, heartbeat, heartbeatAckChan)
		}()
	}

	wg.Wait()
}

func (gs *TcpConnectionKeeper) keepaliveLoop(ctx context.Context, heartbeat time.Duration, heartbeatAckChan chan ControlPacket) {
	ticker := time.NewTicker(heartbeat)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			packet := NewControlPacket(Heartbeat)
			select {
			case gs.writeChan <- packet:
				ctxHb, ctxHbCancel := context.WithTimeout(ctx, heartbeat)
				select {
				case <-ctxHb.Done():
					ctxHbCancel()
					gslog.Error("[TcpConnectionKeeper] wait heartbeat ack timeout", "connID", gs.connID, "heartbeat", heartbeat)
					return
				case heartbeatAckChan <- packet:
					ctxHbCancel()
				}
			case <-ctx.Done():
				return
			}
		}
	}
}

func (gs *TcpConnectionKeeper) readLoop(ctx context.Context, broker ConnectionBroker, heartbeatAckChan chan ControlPacket) {
	for {
		packet, err := broker.ReadPacket()
		if err != nil {
			if IsNetTimeout(err) {
				time.Sleep(TimeoutWaitInterval)
				continue
			}
			break
		}
		gslog.Trace("[TcpConnectionKeeper] readLoop receiver packet", "connID", gs.connID, "packet", packet.String())

		switch packet.(type) {
		case *HeartbeatPacket:
			// send heartbeat ack
			heartbeatAck := NewControlPacket(HeartbeatAck)
			select {
			case <-ctx.Done():
				return
			case gs.writeChan <- heartbeatAck:
			}
		case *HeartbeatAckPacket:
			// dispatch to heart ack chan
			select {
			case <-ctx.Done():
				return
			case gs.readChan <- packet:
			}
		case *DisConnectPacket:
			// fin
			return
		default:
			select {
			case <-ctx.Done():
				return
			case gs.readChan <- packet:
			}
		}
	}
}

func (gs *TcpConnectionKeeper) writeLoop(ctx context.Context, broker ConnectionBroker) {
	var err error
	for {
		select {
		case <-ctx.Done():
			return
		case packet, ok := <-gs.writeChan:
			if !ok {
				// write chan close
				return
			}
			err = broker.WritePacket(packet)
			if err != nil {
				if IsNetTimeout(err) {
					time.Sleep(TimeoutWaitInterval)
					continue
				}
				return
			}
			// write success...
			gslog.Trace("[TcpConnectionKeeper] write packet success...", "connID", gs.connID, "packet", packet.String())
		}
	}
}

func (gs *TcpConnectionKeeper) ConnectionID() string {
	gs.lock.RLock()
	defer gs.lock.RUnlock()
	return gs.connID
}

func (gs *TcpConnectionKeeper) Close() error {
	// send disconnect
	ctxDisconnect, ctxDisconnectCancel := context.WithTimeout(gs.ctx, 5*time.Second)
	err := gs.WritePacket(ctxDisconnect, NewControlPacket(DisConnect))
	ctxDisconnectCancel()

	gs.ctxCancel()

	// wait to exit
	<-gs.stopChan

	return err
}

func (gs *TcpConnectionKeeper) Read() chan ControlPacket {
	return gs.readChan
}

func (gs *TcpConnectionKeeper) WritePacket(ctx context.Context, packet ControlPacket) error {
	defer func() {
		if err := recover(); err != nil {
			gslog.Critical("[TcpConnectionKeeper] WritePacket recover error", "err", err, "packet", packet.String())
		}
	}()

	select {
	case <-ctx.Done():
		return ErrOperationCancel
	case <-gs.ctx.Done():
		return ErrConnectionLayerClosed
	case gs.readChan <- packet:
	}

	return nil
}
