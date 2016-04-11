package opt

import (
	"encoding/json"
	"io/ioutil"
)

type _Val struct {
	RunMode     string `json:"run_mode"`
	SSLCertPath string `json:"ssl_cert_path"`
	SSLKeyPath  string `json:"ssl_key_path"`
	HTTPPort    int    `json:"http_port"`
	HTTPSPort   int    `json:"https_port"`
	DBHost      string `json:"db_host"`
	DBUser      string `json:"db_user"`
	DBPass      string `json:"db_pass"`
	DBName      string `json:"db_name"`
}

type _Json struct {
	Val _Val `json:"portal"`
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

func GetRunMode() string {
	return _json.Val.RunMode
}

func GetSSLCertPath() string {
	return _json.Val.SSLCertPath
}

func GetSSLKeyPath() string {
	return _json.Val.SSLKeyPath
}

func GetHTTPPort() int {
	return _json.Val.HTTPPort
}

func GetHTTPSPort() int {
	return _json.Val.HTTPSPort
}

func GetDBHost() string {
	return _json.Val.DBHost
}

func GetDBUser() string {
	return _json.Val.DBUser
}

func GetDBPass() string {
	return _json.Val.DBPass
}

func GetDBName() string {
	return _json.Val.DBName
}
