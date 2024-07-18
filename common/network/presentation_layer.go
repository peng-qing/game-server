package network

import (
	"encoding/json"

	"github.com/golang/protobuf/proto"

	"GameServer/gslog"
)

// 提供表示层的相关实现

// JsonPresentation Json形式
type JsonPresentation struct{}

func NewJsonPresentation() PresentationLayer {
	return &JsonPresentation{}
}

func (j JsonPresentation) Decode(src []byte, dst any) error {
	return json.Unmarshal(src, dst)
}

func (j JsonPresentation) Encode(src any) (dst []byte, err error) {
	return json.Marshal(src)
}

// PBPresentation Proto Buffer
type PBPresentation struct{}

func NewPBPresentation() PresentationLayer {
	return &PBPresentation{}
}

func (P PBPresentation) Decode(src []byte, dst any) error {
	// must pb.Message
	return proto.Unmarshal(src, dst.(proto.Message))
}

func (P PBPresentation) Encode(src any) (dst []byte, err error) {
	switch vv := src.(type) {
	case []byte:
		return vv, err
	case proto.Message:
		return proto.Marshal(vv)
	}

	gslog.Critical("[PBPresentation] proto marshal failed...", "src", src)
	return nil, err
}
