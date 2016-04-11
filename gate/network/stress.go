package network

import (
	"crypto/tls"
	"fmt"
	"time"

	"github.com/scp1513/ec/net/http"
	"github.com/scp1513/san/common"
	"github.com/scp1513/san/gate/opt"
)

var (
	gateID           int
	isStressFeedback bool
)

func startFeedback() {
	queryGateID()
	isStressFeedback = true
	go func() {
		for isStressFeedback {
			time.Sleep(time.Minute)
			updateStress()
		}
	}()
}

func stopFeedback() {
	isStressFeedback = false
	releaseGate()
}

func queryGateID() {
	request := http.Post(common.GetIDUrl)
	request.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	request.SetTimeout(time.Second*3, time.Second*3)
	request.Param("ip", opt.GetIP())
	request.Param("port", fmt.Sprintf("%d", opt.GetPort()))

	type Resp struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		ID   int    `json:"id"`
	}
	var resp Resp
	err := request.ToJson(&resp)
	if err != nil {
		panic(err)
	}
	if resp.Code != 0 {
		panic(fmt.Errorf("%d: %s", resp.Code, resp.Msg))
	}

	gateID = resp.ID
}

func updateStress() {
	request := http.Post(common.StressUrl)
	request.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	request.SetTimeout(time.Second*3, time.Second*3)
	request.Param("id", fmt.Sprintf("%d", gateID))
	request.Param("stress", fmt.Sprintf("%d", GetClientNum()))

	type Resp struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}
	var resp Resp
	err := request.ToJson(&resp)
	if err != nil {
		panic(err)
	}
	if resp.Code != 0 {
		panic(fmt.Errorf("%d: %s", resp.Code, resp.Msg))
	}
}

func releaseGate() {
	request := http.Post(common.ReleaseUrl)
	request.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	request.SetTimeout(time.Second*3, time.Second*3)
	request.Param("id", fmt.Sprintf("%d", gateID))

	type Resp struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}
	var resp Resp
	err := request.ToJson(&resp)
	if err != nil {
		panic(err)
	}
	if resp.Code != 0 {
		panic(fmt.Errorf("%s", resp.Msg))
	}
}
