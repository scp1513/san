package network

import (
	"github.com/scp1513/san/proto/inner"
	"log"
	"reflect"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"

	"github.com/scp1513/ec/net"
	"github.com/scp1513/san/common"
	"github.com/scp1513/san/gate/g"
	"github.com/scp1513/san/gate/opt"
	"github.com/scp1513/san/proto/nfy"
)

const (
	kTimeSecond      = 30
	kCheckLoopSecond = 10
	kCheckSize       = 64
)

type waitingConnInfo struct {
	Timeout int64
}

type clientConnInfo struct {
	connID uint
	key    string
	sender func([]byte)
}

var (
	clientTcpSrv     *net.TCPServer // 客户端连接
	waitingConn      = make(map[uint]*waitingConnInfo)
	clientConn       = make(map[uint]*clientConnInfo)
	waitingConnMutex sync.RWMutex
	clientConnMutex  sync.RWMutex
)

// startListenClient 开始监听客户端
func startListenClient() (err error) {
	log.Println("startListenClient")
	clientTcpSrv, err = net.NewTCPServer(map[string]interface{}{
		"ip":          opt.GetIP(),
		"port":        opt.GetPort(),
		"onConnected": onClientConnected,
		"onRecved":    onClientRecved,
		"onClosed":    onClientClosed,
	})
	if err != nil {
		return err
	}

	clientTcpSrv.Start()
	go checkWaitingConn()
	<-sendClientBegin
	return nil
}

// stopListenClient 停止监听客户端
func stopListenClient() error {
	log.Println("StopListenClient")
	clientTcpSrv.StopListen()
	return nil
}

// closeAllClient 关闭所有客户端连接
func closeAllClient() error {
	log.Println("closeAllClient")
	clientTcpSrv.CloseAll()
	return nil
}

// notifyClientShutdown 通知所有客户端停止服务器
func notifyClientShutdown() {
	log.Println("notifyClientShutdown")
	data, err := proto.Marshal(&nfy.Shutdown{})
	if err != nil {
		panic(err)
	}
	clientConnMutex.RLock()
	for _, v := range clientConn {
		v.sender(data)
	}
	clientConnMutex.RUnlock()
}

// 客户端连接成功回调
func onClientConnected(connID uint) {
	log.Println("connected", connID)
	waitingConnMutex.Lock()
	defer waitingConnMutex.Unlock()
	waitingConn[connID] = &waitingConnInfo{Timeout: time.Now().Unix() + kTimeSecond}
}

// 客户端接收数据回调
func onClientRecved(connID uint, data []byte) {
	if len(data) < 4 {
		log.Println("invalid len", len(data))
		return
	}
	header, _ := common.U32(data)

	if _, ok := regProto[header]; !ok {
		// 转发游戏服务器
		sendClientMsgToGame(connID, header, data[4:])
	} else {
		// 加入处理队列
		recvSerial.Post(func() { recvClientHandler(connID, data) })
	}
}

// 客户端连接关闭回调
func onClientClosed(connID uint) {
	log.Println("disconnected ", connID)
	waitingConnMutex.Lock()
	delete(waitingConn, connID)
	waitingConnMutex.Unlock()
	clientConnMutex.Lock()
	delete(clientConn, connID)
	clientConnMutex.Unlock()
	SendToGame(&inner.NfyCliDiscon{Sid: proto.Uint64(common.MakeU64(srvID, uint32(connID)))})
}

// 检查&断开超时链接
func checkWaitingConn() {
	for {
		var deleteWaitingConn [kCheckSize]uint
		waitingConnMutex.RLock()
		now := time.Now().Unix()
		size := 0
		for k, v := range waitingConn {
			if now >= v.Timeout {
				deleteWaitingConn[size] = k
				size++
				if size >= kCheckSize {
					break
				}
			}
		}
		waitingConnMutex.RUnlock()
		for i := 0; i < size; i++ {
			clientTcpSrv.Close(deleteWaitingConn[i])
		}
		time.Sleep(time.Second * kCheckLoopSecond)
	}
}

func recvClientHandler(connID uint, data []byte) {
	// TODO: 解密
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
	g.Serial.Post(func() { v.h(connID, req) })
}

func OnClientEncInfo(connID uint) string {
	waitingConnMutex.Lock()
	delete(waitingConn, connID)
	waitingConnMutex.Unlock()

	sender := func(data []byte) {
		// TODO: 加密
		clientTcpSrv.Send(connID, data)
	}
	key := "abcdefghijklmnopqrstuvwxyz"
	connInfo := &clientConnInfo{connID, key, sender}

	clientConnMutex.Lock()
	clientConn[connID] = connInfo
	clientConnMutex.Unlock()
	return key
}

// GetClientNum 获取客户端连接数
func GetClientNum() uint {
	return clientTcpSrv.GetConnCount()
}
