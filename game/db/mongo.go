package db

import (
	"fmt"

	"github.com/scp1513/san/game/g"
	"github.com/scp1513/san/game/opt"
	"github.com/scp1513/san/mongo"
)

var mgoProxy mongo.Proxy

// Init 初始化
func Init() {
	g.L.Println("Init DB proxy.")
	err := mgoProxy.Dial(opt.Val().MgoHost, opt.Val().MgoUser, opt.Val().MgoPass, opt.Val().MgoDBName)
	if err != nil {
		//panic(err)
		fmt.Println(err.Error())
	}
}
