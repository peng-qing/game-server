package network

import (
	"context"
	"encoding/binary"
	"io"
	"time"
)

type (
	// ApplicationLayer 应用层
	// 最顶层业务 表示一个服务
	ApplicationLayer interface {
	}

	//SessionLayer 会话层
	// 管理&控制两个通信之间的会话,数据交换同步
	SessionLayer interface {
	}

	// PresentationLayer 表示层
	// 应用层数据格式转换, 加密和压缩
	PresentationLayer interface {
		Decode(src []byte, dst any) error
		Encode(src any) (dst []byte, err error)
	}

	// ConnectionLayer 连接层
	ConnectionLayer interface {
		// ConnectionID 连接ID
		ConnectionID() string
		// Close 关闭
		Close() error
		// Read 收包队列
		Read() chan ControlPacket
		// ReadPacket 读取单个包
		ReadPacket() ControlPacket
		// WritePacket 写入包
		WritePacket(ctx context.Context, packet ControlPacket) error
	}

	//ConnectionBroker 连接代理
	ConnectionBroker interface {
		// ConnectionID 连接ID
		ConnectionID() string
		// Keepalive 心跳时间
		Keepalive() time.Duration
		// WritePacket 发包
		WritePacket(packet ControlPacket) error
		// ReadPacket 读取单个包
		ReadPacket() (ControlPacket, error)
		// LocalAddr 本地地址
		LocalAddr() string
		// RemoteAddr 远端地址
		RemoteAddr() string
		// Close 关闭
		Close() error
	}

	// ControlPacket 连接层控制报文
	ControlPacket interface {
		Name() string
		String() string
		Validate() int
		Pack(order binary.ByteOrder) ([]byte, error)
		Unpack(r io.Reader, order binary.ByteOrder) error
		WriteTo(w io.Writer, order binary.ByteOrder) (int64, error)
	}
)
