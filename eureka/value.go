package eureka

import (
	"encoding/json"

	"github.com/iooikaak/frame/utils"
)

type Value struct {
	Val []*Instance
	Md5 string
}

func NewValue(val []*Instance) *Value {
	for _, v := range val {
		v.LastDirtyTimestamp = ""
		v.LastUpdatedTimestamp = ""
		v.IsCoordinatingDiscoveryServer = ""
	}
	b, _ := json.Marshal(val)
	md5 := utils.Md5Encrypt(string(b))
	return &Value{
		Val: val,
		Md5: md5,
	}
}
