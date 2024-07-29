package network

import (
	"encoding/binary"
	"io"
	"net"
)

// ConnectionHook 新建连接的钩子函数
type ConnectionHook func(conn net.Conn) error

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
		// GetConnectionHook 获取一个新建立连接的钩子函数
		GetConnectionHook() ConnectionHook
		//Read() chan ControlPacket
		//ReadPacket() ControlPacket
		//Write(ctx context.Context, packet ControlPacket) error
		//Close() error
	}

	////Connection 连接抽象封装
	//Connection interface {
	//	net.Conn
	//	// ConnectionID 连接标识
	//	ConnectionID() string
	//	// Heartbeat 心跳时间
	//	Heartbeat() time.Duration
	//	//ReadPacket() ControlPacket
	//	//WritePacket(ctx context.Context, packet ControlPacket) error
	//}

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
