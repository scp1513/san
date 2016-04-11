package master

import (
	"crypto/tls"
	"log"
	"net"
	"net/http"

	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/netutil"
)

// TODO: 测试最大链接数
const (
	MAX_CONN_NUM = 3000
)

// copy from gin.
func debugPrint(format string, values ...interface{}) {
	if gin.IsDebugging() {
		log.Printf("[GIN-debug] "+format, values...)
	}
}

func debugPrintError(err error) {
	if err != nil {
		debugPrint("[ERROR] %v\n", err)
	}
}

// 1. endless热更新（所有request处理完才退出）
// 2. netutil最大链接数（chan）
func RunHttp(r *gin.Engine, addr string) (err error) {
	defer func() { debugPrintError(err) }()

	debugPrint("Listening and serving HTTP on %s\n", addr)
	err = endless.ListenAndServe(addr, r)

	return
}

func RunHttpTLS(r *gin.Engine, addr string, certFile string, keyFile string) (err error) {
	debugPrint("Listening and serving HTTPS on %s\n", addr)
	defer func() { debugPrintError(err) }()

	err = endless.ListenAndServeTLS(addr, certFile, keyFile, r)
	return
}

func RunHttp2(handler http.Handler, addr string) (err error) {
	defer func() { debugPrintError(err) }()
	debugPrint("Listening and serving HTTP on %s\n", addr)

	l, err := net.Listen("tcp", addr)
	if err != nil {
		return
	}

	l = netutil.LimitListener(TCPKeepAliveListener(l), MAX_CONN_NUM)
	server := &http.Server{Addr: addr, Handler: handler}
	//server := endless.NewServer(addr, handler)
	err = server.Serve(l)
	return
}

func RunHttpTLS2(handler http.Handler, addr string, certFile string, keyFile string) (err error) {
	debugPrint("Listening and serving HTTPS on %s\n", addr)
	defer func() { debugPrintError(err) }()

	certificate, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return err
	}

	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	server := &http.Server{Addr: addr, Handler: handler}
	config := &tls.Config{
		NextProtos:   []string{"http/1.1"},
		Certificates: []tls.Certificate{certificate},
	}

	l = tls.NewListener(TCPKeepAliveListener(l), config)
	l = netutil.LimitListener(l, MAX_CONN_NUM)

	err = server.Serve(l)
	return
}
