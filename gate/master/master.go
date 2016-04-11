package master

import (
	"log"

	"github.com/scp1513/san/gate/g"
	"github.com/scp1513/san/gate/network"
	"github.com/scp1513/san/gate/opt"

	_ "github.com/scp1513/san/gate/handler"
)

func Run() {
	opt.Load("res/opt.json")
	globalSerialFinish := make(chan struct{}, 0)
	go func() {
		log.Println("Global serial start.")
		g.Serial.Run()
		log.Println("Global serial finish.")
		<-globalSerialFinish
	}()
	err := network.Run()
	if err != nil {
		panic(err)
	}
	g.Serial.Stop()
	log.Println("Global serial stop.")
	globalSerialFinish <- struct{}{}
}
