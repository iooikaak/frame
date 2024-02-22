package job

import (
	"encoding/json"
	"fmt"
)

type XxlJobConfig struct {
	Addresses   []string `json:"adds" yaml:"adds"`                 // []string{}
	AccessToken string   `json:"access_token" yaml:"access_token"` //
	AppName     string   `json:"app_name" yaml:"app_name"`         //执行器名称
	Port        int      `json:"port" yaml:"port"`                 //端口
}

func newConfig(configStr string) (conf *XxlJobConfig, err error) {

	conf = &XxlJobConfig{}
	err = json.Unmarshal([]byte(configStr), conf)
	if err != nil {
		return
	}

	if len(conf.Addresses) == 0 {
		err = fmt.Errorf("xxl-job: addresses cannot be empty")
	}

	if len(conf.AppName) == 0 {
		err = fmt.Errorf("xxl-job: app_name cannot be empty")
	}

	return
}
