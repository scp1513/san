// Package g 全局数据
package g

import (
	"fmt"
	"runtime"

	"github.com/scp1513/ec/serial"
)

// Serial 函数串行化
var Serial = serial.New(4096)

func init() {
	Serial.SetRecoverFunc(recoverFn)
}

func recoverFn(e interface{}) {
	buf := make([]byte, 4096)
	size := runtime.Stack(buf, false)
	fmt.Printf("%#v\n%s\n", e, buf[:size])
}
