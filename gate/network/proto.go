package network

import (
	"fmt"
	"reflect"
)

type regInfo struct {
	m  reflect.Type
	h  func(uint, interface{})
	ih func(interface{})
}

var regProto = make(map[uint32]*regInfo)

// RegisterProto 注册协议
func RegisterProto(msg _Message, handler interface{}) {
	header := msg.GetHeader()
	_, ok := regProto[header]
	if ok {
		panic(fmt.Errorf("重复注册协议 %d", header))
	}

	switch h := handler.(type) {
	default:
		panic(fmt.Errorf("无效handler %d", header))
	case func(uint, interface{}):
		pi := &regInfo{h: h}
		if msg != nil {
			pi.m = reflect.TypeOf(msg).Elem()
		}
		regProto[header] = pi
	case func(interface{}):
		pi := &regInfo{ih: h}
		if msg != nil {
			pi.m = reflect.TypeOf(msg).Elem()
		}
		regProto[header] = pi
	}
}
