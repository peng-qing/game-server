package network

import (
	"GameServer/gslog"
	"encoding/binary"
	"net"
	"time"
)

// ConnectBroker 客户端连接到服务器
func ConnectBroker(conn net.Conn, cfg *BrokerConf) (ConnectionBroker, error) {
	var err error
	// send connect
	packet := NewControlPacket(Connect).(*ConnectPacket)
	packet.Keepalive = cfg.KeepaliveInterval
	packet.ClientIdentifier = cfg.ConnectionID
	packet.ProtocolVersion = cfg.Version

	if cfg.WriteTimeout > 0 {
		_ = conn.SetWriteDeadline(time.Now().Add(cfg.WriteTimeout))
	}
	_, err = packet.WriteTo(conn, cfg.ByteOrder)
	if err != nil {
		return nil, err
	}
	if cfg.WriteTimeout > 0 {
		_ = conn.SetWriteDeadline(time.Time{})
	}
	// read connect ack
	if cfg.ReadTimeout > 0 {
		_ = conn.SetReadDeadline(time.Now().Add(cfg.ReadTimeout))
	}
	ackPacket, err := ReadPacket(conn, cfg.ByteOrder)
	if err != nil {
		return nil, err
	}
	if cfg.ReadTimeout > 0 {
		_ = conn.SetReadDeadline(time.Time{})
	}
	if ret := ackPacket.Validate(); ret != Accepted {
		return nil, RetCodeErrors[ret]
	}
	connectAck, ok := ackPacket.(*ConnectAckPacket)
	if !ok {
		gslog.Warn("read ack packet not connect ack", "ack", ackPacket.String())
		return nil, err
	}
	if connectAck.ReturnCode != Accepted {
		return nil, RetCodeErrors[connectAck.ReturnCode]
	}

	return NewConnectionBroker(conn, cfg), nil
}

// AcceptBroker 接受客户端连接
func AcceptBroker(conn net.Conn, cfg *BrokerConf) (ConnectionBroker, error) {
	// read connect packet
	if cfg.ReadTimeout > 0 {
		_ = conn.SetReadDeadline(time.Now().Add(cfg.ReadTimeout))
	}
	packet, err := ReadPacket(conn, cfg.ByteOrder)
	if err != nil {
		return nil, err
	}
	if cfg.ReadTimeout > 0 {
		_ = conn.SetReadDeadline(time.Time{})
	}
	connect, ok := packet.(*ConnectPacket)
	if !ok {
		return nil, err
	}
	cfg.Version = connect.ProtocolVersion
	cfg.ConnectionID = connect.ClientIdentifier
	cfg.KeepaliveInterval = connect.Keepalive

	// send connect ack
	connectAck := NewControlPacket(ConnectAck).(*ConnectAckPacket)
	connectAck.ReturnCode = Accepted
	if cfg.WriteTimeout > 0 {
		_ = conn.SetWriteDeadline(time.Now().Add(cfg.WriteTimeout))
	}
	_, err = connectAck.WriteTo(conn, cfg.ByteOrder)
	if err != nil {
		return nil, err
	}
	if cfg.WriteTimeout > 0 {
		_ = conn.SetWriteDeadline(time.Time{})
	}

	return NewConnectionBroker(conn, cfg), nil
}

type connBroker struct {
	conn          net.Conn
	connectionID  string
	version       int
	keepalive     time.Duration
	readTimeout   time.Duration
	writeTimeout  time.Duration
	byteOrder     binary.ByteOrder
	closeCallback OnConnectionCloseCallback
}

func NewConnectionBroker(conn net.Conn, cfg *BrokerConf) ConnectionBroker {
	return &connBroker{
		conn:          conn,
		connectionID:  cfg.ConnectionID,
		version:       cfg.Version,
		keepalive:     time.Duration(cfg.KeepaliveInterval),
		readTimeout:   cfg.ReadTimeout,
		writeTimeout:  cfg.WriteTimeout,
		byteOrder:     cfg.ByteOrder,
		closeCallback: cfg.OnCloseCallback,
	}
}

func (gs *connBroker) ConnectionID() string {
	return gs.connectionID
}

func (gs *connBroker) Keepalive() time.Duration {
	return gs.keepalive
}

func (gs *connBroker) WritePacket(packet ControlPacket) error {
	if gs.writeTimeout > 0 {
		_ = gs.conn.SetWriteDeadline(time.Now().Add(gs.writeTimeout))
	}
	_, err := packet.WriteTo(gs.conn, gs.byteOrder)
	if err != nil {
		return err
	}
	if gs.writeTimeout > 0 {
		_ = gs.conn.SetWriteDeadline(time.Time{})
	}

	return nil
}

func (gs *connBroker) ReadPacket() (ControlPacket, error) {
	if gs.readTimeout > 0 {
		_ = gs.conn.SetReadDeadline(time.Now().Add(gs.readTimeout))
	}
	packet, err := ReadPacket(gs.conn, gs.byteOrder)
	if err != nil {
		return nil, err
	}
	if gs.readTimeout > 0 {
		_ = gs.conn.SetReadDeadline(time.Time{})
	}

	return packet, nil
}

func (gs *connBroker) LocalAddr() string {
	return gs.conn.LocalAddr().String()
}

func (gs *connBroker) RemoteAddr() string {
	return gs.conn.RemoteAddr().String()
}

func (gs *connBroker) Close() error {
	err := gs.conn.Close()
	if err == nil && gs.closeCallback != nil {
		gs.closeCallback(gs.connectionID)
	}

	return err
}
