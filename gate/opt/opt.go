package opt

import (
	"encoding/json"
	"io/ioutil"
)

type _Val struct {
	IP           string `json:"ip"`
	Port         int    `json:"port"`
	GameIP       string `json:"game_ip"`
	GamePort     int    `json:"game_port"`
	CheckSrvTime bool   `json:"check_srv_time"`
}

type _Json struct {
	Val _Val `json:"gate"`
}

var _json _Json

func Load(fn string) error {
	data, err := ioutil.ReadFile(fn)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &_json)
	if err != nil {
		return err
	}
	return nil
}

func GetIP() string {
	return _json.Val.IP
}

func GetPort() int {
	return _json.Val.Port
}

func GetGameIP() string {
	return _json.Val.GameIP
}

func GetGamePort() int {
	return _json.Val.GamePort
}

func CheckSrvTime() bool {
	return _json.Val.CheckSrvTime
}
