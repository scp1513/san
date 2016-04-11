package register

import (
	"net/url"
	"time"

	"github.com/scp1513/ec/util"
	"github.com/scp1513/san/portal/models/db"
	"github.com/scp1513/san/portal/models/def"
	"github.com/scp1513/san/portal/models/gen"
)

// Visitor 创建游客账号
func Visitor(vals url.Values) map[string]interface{} {
	var (
		account  = ""
		salt     = gen.RandSalt()
		srcPwd   = gen.RandPassword()
		password = util.MD5(srcPwd)
		pwdSalt  = util.MD5(password, salt)
		token    = ""
		tokenExp = time.Now().Unix() + def.TokenExpSecond
		flag     = 1
		err      error
	)

	for i := 0; i < 10; i++ {
		account = gen.RandAccount()
		token = gen.RandToken()
		_, err = db.Get().Exec("INSERT INTO `account`(account,password,salt,token,token_exp,flags) VALUES(?,?,?,?,?,?)", account, pwdSalt, salt, token, tokenExp, flag)
		if err == nil {
			break
		}
	}
	if err != nil {
		return map[string]interface{}{
			"code": 1,
			"msg":  "创建游客账号失败，请重新尝试",
		}
	}
	return map[string]interface{}{
		"code":      0,
		"msg":       "",
		"account":   account,
		"password":  srcPwd,
		"token":     token,
		"token_exp": tokenExp - time.Now().Unix(),
	}
}

// Upgrade 升级账号
// tmpAcc
// tmpPwd
// account
// password
func Upgrade(vals url.Values) map[string]interface{} {
	id, checkPwd, salt, flag := uint32(0), "", "", uint32(0)
	err := db.Get().QueryRow("SELECT id, password, salt, flags FROM `account` WHERE account = ?", vals.Get("tmpAcc")).Scan(&id, &checkPwd, &salt, &flag)
	if err != nil {
		return map[string]interface{}{
			"code": 1,
			"msg":  "无此账号",
		}
	}
	if flag != 1 {
		return map[string]interface{}{
			"code": 2,
			"msg":  "非游客账号不能升级",
		}
	}
	if checkPwd != util.MD5(vals.Get("tmpPwd"), salt) {
		return map[string]interface{}{
			"code": 3,
			"msg":  "旧密码错误",
		}
	}
	flag = 0
	pwdSalt := util.MD5(vals.Get("password"), salt)
	token := gen.RandToken()
	tokenExp := time.Now().Unix()
	_, err = db.Get().Exec("UPDATE `account` SET account = ?, password = ?, token = ?, token_exp = ?, flags = ? WHERE id = ?", vals.Get("account"), pwdSalt, token, tokenExp, flag, id)
	if err != nil {
		return map[string]interface{}{
			"code": 4,
			"msg":  "账号已存在",
		}
	}
	return map[string]interface{}{
		"code":      0,
		"token":     token,
		"token_exp": tokenExp - time.Now().Unix(),
	}
}

// Reg 注册
// account
// password
func Reg(vals url.Values) map[string]interface{} {
	var (
		salt     = gen.RandSalt()
		pwdSalt  = util.MD5(vals.Get("password"), salt)
		token    = ""
		tokenExp = time.Now().Unix()
		err      error
	)
	for i := 0; i < 2; i++ {
		token = gen.RandToken()
		_, err = db.Get().Exec("INSERT INTO `account`(account,password,salt,token,tokenExp) VALUES(?,?,?,?)", vals.Get("account"), pwdSalt, salt, token, tokenExp)
		if err == nil {
			break
		}
	}
	if err != nil {
		return map[string]interface{}{
			"code": 1,
			"msg":  "账号已存在",
		}
	}
	return map[string]interface{}{
		"code":      0,
		"token":     token,
		"token_exp": tokenExp - time.Now().Unix(),
	}
}
