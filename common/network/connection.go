package network

import (
	"context"
	"fmt"
	"net"

	"GameServer/gslog"
)

func AcceptTcpConn(cfg *ConnectionConfig) (TcpConnFactory, error) {
	addr, err := net.ResolveTCPAddr(cfg.IPVersion, fmt.Sprintf("%s:%d", cfg.IP, cfg.Port))
	if err != nil {
		return nil, err
	}
	listener, err := net.ListenTCP(cfg.IPVersion, addr)
	if err != nil {
		gslog.Error("[AcceptTcpConn] listen tcp failed", "ipVersion", cfg.IPVersion, "ip", cfg.IP, "port", cfg.Port)
		return nil, err
	}

	// 是否应该是一个 Connection Hook 而不是所有参数都传递进去...
	return func(ctx context.Context, hook ConnectionHook) *net.TCPConn {
		conn, err := listener.AcceptTCP()
		if err != nil {
			gslog.Critical("[TcpConnFactory] accept tcp conn failed", "addr", listener.Addr().String())
			return nil
		}
		if hook != nil {
			err = hook(conn)
			if err != nil {
				return nil
			}
		}
		return conn
	}, nil
}

// import (
//
//	"encoding/binary"
//	"errors"
//	"io"
//	"net"
//	"time"
//
// )

//	type ConnectionConfig struct {
//		ConnectionID      string
//		KeepaliveInterval int
//		Version           int
//		WriteTimeout      time.Duration
//		ReadTimeout       time.Duration
//		ByteOrder         binary.ByteOrder
//	}

//
//// ConnectServer 连接到服务器
//func ConnectServer(conn net.Conn, cfg *ConnectionConfig) (Connection, error) {
//	// send connect packet
//	packet := NewControlPacket(Connect).(*ConnectPacket)
//	packet.ClientIdentifier = cfg.ConnectionID
//	packet.Keepalive = cfg.KeepaliveInterval
//	packet.ProtocolVersion = cfg.Version
//
//	if cfg.WriteTimeout != 0 {
//		_ = conn.SetWriteDeadline(time.Now().Add(cfg.WriteTimeout))
//	}
//	_, err := packet.WriteTo(conn, cfg.ByteOrder)
//	if err != nil {
//		return nil, err
//	}
//	if cfg.WriteTimeout != 0 {
//		_ = conn.SetWriteDeadline(time.Time{})
//	}
//	// read connect_ack packet
//	if cfg.ReadTimeout != 0 {
//		_ = conn.SetReadDeadline(time.Now().Add(cfg.ReadTimeout))
//	}
//	ackPacket, err := ReadPacket(conn, cfg.ByteOrder)
//	if err != nil {
//		return nil, err
//	}
//	ack, ok := ackPacket.(*ConnectAckPacket)
//	if !ok {
//		return nil, err
//	}
//	if cfg.ReadTimeout != 0 {
//		_ = conn.SetReadDeadline(time.Time{})
//	}
//	if ack.ReturnCode != Accepted {
//		return nil, RetCodeErrors[ack.ReturnCode]
//	}
//
//	return NewConnection(conn, cfg), nil
//}

//
//type connection struct {
//	net.Conn
//	connID       string
//	keepalive    time.Duration
//	byteOrder    binary.ByteOrder
//	version      int
//	writeTimeout time.Duration
//	readTimeout  time.Duration
//}
//
//func NewConnection(conn net.Conn, cfg *ConnectionConfig) Connection {
//	return &connection{
//		Conn:         conn,
//		connID:       cfg.ConnectionID,
//		keepalive:    time.Duration(cfg.KeepaliveInterval),
//		byteOrder:    cfg.ByteOrder,
//		version:      cfg.Version,
//		writeTimeout: cfg.WriteTimeout,
//		readTimeout:  cfg.ReadTimeout,
//	}
//}
//
//func (gs *connection) ConnectionID() string {
//	return gs.connID
//}
//
//func (gs *connection) Heartbeat() time.Duration {
//	return gs.keepalive
//}
//
//// Write 重写覆盖net.Conn Write
//func (gs *connection) Write(b []byte) (n int, err error) {
//	if gs.writeTimeout != 0 {
//		_ = gs.Conn.SetWriteDeadline(time.Now().Add(gs.writeTimeout))
//	}
//	n, err = gs.Conn.Write(b)
//	if gs.writeTimeout != 0 {
//		_ = gs.Conn.SetWriteDeadline(time.Time{})
//	}
//
//	return
//}
//
//// Read 重写覆盖net.Conn Read
//func (gs *connection) Read(b []byte) (n int, err error) {
//	if gs.readTimeout != 0 {
//		_ = gs.Conn.SetReadDeadline(time.Now().Add(gs.readTimeout))
//	}
//	n, err = gs.Conn.Read(b)
//	if gs.readTimeout != 0 {
//		_ = gs.Conn.SetReadDeadline(time.Time{})
//	}
//	return
//}
