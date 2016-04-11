package network

import (
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/golang/protobuf/proto"

	"github.com/scp1513/ec/net"
	"github.com/scp1513/san/common"
	"github.com/scp1513/san/gate/g"
	"github.com/scp1513/san/gate/opt"
)

var (
	gameConnID uint   // 游戏服连接ID
	srvID      uint32 // 网关服务器ID
)

// 连接到游戏服
func connectGame() error {
	log.Println("connectGame")
	connectFinished := make(chan error, 1)
	connected := false
	for i := 0; i < 20; i++ {
		gameConnID = net.ConnectTo(map[string]interface{}{
			"ip":       opt.GetGameIP(),
			"port":     opt.GetGamePort(),
			"onStatus": func(connID uint, err error) { connectFinished <- err; onGameStatus(connID, err) },
			"onRecved": onGameRecved,
			"onClosed": onGameClosed,
		})
		if err := <-connectFinished; err != nil {
			gameConnID = 0
			log.Println("Reconnect game server in 3s.")
			time.Sleep(time.Second * 3)
			continue
		}
		connected = true
		break
	}
	if !connected {
		return fmt.Errorf("Connect game server failed!")
	}
	return nil
}

func disconnectGame() {
	log.Println("disconnectGame")
	net.Disconnect(gameConnID)
}

// 游戏服连接状态回调
func onGameStatus(connID uint, err error) {
	if err != nil {
		return
	}

	log.Println("Connect game success.")
	sendReqSrvLoginMsg()
}

// 游戏服接收数据回调
func onGameRecved(connID uint, data []byte) {
	if len(data) < 4 {
		log.Println("无效数据长度", len(data))
		return
	}
	header, _ := common.U32(data)
	if _, ok := regProto[header]; !ok {
		// 转发给客户端
		clientSendSerial.Post(func() { proxyToClientHandler(data) })
	} else {
		// 游戏服的消息
		recvSerial.Post(func() { recvGameHandler(data) })
	}
}

// 游戏服连接关闭回调
func onGameClosed(connID uint) {
	// 根据状态判断是否需要重连
	log.Println("onGameClosed")
	go func() {
		for status != statusStopping {
			status = statusReconnect
			log.Println("Reconnect game")
			if err := connectGame(); err == nil {
				break
			}
		}
	}()
}

// 主动关闭游戏服连接
func closeGameConn() {
}

func recvGameHandler(data []byte) {
	if len(data) < 4 {
		log.Println("数据长度无效:", len(data))
		return
	}
	header, _ := common.U32(data)
	v, ok := regProto[header]
	if !ok {
		log.Println("无效协议号:", header)
		return
	}
	var req interface{}
	if v.m != nil {
		req = reflect.New(v.m).Interface()
		err := proto.Unmarshal(data[4:], req.(proto.Message))
		if err != nil {
			log.Println("proto.Unmarshal 错误:", header)
			return
		}
	}
	g.Serial.Post(func() { v.ih(req) })
}
