package common

import (
	"fmt"

	"github.com/scp1513/ec/util"
)

// 服务器登陆签名
func SrvLoginSign(t int64) string {
	return util.MD5(fmt.Sprintf("84#fa!;[4a5we\bfa9\r8w4r#%d", t))
}
