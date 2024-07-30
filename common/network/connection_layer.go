package network

//
//import (
//	"context"
//	"encoding/binary"
//	"errors"
//	"net"
//	"sync"
//	"time"
//
//	"GameServer/gslog"
//)
//
//var (
//	TReadTimeoutWaitInterval = 20 * time.Microsecond
//
//	ErrNotSameConnectionID = errors.New("not the same connection identify")
//)
//
//type TcpConnFactory func(ctx context.Context, hook ConnectionHook) *net.TCPConn
//
//type TcpConnectionKeeper struct {
//	connID         string
//	version        int
//	keepalive      time.Duration
//	byteOrder      binary.ByteOrder
//	writeTimeout   time.Duration
//	readTimeout    time.Duration
//	tcpConnFactory TcpConnFactory
//
//	stopChan  chan struct{}
//	closed    bool
//	readChan  chan ControlPacket
//	writeChan chan ControlPacket
//
//	ctx       context.Context
//	ctxCancel context.CancelFunc
//
//	lock sync.RWMutex
//}
//
//func NewTcpConnectionKeeper(ctx context.Context, cfg *ConnectionLayerConfig, tcpConnFactory TcpConnFactory) *TcpConnectionKeeper {
//	instance := &TcpConnectionKeeper{
//		connID:         cfg.ConnectionID,
//		version:        cfg.Version,
//		keepalive:      time.Duration(cfg.KeepaliveInterval),
//		byteOrder:      cfg.ByteOrder,
//		writeTimeout:   cfg.WriteTimeout,
//		readTimeout:    cfg.ReadTimeout,
//		tcpConnFactory: tcpConnFactory,
//		stopChan:       make(chan struct{}),
//		readChan:       make(chan ControlPacket),
//		writeChan:      make(chan ControlPacket),
//	}
//
//	instance.ctx, instance.ctxCancel = context.WithCancel(ctx)
//
//	instance.loop()
//
//	return instance
//}
//
//func (gs *TcpConnectionKeeper) loop() {
//	if gs.connID == "" {
//		gslog.Error("[TcpConnectionKeeper] tcp connection keeper loop fail for invalid cid", "connID", gs.connID)
//		return
//	}
//
//	var tcpConn *net.TCPConn
//	//tcpConn := gs.tcpConnFactory(gs.ctx, gs.GetConnectionHook())
//	//if tcpConn == nil {
//	//	gslog.Error("[TcpConnectionKeeper] tcp connection keeper loop fail for create tcp conn", "connID", gs.connID)
//	//	return
//	//}
//
//	go func() {
//		defer func() {
//			close(gs.readChan)
//			close(gs.writeChan)
//			close(gs.stopChan)
//			gs.ctxCancel()
//		}()
//		for {
//			gs.start(tcpConn)
//			if gs.IsClosed() {
//				break
//			}
//			//tcpConn = gs.tcpConnFactory(gs.ctx, gs.GetConnectionHook())
//		}
//	}()
//}
//
//func (gs *TcpConnectionKeeper) start(tcpConn *net.TCPConn) {
//	defer tcpConn.Close()
//
//	heartbeatChan := make(chan ControlPacket)
//	taskCtx, taskCtxCancel := context.WithCancel(gs.ctx)
//	defer taskCtxCancel()
//
//	wg := sync.WaitGroup{}
//	if gs.keepalive > 0 {
//		wg.Add(1)
//		go func() {
//			defer wg.Done()
//			defer taskCtxCancel()
//			gs.keepaliveLoop(taskCtx, tcpConn, gs.keepalive, heartbeatChan)
//		}()
//	}
//
//	wg.Add(1)
//	go func() {
//		defer wg.Done()
//		defer taskCtxCancel()
//		gs.readLoop(taskCtx, tcpConn, heartbeatChan)
//	}()
//
//	wg.Add(1)
//	go func() {
//		defer wg.Done()
//		defer taskCtxCancel()
//		gs.writeLoop(taskCtx, tcpConn)
//	}()
//}
//
//func (gs *TcpConnectionKeeper) IsClosed() bool {
//	gs.lock.RLock()
//	defer gs.lock.RUnlock()
//	return gs.closed
//}
//
//func (gs *TcpConnectionKeeper) keepaliveLoop(ctx context.Context, conn *net.TCPConn, hbInterval time.Duration, heartbeatChan chan ControlPacket) {
//	ticker := time.NewTicker(hbInterval)
//	defer ticker.Stop()
//
//	for {
//		select {
//		case <-ctx.Done():
//			return
//		case <-ticker.C:
//			// heartbeat req
//			packet := NewControlPacket(Heartbeat)
//			select {
//			case gs.writeChan <- packet:
//				hbCtx, hbCtxCancel := context.WithCancel(ctx)
//				select {
//				case <-hbCtx.Done():
//					hbCtxCancel()
//					gslog.Error("[TcpConnectionKeeper] connect heartbeat timeout", "connID", gs.connID, "heartbeat", hbInterval)
//					return
//				case heartbeatChan <- packet:
//					hbCtxCancel()
//				}
//			}
//		}
//	}
//}
//
//func (gs *TcpConnectionKeeper) readLoop(ctx context.Context, conn *net.TCPConn, heartbeatChan chan ControlPacket) {
//	for {
//		if gs.readTimeout > 0 {
//			_ = conn.SetReadDeadline(time.Now().Add(gs.readTimeout))
//		}
//		packet, err := ReadPacket(conn, gs.byteOrder)
//		if err != nil {
//			var netErr net.Error
//			if errors.As(err, &netErr) && netErr.Timeout() {
//				time.Sleep(TReadTimeoutWaitInterval)
//				continue
//			}
//			break
//		}
//		if gs.readTimeout > 0 {
//			_ = conn.SetReadDeadline(time.Time{})
//		}
//		gslog.Trace("[TcpConnectionKeeper] readLoop receiver", "connID", gs.connID, "packet", packet.String())
//		switch msg := packet.(type) {
//		case *HeartbeatPacket:
//			// send heartbeat ack
//			ackPacket := NewControlPacket(HeartbeatAck)
//			select {
//			case gs.writeChan <- ackPacket:
//			case <-ctx.Done():
//				return
//			}
//		case *HeartbeatAckPacket:
//			// dispatch to heartbeat channel
//			select {
//			case heartbeatChan <- packet:
//			case <-ctx.Done():
//				return
//			}
//		case *DisConnectPacket:
//			return
//		default:
//			select {
//			case <-ctx.Done():
//				return
//			case gs.readChan <- msg:
//			}
//		}
//	}
//	return
//}
//
//func (gs *TcpConnectionKeeper) writeLoop(ctx context.Context, conn *net.TCPConn) {
//	var err error
//	for {
//		select {
//		case <-ctx.Done():
//			return
//		case packet, ok := <-gs.writeChan:
//			if !ok {
//				// channel close
//				break
//			}
//			if gs.writeTimeout > 0 {
//				_ = conn.SetWriteDeadline(time.Now().Add(gs.writeTimeout))
//			}
//			_, err = packet.WriteTo(conn, gs.byteOrder)
//			if err != nil {
//				var netErr net.Error
//				if errors.As(err, &netErr) && netErr.Timeout() {
//					time.Sleep(TReadTimeoutWaitInterval)
//					continue
//				}
//				return
//			}
//			if gs.writeTimeout > 0 {
//				_ = conn.SetWriteDeadline(time.Time{})
//			}
//			gslog.Trace("[TcpConnectionKeeper] writeLoop sender", "connID", gs.connID, "packet", packet.String())
//		}
//	}
//}
//
//func (gs *TcpConnectionKeeper) ConnectionID() string {
//	gs.lock.RLock()
//	defer gs.lock.RUnlock()
//	return gs.connID
//}
//
////func (gs *TcpConnectionKeeper) GetConnectionHook() ConnectionHook {
////	return func(conn net.Conn) error {
////		// read connect packet
////		if gs.readTimeout > 0 {
////			_ = conn.SetReadDeadline(time.Now().Add(gs.readTimeout))
////		}
////		msg, err := ReadPacket(conn, gs.byteOrder)
////		if err != nil {
////			return err
////		}
////		if gs.readTimeout > 0 {
////			_ = conn.SetReadDeadline(time.Time{})
////		}
////		packet, ok := msg.(*ConnectPacket)
////		if !ok || packet.Validate() != Accepted {
////			return err
////		}
////		if gs.ConnectionID() != packet.ClientIdentifier {
////			return ErrNotSameConnectionID
////		}
////		gs.lock.Lock()
////		gs.version = packet.ProtocolVersion
////		gs.keepalive = time.Duration(packet.Keepalive)
////		gs.lock.Unlock()
////
////		// send connection ack packet
////		ackPacket := NewControlPacket(ConnectAck).(*ConnectAckPacket)
////		if gs.writeTimeout > 0 {
////			_ = conn.SetWriteDeadline(time.Now().Add(gs.writeTimeout))
////		}
////		_, err = ackPacket.WriteTo(conn, gs.byteOrder)
////		if err != nil {
////			return err
////		}
////		if gs.writeTimeout > 0 {
////			_ = conn.SetWriteDeadline(time.Time{})
////		}
////
////		return nil
////	}
////}
