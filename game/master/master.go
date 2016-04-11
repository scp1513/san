package master

import (
	"github.com/scp1513/san/game/cnf"
	"github.com/scp1513/san/game/db"
	"github.com/scp1513/san/game/g"
	"github.com/scp1513/san/game/network"
	"github.com/scp1513/san/game/opt"
	"github.com/scp1513/san/mode"

	_ "github.com/scp1513/san/game/handler"
)

func Run() {
	opt.Load("res/opt.json")
	if err := g.InitLogger(""); err != nil {
		panic(err)
	}
	mode.Set(opt.Val().Mode)
	cnf.LoadAll()
	db.Init()
	finished := make(chan struct{}, 0)
	g.GO(func() {
		g.L.Println("Global serial start.")
		g.Serial.Run()
		g.L.Println("Global serial finish.")
		<-finished
	})
	if err := network.Run(); err != nil {
		panic(err)
	}
	g.Serial.Stop()
	g.L.Println("Global serial stop.")
	finished <- struct{}{}
	g.ReleaseLogger()
}
