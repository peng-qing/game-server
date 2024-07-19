package network

import (
	"GameServer/utils"
	"bytes"
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
)

// FixedHeader 固定包头 所有控制报文都含有
type FixedHeader struct {
	PacketType   PacketType
	RemainLength int
}

func (gs *FixedHeader) String() string {
	return fmt.Sprintf("PacketType: %d, RemainLength: %d", gs.PacketType, gs.RemainLength)
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
