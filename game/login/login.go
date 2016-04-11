package login

import (
	"crypto/tls"
	"fmt"
	"time"

	"github.com/golang/protobuf/proto"

	"github.com/scp1513/ec/net/http"
	"github.com/scp1513/san/common"
	"github.com/scp1513/san/game/g"
	"github.com/scp1513/san/game/network"
	"github.com/scp1513/san/game/session"
	"github.com/scp1513/san/proto/rsp"
)

func LoginVerify(token string, callback func(uint32, error)) {
	g.GO(func() {
		request := http.Post(common.LoginVerifyUrl)
		request.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
		request.SetTimeout(time.Second*3, time.Second*3)
		request.Param("token", token)

		var resp map[string]interface{}
		err := request.ToJson(&resp)
		if err != nil {
			panic(err)
		}
		code := int(resp["code"].(float64))
		if code != 0 {
			err = fmt.Errorf("%d: %s", code, resp["msg"].(string))
		}
		g.Serial.Post(func() {
			aid := uint32(0)
			if resp["aid"] != nil {
				aid = uint32(resp["aid"].(float64))
			}
			callback(aid, err)
		})
	})
}

func OnLoginVerifyFinished(sid uint64, aid uint32, err error) {
	if err != nil {
		network.SendClient(sid, &rsp.LoginVerify{
			Success:   proto.Bool(false),
			ErrMsg:    proto.String(err.Error()),
			AccountID: proto.Uint32(0),
		})
		return
	}

	sess, err := session.OnLoginSuccess(sid, aid)
	if err != nil {
		network.SendClient(sid, &rsp.LoginVerify{
			Success:   proto.Bool(false),
			ErrMsg:    proto.String(err.Error()),
			AccountID: proto.Uint32(0),
		})
		return
	}

	sess.Send(&rsp.LoginVerify{
		Success:   proto.Bool(true),
		AccountID: proto.Uint32(aid),
	})
}
