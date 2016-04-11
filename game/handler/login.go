package handler

import (
	"github.com/golang/protobuf/proto"

	"github.com/scp1513/san/def"
	"github.com/scp1513/san/game/db"
	"github.com/scp1513/san/game/login"
	"github.com/scp1513/san/game/network"
	"github.com/scp1513/san/game/session"
	"github.com/scp1513/san/mode"
	"github.com/scp1513/san/proto/req"
	"github.com/scp1513/san/proto/rsp"
)

func init() {
	network.RegisterProto(&req.LoginVerify{}, handleLoginVerify)
	network.RegisterProto(&req.RoleList{}, handleRoleList)
}

func handleLoginVerify(sid uint64, msg interface{}) {
	m := msg.(*req.LoginVerify)
	if mode.Val() == mode.Test {
		db.CheckAccount(m.GetAccount(), m.GetPassword(), func(aid uint32, err error) {
			login.OnLoginVerifyFinished(sid, aid, err)
		})
	} else {
		login.LoginVerify(m.GetToken(), func(aid uint32, err error) {
			login.OnLoginVerifyFinished(sid, aid, err)
		})
	}
}

func handleRoleList(sid uint64, msg interface{}) {
	sess := session.Get(sid)
	if sess == nil {
		return
	}

	db.FindRoleList(sid, sess.Aid, func(info []*def.SelectRoleInfo, err error) {
		sess := session.Get(sid)
		if sess == nil {
			return
		}
		if err != nil {
			sess.Send(&rsp.RoleList{
				Success: proto.Bool(false),
				ErrMsg:  proto.String(err.Error()),
			})
		}
		mrsp := &rsp.RoleList{Success: proto.Bool(true)}
		mrsp.RoleInfos = make([]*rsp.RoleList_RoleInfo, len(info))
		for i, v := range info {
			mrsp.RoleInfos[i] = &rsp.RoleList_RoleInfo{
				Rid:  proto.Uint64(v.Rid),
				Name: proto.String(v.Name),
			}
		}
		sess.Send(mrsp)
	})
}
