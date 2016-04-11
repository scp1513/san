package server

import (
	"net/url"

	"github.com/scp1513/san/portal/models/db"
)

func List(vals url.Values) map[string]interface{} {
	rows, err := db.Get().Query("SELECT realm_id, name, `desc`, flags FROM `server` GROUP BY realm_id")
	if err != nil {
		return map[string]interface{}{
			"code": 1,
			"msg":  err.Error(),
		}
	}
	type RealmInfo struct {
		RealmID int    `json:"realmID"`
		Name    string `json:"name"`
		Desc    string `json:"desc"`
		Flags   uint   `json:"flags"`
	}
	var realmInfos []*RealmInfo
	for rows.Next() {
		info := new(RealmInfo)
		rows.Scan(&info.RealmID, &info.Name, &info.Desc, &info.Flags)
		realmInfos = append(realmInfos, info)
	}
	return map[string]interface{}{
		"code": 0,
		"list": realmInfos,
	}
}

func GetID(vals url.Values) map[string]interface{} {
	id := uint32(0)
	err := db.Get().QueryRow("SELECT id FROM `server` WHERE ip = ? AND port = ?", vals.Get("ip"), vals.Get("port")).Scan(&id)
	if err != nil {
		return map[string]interface{}{
			"code": 1,
			"msg":  "无此服务器",
		}
	}
	db.Get().Exec("UPDATE `server` SET valid = 1 WHERE id = ?", id)
	return map[string]interface{}{
		"code": 0,
		"id":   id,
	}
}

func Release(vals url.Values) map[string]interface{} {
	_, err := db.Get().Exec("UPDATE `server` SET valid = 0 WHERE id = ?", vals.Get("id"))
	if err != nil {
		return map[string]interface{}{
			"code": 1,
			"msg":  "无此服务器: " + vals.Get("id"),
		}
	}
	return map[string]interface{}{
		"code": 0,
	}
}

func Stress(vals url.Values) map[string]interface{} {
	_, err := db.Get().Exec("UPDATE `server` SET stress = ? WHERE id = ?", vals.Get("stress"), vals.Get("id"))
	if err != nil {
		return map[string]interface{}{
			"code": 1,
			"msg":  "压力反馈失败 " + err.Error(),
		}
	}
	return map[string]interface{}{
		"code": 0,
	}
}

func GetAddr(vals url.Values) map[string]interface{} {
	rows, err := db.Get().Query("SELECT ip, port, stress FROM `server` WHERE realm_id = ?", vals.Get("realm_id"))
	if err != nil {
		return map[string]interface{}{
			"code": 1,
			"msg":  "获取服务器地址失败",
		}
	}

	type Info struct {
		IP   string
		Port int
	}
	var info Info
	minStress := uint32(0xFFFFffff)
	for rows.Next() {
		ip, port, stress := "", 0, uint32(0xFFFFffff)
		rows.Scan(&ip, &port, &stress)
		if minStress > stress {
			minStress = stress
			info.IP = ip
			info.Port = port
		}
	}
	if info.IP == "" {
		return map[string]interface{}{
			"code": 2,
			"msg":  "服务器已关闭",
		}
	}
	return map[string]interface{}{
		"code": 0,
		"ip":   info.IP,
		"port": info.Port,
	}
}
