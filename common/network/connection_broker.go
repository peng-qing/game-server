package network

import (
	"net"
	"time"
)

type Broker struct {
	conn          net.Conn
	connectionID  string
	version       int
	keepalive     time.Duration
	readTimeout   time.Duration
	writeTimeout  time.Duration
	closeCallback OnConnectionCloseCallback
}

func (gs *Broker) ConnectionID() string {
	return gs.connectionID
}

func (gs *Broker) Keepalive() time.Duration {
	return gs.keepalive
}

func (gs *Broker) WritePacket(packet ControlPacket) error {
	//TODO implement me
	panic("implement me")
}

func (gs *Broker) ReadPacket() (ControlPacket, error) {
	//TODO implement me
	panic("implement me")
}

func (gs *Broker) LocalAddr() string {
	return gs.conn.LocalAddr().String()
}

func (gs *Broker) RemoteAddr() string {
	return gs.conn.RemoteAddr().String()
}

func (gs *Broker) Close() error {
	err := gs.conn.Close()
	if err == nil && gs.closeCallback != nil {
		gs.closeCallback(gs.connectionID)
	}
	return err
}

//func AcceptTcpConn(cfg *ConnectionConfig) (TcpConnFactory, error) {
//	addr, err := net.ResolveTCPAddr(cfg.IPVersion, fmt.Sprintf("%s:%d", cfg.IP, cfg.Port))
//	if err != nil {
//		return nil, err
//	}
//	listener, err := net.ListenTCP(cfg.IPVersion, addr)
//	if err != nil {
//		gslog.Error("[AcceptTcpConn] listen tcp failed", "ipVersion", cfg.IPVersion, "ip", cfg.IP, "port", cfg.Port)
//		return nil, err
//	}
//
//	return func(ctx context.Context, hook ConnectionHook) *net.TCPConn {
//		conn, err := listener.AcceptTCP()
//		if err != nil {
//			gslog.Critical("[TcpConnFactory] accept tcp conn failed", "addr", listener.Addr().String(), "ipVersion", cfg.IPVersion, "err", err)
//			return nil
//		}
//		if hook != nil {
//			err = hook(conn)
//			if err != nil {
//				return nil
//			}
//		}
//		return conn
//	}, nil
//}
//
//func ConnectTcpConn(cfg *ConnectionConfig, hook ConnectionHook) (TcpConnFactory, error) {
//	addr, err := net.ResolveTCPAddr(cfg.IPVersion, fmt.Sprintf("%s:%d", cfg.IP, cfg.Port))
//	if err != nil {
//		return nil, err
//	}
//	return func(ctx context.Context, hook ConnectionHook) *net.TCPConn {
//		conn, err := net.DialTCP(cfg.IPVersion, nil, addr)
//		if err != nil {
//			gslog.Critical("[TcpConnFactory] connect tcp failed", "addr", addr.String(), "ipVersion", cfg.IPVersion, "err", err)
//			return nil
//		}
//		if hook != nil {
//			err = hook(conn)
//			if err != nil {
//				return nil
//			}
//		}
//		return conn
//	}, nil
//}
