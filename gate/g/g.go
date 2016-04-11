// 全局数据
// 不允许import其他"san/***"内的package
package g

import (
	"github.com/scp1513/ec/serial"
)

var (
	// 函数串行化
	Serial = serial.New(4096)
)
