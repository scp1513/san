package opt

import (
	"encoding/json"
	"io/ioutil"
	"strconv"
)

type _Val struct {
	IP               string   `json:"ip"`
	Port             int      `json:"port"`
	Mode             string   `json:"mode"`
	MgoHost          string   `json:"mgo_host"`
	MgoUser          string   `json:"mgo_user"`
	MgoPass          string   `json:"mgo_pass"`
	MgoDBName        string   `json:"mgo_dbname"`
	IPList           []string `json:"ip_list"`
	ResDir           string   `json:"res_dir"`
	DailyResetTime   string   `json:"daily_reset_time"`
	WeeklyResetTime  []string `json:"weekly_reset_time"`
	MonthlyResetTime []string `json:"monthly_reset_time"`
}

type _Json struct {
	Val _Val `json:"game"`
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

func Val() *_Val {
	return &_json.Val
}

func IsIPValid(ip string) bool {
	for _, v := range _json.Val.IPList {
		if ip == v {
			return true
		}
	}
	return false
}

func GetWeeklyResetTime() (int, string) {
	i, err := strconv.Atoi(_json.Val.WeeklyResetTime[0])
	if err != nil {
		i = 1
	}
	return i, _json.Val.WeeklyResetTime[1]
}

func GetMonthlyResetTime() (int, string) {
	i, err := strconv.Atoi(_json.Val.MonthlyResetTime[0])
	if err != nil {
		i = 1
	}
	return i, _json.Val.MonthlyResetTime[1]
}
