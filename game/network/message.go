package network

import (
	"github.com/golang/protobuf/proto"
)

type Message interface {
	proto.Message
	GetHeader() uint32
}
