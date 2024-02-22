package json

import jsoniter "github.com/json-iterator/go"

var (
	json          = jsoniter.ConfigCompatibleWithStandardLibrary
	Marshal       = json.Marshal
	MarshalIndent = json.MarshalIndent
	NewDecoder    = json.NewDecoder
	Unmarshal     = json.Unmarshal
)
