package network

import (
	"GameServer/utils"
	"bytes"
	"encoding/binary"
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

	ErrInvalidPacketType        = errors.New("invalid packet type")
	ErrBadProtocolVersion       = errors.New("connection refused: bad protocol version")
	ErrRefusedInvalidIdentifier = errors.New("connection refused: invalid client identifier")

	RetCodeErrors = map[int]error{
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

func (gs *FixedHeader) Name() string {
	return PacketNames[gs.PacketType]
}

// Pack 打包固定包头
func (gs *FixedHeader) Pack() bytes.Buffer {
	var header bytes.Buffer

	header.WriteByte(byte(gs.PacketType))
	header.Write(utils.EncodeVariableInt(int64(gs.RemainLength)))

	return header
}

// UnPack 解包固定包头
func (gs *FixedHeader) UnPack(packetType PacketType, r io.Reader) error {
	var err error
	gs.PacketType = packetType
	gs.RemainLength, err = utils.DecodeReaderVariableInt(r)

	return err
}

// NewControlPacket 创建协议包
func NewControlPacket(packetType PacketType) ControlPacket {
	fixedHeader := FixedHeader{PacketType: packetType}

	switch packetType {
	case Connect:
		return &ConnectPacket{FixedHeader: fixedHeader}
	case ConnectAck:
		return &ConnectAckPacket{FixedHeader: fixedHeader}
	case Heartbeat:
		return &HeartbeatPacket{FixedHeader: fixedHeader}
	case HeartbeatAck:
		return &HeartbeatAckPacket{FixedHeader: fixedHeader}
	case Publish:
		return &PublishPacket{FixedHeader: fixedHeader}
	case PublishAck:
		return &PublishAckPacket{FixedHeader: fixedHeader}
	case DisConnect:
		return &DisConnectPacket{FixedHeader: fixedHeader}
	case Invalid:
		fallthrough
	default:
		return nil
	}
}

func NewControlPacketWithHeader(fh FixedHeader) (ControlPacket, error) {
	switch fh.PacketType {
	case Connect:
		return &ConnectPacket{FixedHeader: fh}, nil
	case ConnectAck:
		return &ConnectAckPacket{FixedHeader: fh}, nil
	case Heartbeat:
		return &HeartbeatPacket{FixedHeader: fh}, nil
	case HeartbeatAck:
		return &HeartbeatAckPacket{FixedHeader: fh}, nil
	case Publish:
		return &PublishPacket{FixedHeader: fh}, nil
	case PublishAck:
		return &PublishAckPacket{FixedHeader: fh}, nil
	case DisConnect:
		return &DisConnectPacket{FixedHeader: fh}, nil
	case Invalid:
		fallthrough
	default:
		return nil, ErrInvalidPacketType
	}
}

type ConnectPacket struct {
	FixedHeader
	ProtocolVersion  int    // 协议版本
	Keepalive        int    // 心跳时间
	ClientIdentifier string // 客户端唯一标识
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

func (gs *ConnectPacket) Pack(order binary.ByteOrder) ([]byte, error) {
	var body bytes.Buffer
	var err error

	err = binary.Write(&body, order, gs.ProtocolVersion)
	err = binary.Write(&body, order, gs.Keepalive)
	err = binary.Write(&body, order, utils.EncodeString(gs.ClientIdentifier, order))

	gs.FixedHeader.RemainLength = body.Len()
	packet := gs.FixedHeader.Pack()
	packet.Write(body.Bytes())

	return packet.Bytes(), err
}

func (gs *ConnectPacket) Unpack(r io.Reader, order binary.ByteOrder) error {
	var err error
	err = binary.Read(r, order, &gs.ProtocolVersion)
	err = binary.Read(r, order, &gs.Keepalive)
	gs.ClientIdentifier, err = utils.DecodeReaderString(r, order)

	return err
}

func (gs *ConnectPacket) WriteTo(w io.Writer, order binary.ByteOrder) (int64, error) {
	data, err := gs.Pack(order)
	if err != nil {
		return 0, err
	}
	n, err := w.Write(data)
	return int64(n), err
}

type ConnectAckPacket struct {
	FixedHeader
	ReturnCode int
}

func (gs *ConnectAckPacket) Validate() int {

	return 0
}

func (gs *ConnectAckPacket) String() string {
	return fmt.Sprintf("%s , returnCode:%d", gs.FixedHeader.String(), gs.ReturnCode)
}

func (gs *ConnectAckPacket) Pack(order binary.ByteOrder) ([]byte, error) {
	var body bytes.Buffer
	var err error

	err = binary.Write(&body, order, gs.ReturnCode)

	gs.FixedHeader.RemainLength = body.Len()
	packet := gs.FixedHeader.Pack()
	packet.Write(body.Bytes())

	return packet.Bytes(), err
}

func (gs *ConnectAckPacket) Unpack(r io.Reader, order binary.ByteOrder) error {
	var err error
	err = binary.Read(r, order, &gs.ReturnCode)

	return err
}

func (gs *ConnectAckPacket) WriteTo(w io.Writer, order binary.ByteOrder) (int64, error) {
	data, err := gs.Pack(order)
	if err != nil {
		return 0, err
	}
	n, err := w.Write(data)
	return int64(n), err
}

type HeartbeatPacket struct {
	FixedHeader
}

func (gs *HeartbeatPacket) Validate() int {
	return 0
}

func (gs *HeartbeatPacket) Pack(order binary.ByteOrder) ([]byte, error) {
	packet := gs.FixedHeader.Pack()

	return packet.Bytes(), nil
}

func (gs *HeartbeatPacket) Unpack(r io.Reader, order binary.ByteOrder) error {
	return nil
}

func (gs *HeartbeatPacket) WriteTo(w io.Writer, order binary.ByteOrder) (int64, error) {
	data, err := gs.Pack(order)
	if err != nil {
		return 0, err
	}
	n, err := w.Write(data)
	return int64(n), err
}

type HeartbeatAckPacket struct {
	FixedHeader
}

func (gs *HeartbeatAckPacket) Validate() int {
	return 0
}

func (gs *HeartbeatAckPacket) Pack(order binary.ByteOrder) ([]byte, error) {
	packet := gs.FixedHeader.Pack()

	return packet.Bytes(), nil
}

func (gs *HeartbeatAckPacket) Unpack(r io.Reader, order binary.ByteOrder) error {
	return nil
}

func (gs *HeartbeatAckPacket) WriteTo(w io.Writer, order binary.ByteOrder) (int64, error) {
	data, err := gs.Pack(order)
	if err != nil {
		return 0, err
	}
	n, err := w.Write(data)
	return int64(n), err
}

type DisConnectPacket struct {
	FixedHeader
}

func (gs *DisConnectPacket) Validate() int {
	return 0
}

func (gs *DisConnectPacket) Pack(order binary.ByteOrder) ([]byte, error) {
	packet := gs.FixedHeader.Pack()

	return packet.Bytes(), nil
}

func (gs *DisConnectPacket) Unpack(r io.Reader, order binary.ByteOrder) error {
	return nil
}

func (gs *DisConnectPacket) WriteTo(w io.Writer, order binary.ByteOrder) (int64, error) {
	data, err := gs.Pack(order)
	if err != nil {
		return 0, err
	}
	n, err := w.Write(data)
	return int64(n), err
}

type PublishPacket struct {
	FixedHeader
	MessageID uint32
	Payload   []byte
}

func (gs *PublishPacket) Validate() int {
	return 0
}

func (gs *PublishPacket) String() string {
	return fmt.Sprintf("%s , MessageID:%d, Payload:%s", gs.FixedHeader.String(), gs.MessageID, string(gs.Payload))
}

func (gs *PublishPacket) Pack(order binary.ByteOrder) ([]byte, error) {
	var body bytes.Buffer
	var err error

	err = binary.Write(&body, order, gs.MessageID)
	err = binary.Write(&body, order, utils.EncodeBytes(gs.Payload, order))

	gs.FixedHeader.RemainLength = body.Len()
	packet := gs.FixedHeader.Pack()
	packet.Write(body.Bytes())

	return packet.Bytes(), err
}

func (gs *PublishPacket) Unpack(r io.Reader, order binary.ByteOrder) error {
	var err error
	err = binary.Read(r, order, &gs.MessageID)
	gs.Payload, err = utils.DecodeReaderBytes(r, order)

	return err
}

func (gs *PublishPacket) WriteTo(w io.Writer, order binary.ByteOrder) (int64, error) {
	data, err := gs.Pack(order)
	if err != nil {
		return 0, err
	}
	n, err := w.Write(data)
	return int64(n), err
}

type PublishAckPacket struct {
	FixedHeader
	MessageID uint32
}

func (gs *PublishAckPacket) Validate() int {
	return 0
}

func (gs *PublishAckPacket) String() string {
	return fmt.Sprintf("%s , MessageID:%d, Payload:%s", gs.FixedHeader.String(), gs.MessageID)
}

func (gs *PublishAckPacket) Pack(order binary.ByteOrder) ([]byte, error) {
	var body bytes.Buffer
	var err error

	err = binary.Write(&body, order, gs.MessageID)

	gs.FixedHeader.RemainLength = body.Len()
	packet := gs.FixedHeader.Pack()
	packet.Write(body.Bytes())

	return packet.Bytes(), err
}

func (gs *PublishAckPacket) Unpack(r io.Reader, order binary.ByteOrder) error {
	var err error
	err = binary.Read(r, order, &gs.MessageID)

	return err
}

func (gs *PublishAckPacket) WriteTo(w io.Writer, order binary.ByteOrder) (int64, error) {
	data, err := gs.Pack(order)
	if err != nil {
		return 0, err
	}
	n, err := w.Write(data)
	return int64(n), err
}
