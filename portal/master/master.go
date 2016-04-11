package master

import (
	"fmt"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/scp1513/san/portal/api"
	"github.com/scp1513/san/portal/models/db"
	"github.com/scp1513/san/portal/opt"
)

// 创建pid文件
func createPid() {
	file, err := os.OpenFile("pid", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	file.WriteString(strconv.Itoa(os.Getpid()))
}

// 删除pid文件
func delPid() {
	os.Remove("pid")
}

// Init 初始化
func Init() {
	if err := opt.Load("res/opt.json"); err != nil {
		panic(err)
	}
	db.Init(opt.GetDBHost(), opt.GetDBUser(), opt.GetDBPass(), opt.GetDBName())
}

// Run 运行portal
func Run() {
	createPid()

	Init()

	// 开发模式
	gin.SetMode(opt.GetRunMode())

	// 注册路由
	router := gin.Default()
	api.SetupRouters(router)

	if gin.IsDebugging() {
		go func() {
			RunHttp(router, fmt.Sprintf(":%d", opt.GetHTTPPort()))
		}()
	}

	//go func() {
	RunHttpTLS(router, fmt.Sprintf(":%d", opt.GetHTTPSPort()), opt.GetSSLCertPath(), opt.GetSSLKeyPath())
	//RunHttpTLSDual(router, ":443", "res/ssl/server.crt", "res/ssl/server.key", "res/ssl/ca.crt")
	//}()
}
