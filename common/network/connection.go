package network

import (
	"encoding/binary"
	"errors"
	"io"
	"net"
	"time"
)

var (
	ErrReadExpectedDataFailed = errors.New("read expected data failed")
)

type ConnectionConfig struct {
	ConnectionID      string
	KeepaliveInterval int
	Version           int
	WriteTimeout      time.Duration
	ReadTimeout       time.Duration
	ByteOrder         binary.ByteOrder
}

func ReadPacket(r io.Reader, order binary.ByteOrder) (ControlPacket, error) {
	var fixedHeader FixedHeader
	var err error
	buf := make([]byte, 1)

	if _, err = io.ReadFull(r, buf); err != nil {
		return nil, err
	}
	err = fixedHeader.UnPack(PacketType(buf[0]), r)
	if err != nil {
		return nil, err
	}
	packet, err := NewControlPacketWithHeader(fixedHeader)
	if err != nil {
		return nil, err
	}
	bodySize := make([]byte, fixedHeader.RemainLength)
	n, err := io.ReadFull(r, bodySize)
	if err != nil {
		return nil, err
	}
	if n != fixedHeader.RemainLength {
		return nil, ErrReadExpectedDataFailed
	}

	err = packet.Unpack(r, order)

	return packet, err
}

// ConnectServer 连接到服务器
func ConnectServer(conn net.Conn, cfg *ConnectionConfig) (ConnectionLayer, error) {
	// send connect packet
	packet := NewControlPacket(Connect).(*ConnectPacket)
	packet.ClientIdentifier = cfg.ConnectionID
	packet.Keepalive = cfg.KeepaliveInterval
	packet.ProtocolVersion = cfg.Version

	if cfg.WriteTimeout != 0 {
		_ = conn.SetWriteDeadline(time.Now().Add(cfg.WriteTimeout))
	}
	_, err := packet.WriteTo(conn, cfg.ByteOrder)
	if err != nil {
		return nil, err
	}
	if cfg.WriteTimeout != 0 {
		_ = conn.SetWriteDeadline(time.Time{})
	}
	// read connect_ack packet
	if cfg.ReadTimeout != 0 {
		_ = conn.SetReadDeadline(time.Now().Add(cfg.ReadTimeout))
	}
	ackPacket, err := ReadPacket(conn, cfg.ByteOrder)
	if err != nil {
		return nil, err
	}
	ack, ok := ackPacket.(*ConnectAckPacket)
	if !ok {
		return nil, err
	}
	if cfg.ReadTimeout != 0 {
		_ = conn.SetReadDeadline(time.Time{})
	}
	if ack.ReturnCode != Accepted {
		return nil, RetCodeErrors[ack.ReturnCode]
	}

	return NewConnectionLayer(conn, cfg), nil
}
