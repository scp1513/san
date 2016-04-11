package db

import (
	"github.com/scp1513/san/def"
	"github.com/scp1513/san/game/g"
)

func FindRoleData(rid uint64, callback func(*def.RoleData, error)) {
	g.GO(func() {
		var ret *def.RoleData
		err := mgoProxy.SelectByID(ROLE_CNAME, rid, nil, &ret)
		g.Serial.Post(func() { callback(ret, err) })
	})
}
