package db

import (
	"fmt"

	"gopkg.in/mgo.v2/bson"

	"github.com/scp1513/san/def"
	"github.com/scp1513/san/game/g"
)

// CheckAccount 判断账号是否存在（Test模式使用）
func CheckAccount(acc, pwd string, callback func(uint32, error)) {
	g.GO(func() {
		var results map[string]interface{}
		err := mgoProxy.SelectOne(
			ACCOUNT_TEST_CNAME,
			bson.M{"account": acc},
			bson.M{"password": 0, "_id": 1},
			&results)
		if err != nil {
			g.Serial.Post(func() { callback(0, err) })
			return
		}

		if pwd != results["password"].(string) {
			g.Serial.Post(func() { callback(0, fmt.Errorf("密码错误")) })
			return
		}

		aid := uint32(results["_id"].(float64))
		g.Serial.Post(func() { callback(aid, nil) })
	})
}

// FindRoleList Find角色列表
func FindRoleList(sid uint64, aid uint32, callback func([]*def.SelectRoleInfo, error)) {
	g.GO(func() {
		var ret *def.AccountData
		err := mgoProxy.SelectByID(ACCOUNT_CNAME, aid, nil, &ret)
		if err != nil {
			g.Serial.Post(func() { callback(nil, err) })
			return
		}

		var roleInfo []*def.SelectRoleInfo
		err = mgoProxy.SelectAllWithParam(
			ROLE_CNAME,
			bson.M{"_id": bson.M{"$in": ret.RoleIDs}},
			"", nil, 0, 0,
			&roleInfo)
		g.Serial.Post(func() { callback(roleInfo, err) })
	})
}
