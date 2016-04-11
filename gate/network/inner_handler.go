package network

import (
	"log"
	"time"

	"github.com/scp1513/san/gate/opt"
	"github.com/scp1513/san/mode"
	"github.com/scp1513/san/proto/inner"
	"github.com/scp1513/san/stime"
)

func init() {
	RegisterProto(&inner.RepSrvTime{}, handleRepSrvTime)
	RegisterProto(&inner.RepSrvLogin{}, handleRepSrvLogin)
	RegisterProto(&inner.NfyShutdown{}, handleNfyShutdown)
}

func handleRepSrvTime(msg interface{}) {
	m := msg.(*inner.RepSrvTime)
	t1 := m.GetReqTime()
	t2 := time.Now().UnixNano()
	srvTime := m.GetSrvTime()
	if ok, offset := stime.Delta(t1, t2, srvTime); ok {
		log.Println("Time offset =", offset)
		if err := startListenClient(); err != nil {
			panic(err)
		}
	} else {
		sendReqSrvTimeMsg()
	}
}

func handleRepSrvLogin(msg interface{}) {
	log.Println("handleRepSrvLogin")
	m := msg.(*inner.RepSrvLogin)
	if !m.GetSuccess() {
		sendReqSrvLoginMsg()
		return
	}
	srvID = m.GetSrvID()
	mode.Set(m.GetMode())
	if status == statusConnecting {
		if !opt.CheckSrvTime() {
			if err := startListenClient(); err != nil {
				panic(err)
			}
		} else {
			sendReqSrvTimeMsg()
		}
	} else if status == statusReconnect {
		status = statusRunning
	}
}

func handleNfyShutdown(msg interface{}) {
	m := msg.(*inner.NfyShutdown)
	sendShutdown <- m.GetFlag()
}
