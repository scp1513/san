package handler

import (
	"github.com/scp1513/san/def"
	"github.com/scp1513/san/game/db"
	"github.com/scp1513/san/game/network"
	"github.com/scp1513/san/game/role"
	"github.com/scp1513/san/game/session"
	"github.com/scp1513/san/proto/req"
)

func init() {
	network.RegisterProto(&req.RoleData{}, handleRoleData)
}

func handleRoleData(sid uint64, msg interface{}) {
	sess := session.Get(sid)
	if sess == nil {
		return
	}
	m := msg.(*req.RoleData)
	r := role.Get(m.GetRoleID())
	if r != nil {
		if r.Data.Aid != sess.Aid {
			return
		}
		sess.RoleID = m.GetRoleID()
		network.SendClient(sid, r.ToProto())
	} else {
		db.FindRoleData(m.GetRoleID(), func(data *def.RoleData, err error) {
			sess := session.Get(sid)
			if sess == nil || sess.Aid != data.Aid {
				return
			}
			r := role.Cache(data)
			sess.RoleID = m.GetRoleID()
			network.SendClient(sid, r.ToProto())
		})
	}
}
