package network

import (
	"GameServer/utils"
	"bytes"
	"errors"
	"fmt"
	"io"
)

type PacketType byte

const (
	Invalid PacketType = iota
	Connect
	ConnectAck
	Heartbeat
	HeartbeatAck
	Publish
	PublishAck
	DisConnect
)

const (
	Accepted                  = 0
	RefusedBadProtocolVersion = 1
	RefusedInvalidIdentifier  = 2
)

var (
	PacketNames = []string{
		"INVALID",
		"CONNECT",
		"CONNECT_ACK",
		"HEARTBEAT",
		"HEARTBEAT_ACK",
		"PUBLISH",
		"PUBLISH_ACK",
		"DISCONNECT",
	}

	ErrBadProtocolVersion       = errors.New("connection refused: bad protocol version")
	ErrRefusedInvalidIdentifier = errors.New("connection refused: invalid client identifier")

	RetcodeErrors = map[int]error{
		Accepted:                  nil,
		RefusedBadProtocolVersion: ErrBadProtocolVersion,
		RefusedInvalidIdentifier:  ErrRefusedInvalidIdentifier,
	}
)

// FixedHeader 固定包头 所有控制报文都含有
type FixedHeader struct {
	PacketType   PacketType
	RemainLength int
}

func (gs *FixedHeader) String() string {
	return fmt.Sprintf("[%s] RemainLength: %d", PacketNames[gs.PacketType], gs.RemainLength)
}

// Pack 打包固定包头
func (gs *FixedHeader) Pack() bytes.Buffer {
	var header bytes.Buffer

	header.WriteByte(byte(gs.PacketType))
	header.Write(utils.EncodeVariable(int64(gs.RemainLength)))

	return header
}

// UnPack 解包固定包头
func (gs *FixedHeader) UnPack(packetType PacketType, r io.Reader) error {
	var err error
	gs.PacketType = packetType
	gs.RemainLength, err = utils.DecodeReaderVariableInt(r)

	return err
}

type ConnectPacket struct {
	FixedHeader
	ProtocolVersion  int    // 协议版本
	Keepalive        int    // 心跳时间
	ClientIdentifier string // 客户端唯一标识
}

func (gs *ConnectPacket) Name() string {
	return PacketNames[gs.FixedHeader.PacketType]
}

func (gs *ConnectPacket) Validate() int {
	if len(gs.ClientIdentifier) <= 0 ||
		len(gs.ClientIdentifier) > 65535 {
		return RefusedInvalidIdentifier
	}

	return 0
}

func (gs *ConnectPacket) String() string {
	return fmt.Sprintf("%s ,protocolVersion:%d, keepalive:%d, clientIdentifier:%s",
		gs.FixedHeader.String(), gs.ProtocolVersion, gs.Keepalive, gs.ClientIdentifier)
}
