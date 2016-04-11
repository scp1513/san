package handler

import (
	"log"

	"github.com/scp1513/san/gate/network"
	"github.com/scp1513/san/proto/req"
	"github.com/scp1513/san/proto/rsp"
)

func init() {
	network.RegisterProto(&req.EncInfo{}, handleEncInfo)
}

func handleEncInfo(connID uint, msg interface{}) {
	log.Println("handleEncInfo")
	key := network.OnClientEncInfo(connID)
	mr := &rsp.EncInfo{Key: &key}
	network.SendToClient(connID, mr)
}
