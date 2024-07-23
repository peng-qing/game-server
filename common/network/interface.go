package network

import (
	"encoding/binary"
	"io"
	"net"
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
	// 连接的抽象
	ConnectionLayer interface {
		net.Conn
		// ConnectionID 连接ID
		ConnectionID() string
		// Heartbeat 心跳时间
		Heartbeat() time.Duration
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
