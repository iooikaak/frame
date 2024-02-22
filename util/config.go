package util

import (
	"encoding/json"
	"os"
)

func NewConfiguration(physicalPath string, configPointer interface{}) {
	file, err := os.Open(physicalPath)
	if err != nil {
		panic(err)
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(configPointer)
	if err != nil {
		panic(err)
	}
}
