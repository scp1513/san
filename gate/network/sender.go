package network

import (
	"log"
	"time"

	"github.com/golang/protobuf/proto"

	"github.com/scp1513/ec/net"
	"github.com/scp1513/san/common"
	"github.com/scp1513/san/def"
	"github.com/scp1513/san/proto/inner"
)

//////////////////////////////////////////////////////////////////////////////
// to game

func doSendGame(data []byte) (n int, err error) {
	for i := 0; i < 60; i++ {
		n, err = net.Send(gameConnID, data)
		if err == nil {
			break
		}
		time.Sleep(time.Second)
	}
	return
}

func SendToGame(m _Message) {
	gameSendSerial.Post(func() {
		data, err := proto.Marshal(m)
		if err != nil {
			log.Println(err.Error())
			return
		}
		b := common.U32Byte(m.GetHeader())
		data = append(b, data...)
		doSendGame(data)
	})
}

func sendClientMsgToGame(connId uint, header uint32, data []byte) (n int, err error) {
	b := make([]byte, len(data)+12)
	copy(b, common.U32Byte(header))
	copy(b[4:], common.U64Byte(common.MakeU64(srvID, uint32(connId))))
	copy(b[12:], data)
	return doSendGame(b)
}

func sendReqSrvTimeMsg() {
	log.Println("sendReqSrvTimeMsg")
	req := &inner.ReqSrvTime{
		ReqTime: proto.Int64(time.Now().UnixNano()),
	}
	SendToGame(req)
}

func sendReqSrvLoginMsg() {
	log.Println("sendReqLoginMsg")
	now := time.Now().Unix()
	req := &inner.ReqSrvLogin{
		Type: proto.Uint32(def.SERVERTYPE_GATE),
		Sign: proto.String(common.SrvLoginSign(now)),
		Time: proto.Int64(now),
	}
	SendToGame(req)
}

//////////////////////////////////////////////////////////////////////////////
// to client

func proxyToClientHandler(data []byte) {
	if len(data) < 12 {
		log.Println("无效数据长度", len(data))
		return
	}
	connID, _ := common.U32(data[4:])
	clientConnMutex.RLock()
	c, ok := clientConn[uint(connID)]
	clientConnMutex.RUnlock()
	if ok {
		data = append(data[:4], data[12:]...)
		c.sender(data)
	}
}

func SendToClient(connID uint, m _Message) {
	clientConnMutex.RLock()
	c, ok := clientConn[connID]
	clientConnMutex.RUnlock()
	if ok {
		header := m.GetHeader()
		data, err := proto.Marshal(m)
		if err != nil {
			panic(err)
		}
		data = append(common.U32Byte(header), data...)
		c.sender(data)
	}
}
