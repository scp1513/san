package network

import (
	"reflect"
	"sync"

	"github.com/golang/protobuf/proto"

	"github.com/scp1513/ec/net"
	"github.com/scp1513/san/common"
	"github.com/scp1513/san/game/g"
	"github.com/scp1513/san/game/opt"
)

type SrvInfo struct {
	Type uint32
}

var (
	innerTCPSrv *net.TCPServer // 其他服务器连接
	connsMutext sync.RWMutex
	conns       = make(map[uint]*SrvInfo)
)

// 开始监听内部其他服务器
func startListenInner() (err error) {
	g.L.Println("Start listen inner.")
	innerTCPSrv, err = net.NewTCPServer(map[string]interface{}{
		"ip":          opt.Val().IP,
		"port":        opt.Val().Port,
		"onConnected": onSrvConnected,
		"onRecved":    onSrvRecved,
		"onClosed":    onSrvClosed,
	})
	if err != nil {
		return err
	}

	innerTCPSrv.Start()
	return nil
}

// 服务器连接成功回调
func onSrvConnected(connID uint) {
	if ip, err := innerTCPSrv.GetIP(connID); err != nil || !opt.IsIPValid(ip) {
		g.L.Println(ip)
		innerTCPSrv.Close(connID)
		return
	}
	connsMutext.Lock()
	defer connsMutext.Unlock()
	conns[connID] = &SrvInfo{}
}

// 服务器接收数据回调
func onSrvRecved(connID uint, data []byte) {
	recvSerial.Post(func() { recvHandler(connID, data) })
}

// 服务器连接关闭回调
func onSrvClosed(connID uint) {
	g.L.Println("Disconnect", connID)
	connsMutext.Lock()
	defer connsMutext.Unlock()
	g.L.Printf("Server disconnect %d, type: %d", connID, conns[connID].Type)
	delete(conns, connID)
}

// OnServerLogin 服务器登陆
func OnServerLogin(connID uint, stype uint32) error {
	g.L.Printf("Server login %d", stype)
	connsMutext.Lock()
	defer connsMutext.Unlock()
	conns[connID].Type = stype
	return nil
}

func sendHandler(connID uint, m Message) {
	data, err := proto.Marshal(m)
	if err != nil {
		g.L.Println(err.Error())
		return
	}
	b := common.U32Byte(m.GetHeader())
	innerTCPSrv.Send(connID, b, data)
}

func proxyHandler(sid uint64, m Message) {
	d, err := proto.Marshal(m)
	if err != nil {
		g.L.Println(err.Error())
		return
	}
	gateID, _ := common.SplitU64(sid)
	hb := common.U32Byte(m.GetHeader())
	sb := common.U64Byte(sid)
	innerTCPSrv.Send(uint(gateID), hb, sb, d)
}

func recvHandler(connID uint, data []byte) {
	if len(data) < 4 {
		g.L.Println("Invalid data len:", len(data))
		return
	}
	header, _ := common.U32(data)
	v, ok := regProto[header]
	if !ok {
		g.L.Println("Invalid header:", header)
		return
	}
	if v.ih != nil { // 服务器之间的协议
		var req interface{}
		if v.m != nil {
			req = reflect.New(v.m).Interface()
			err := proto.Unmarshal(data[4:], req.(proto.Message))
			if err != nil {
				g.L.Println("proto.Unmarshal error:", header, err.Error())
				return
			}
		}
		g.Serial.Post(func() { v.ih(connID, req) })
	} else {
		if len(data) < 12 {
			g.L.Println("Invlid client data len:", len(data))
			return
		}
		sid, _ := common.U64(data[4:])
		var req interface{}
		if v.m != nil {
			req = reflect.New(v.m).Interface()
			err := proto.Unmarshal(data[12:], req.(proto.Message))
			if err != nil {
				g.L.Println("proto.Unmarshal error:", header, err.Error())
				return
			}
		}
		g.Serial.Post(func() { v.h(sid, req) })
	}
}

// PostSrvProto 手动添加服务器消息
func PostSrvProto(connID uint, m Message) {
	g.Serial.Post(func() { regProto[m.GetHeader()].ih(connID, m) })
}

// PostClientProto 手动添加客户端消息
func PostClientProto(sid uint64, m Message) {
	g.Serial.Post(func() { regProto[m.GetHeader()].h(sid, m) })
}
