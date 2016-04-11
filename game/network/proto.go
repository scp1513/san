package network

import (
	"fmt"
	"reflect"
)

type regInfo struct {
	m  reflect.Type              // msg
	h  func(uint64, interface{}) // handler
	ih func(uint, interface{})   // inner handler
}

var regProto = make(map[uint32]*regInfo)

// 注册协议
func RegisterProto(msg Message, handler interface{}) {
	header := msg.GetHeader()
	_, ok := regProto[header]
	if ok {
		panic(fmt.Errorf("重复注册协议 %d", header))
	}

	switch h := handler.(type) {
	default:
		panic(fmt.Errorf("无效handler %d", header))
	case func(uint64, interface{}):
		pi := &regInfo{h: h}
		if msg != nil {
			pi.m = reflect.TypeOf(msg).Elem()
		}
		regProto[header] = pi
	case func(uint, interface{}):
		pi := &regInfo{ih: h}
		if msg != nil {
			pi.m = reflect.TypeOf(msg).Elem()
		}
		regProto[header] = pi
	}
}
