package network

import (
	"GameServer/gslog"
	"context"
	"sync"
	"time"
)

type TcpConnectionKeeper struct {
	connectionID string
	ctx          context.Context
	cancel       context.CancelFunc
	stopChan     chan struct{}
	readChan     chan ControlPacket
	writeChan    chan ControlPacket
	closed       bool
	lock         sync.RWMutex
}

func NewTcpConnectionKeeper(ctx context.Context, connFactory func(ctx context.Context) Connection) *TcpConnectionKeeper {
	instance := &TcpConnectionKeeper{
		stopChan:  make(chan struct{}),
		readChan:  make(chan ControlPacket),
		writeChan: make(chan ControlPacket),
		closed:    false,
	}
	instance.ctx, instance.cancel = context.WithCancel(ctx)

	instance.loop(connFactory)

	return instance
}

func (gs *TcpConnectionKeeper) loop(connFactory func(ctx context.Context) Connection) {
	conn := connFactory(gs.ctx)
	if conn == nil {
		return
	}
	gs.connectionID = conn.ConnectionID()

	go func() {
		defer func() {
			close(gs.readChan)
			close(gs.writeChan)
			close(gs.stopChan)
			gs.cancel()
		}()
		for {
			gs.start(conn)
			if gs.isClosed() {
				break
			}
			conn = connFactory(gs.ctx)
		}
	}()
}

func (gs *TcpConnectionKeeper) start(conn Connection) {
	defer conn.Close()
	connectionID := conn.ConnectionID()
	if gs.connectionID != "" && gs.connectionID != connectionID {
		gslog.Critical("[TcpConnectionKeeper] invalid connection id", "connectionID", connectionID)
		return
	}
	gs.connectionID = connectionID
	heartbeatRspChan := make(chan ControlPacket)
	defer close(heartbeatRspChan)

	taskCtx, taskCtxCancel := context.WithCancel(gs.ctx)
	defer taskCtxCancel()

	wg := sync.WaitGroup{}

	heartbeat := conn.Heartbeat()
	if heartbeat > 0 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer taskCtxCancel()
			gs.keepaliveLoop(taskCtx, heartbeat, heartbeatRspChan)
		}()
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer taskCtxCancel()
		gs.readLoop(taskCtx, conn, heartbeatRspChan)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer taskCtxCancel()
		gs.writeLoop(taskCtx, conn)
	}()

	wg.Wait()
}

func (gs *TcpConnectionKeeper) readLoop(ctx context.Context, conn Connection, heartbeatRspChan chan ControlPacket) {

}

func (gs *TcpConnectionKeeper) writeLoop(ctx context.Context, conn Connection) {

}

func (gs *TcpConnectionKeeper) keepaliveLoop(ctx context.Context, heartbeat time.Duration, heartbeatRspChan chan ControlPacket) {
	ticker := time.NewTicker(heartbeat)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			
		case <-ctx.Done():
			return
		}
	}
}

func (gs *TcpConnectionKeeper) isClosed() bool {
	gs.lock.RLock()
	defer gs.lock.RUnlock()

	return gs.closed
}

func (gs *TcpConnectionKeeper) ConnectionID() string {
	return gs.connectionID
}

func (gs *TcpConnectionKeeper) Read() chan ControlPacket {
	//TODO implement me
	panic("implement me")
}

func (gs *TcpConnectionKeeper) ReadPacket() ControlPacket {
	//TODO implement me
	panic("implement me")
}

func (gs *TcpConnectionKeeper) Write(ctx context.Context, packet ControlPacket) error {
	//TODO implement me
	panic("implement me")
}

func (gs *TcpConnectionKeeper) Close() error {
	//TODO implement me
	panic("implement me")
}
