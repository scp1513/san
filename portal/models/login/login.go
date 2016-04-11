package login

import (
	"log"
	"net/url"
	"time"

	"github.com/scp1513/ec/util"
	"github.com/scp1513/san/portal/models/db"
	"github.com/scp1513/san/portal/models/def"
	"github.com/scp1513/san/portal/models/gen"
)

// Login 登陆
// account
// password
func Login(vals url.Values) map[string]interface{} {
	id, checkPwd, salt, token, tokenExp := uint32(0), "", "", "", int64(0)
	err := db.Get().QueryRow("SELECT id, password, salt, token, token_exp FROM `account` WHERE account = ?", vals.Get("account")).Scan(&id, &checkPwd, &salt, &token, &tokenExp)
	if err != nil {
		log.Println(err.Error())
		return map[string]interface{}{
			"code": 1,
			"msg":  "账号无效",
		}
	}
	if checkPwd != util.MD5(vals.Get("password"), salt) {
		return map[string]interface{}{
			"code": 2,
			"msg":  "密码错误",
		}
	}

	now := time.Now().Unix()
	if tokenExp < now {
		newtokenExp := now + def.TokenExpSecond
		for i := 0; i < 10; i++ {
			newtoken := gen.RandToken()
			_, err := db.Get().Exec("UPDATE `account` SET token = ?, token_exp = ? WHERE id = ?", newtoken, newtokenExp, id)
			if err == nil {
				token = newtoken
				tokenExp = newtokenExp
				break
			}
		}
	}

	return map[string]interface{}{
		"code":     0,
		"token":    token,
		"tokenExp": tokenExp - time.Now().Unix(),
	}
}

// LoginVerify 登陆验证
// token
func LoginVerify(vals url.Values) map[string]interface{} {
	id, tokenExp := uint32(0), int64(0)
	err := db.Get().QueryRow("SELECT id, token_exp FROM `account` WHERE token = ?", vals.Get("token")).Scan(&id, &tokenExp)
	if err != nil {
		return map[string]interface{}{
			"code": 1,
			"msg":  "无效token",
		}
	}
	now := time.Now().Unix()
	// 允许超时1分钟
	if now > tokenExp+60 {
		return map[string]interface{}{
			"code": 2,
			"msg":  "token过期",
		}
	}
	return map[string]interface{}{
		"code": 0,
		"aid":  id,
	}
}
