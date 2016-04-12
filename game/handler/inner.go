package handler

import (
	"github.com/golang/protobuf/proto"

	"github.com/scp1513/san/common"
	"github.com/scp1513/san/game/network"
	"github.com/scp1513/san/game/opt"
	"github.com/scp1513/san/game/session"
	"github.com/scp1513/san/proto/inner"
	"github.com/scp1513/san/stime"
)

func init() {
	network.RegisterProto(&inner.ReqSrvTime{}, handleReqSrvTime)
	network.RegisterProto(&inner.ReqSrvLogin{}, handleReqSrvLogin)
	network.RegisterProto(&inner.NfyCliDiscon{}, handleNfyCliDiscon)
}

func handleReqSrvTime(connID uint, msg interface{}) {
	m := msg.(*inner.ReqSrvTime)
	network.Send(connID, &inner.RspSrvTime{
		ReqTime: proto.Int64(m.GetReqTime()),
		SrvTime: proto.Int64(stime.Now().UnixNano()),
	})
}

func handleReqSrvLogin(connID uint, msg interface{}) {
	m := msg.(*inner.ReqSrvLogin)
	if stime.Now().Unix() > m.GetTime()+10 {
		network.Send(connID, &inner.RspSrvLogin{
			Success: proto.Bool(false),
			SrvID:   proto.Uint32(0),
		})
		return
	}

	sign := common.SrvLoginSign(m.GetTime())
	if m.GetSign() != sign {
		network.Send(connID, &inner.RspSrvLogin{
			Success: proto.Bool(false),
			SrvID:   proto.Uint32(0),
		})
		return
	}

	if err := network.OnServerLogin(connID, m.GetType()); err != nil {
		network.Send(connID, &inner.RspSrvLogin{
			Success: proto.Bool(false),
			SrvID:   proto.Uint32(0),
		})
		return
	}

	network.Send(connID, &inner.RspSrvLogin{
		Success: proto.Bool(true),
		SrvID:   proto.Uint32(uint32(connID)),
		Mode:    proto.String(opt.Val().Mode),
	})
}

func handleNfyCliDiscon(connID uint, msg interface{}) {
	m := msg.(*inner.NfyCliDiscon)
	session.Logout(m.GetSid())
}
