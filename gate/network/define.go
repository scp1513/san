package network

import (
	"github.com/golang/protobuf/proto"
)

type _Message interface {
	proto.Message
	GetHeader() uint32
}
