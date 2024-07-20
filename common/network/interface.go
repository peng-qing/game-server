package network

import "io"

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
		// ConnectionID 连接ID
		ConnectionID() int64
	}

	// ControlPacket 连接层控制报文
	ControlPacket interface {
		Name() string
		String() string
		Validate() int
	}

	// ControlPacker 连接层控制报文封包/解包
	ControlPacker interface {
		Pack(packet ControlPacket) ([]byte, error)
		Unpack(r io.Reader) (ControlPacket, error)
		WriteTo(w io.Writer, packet ControlPacket) (int64, error)
	}
)
