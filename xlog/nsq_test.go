package xlog

import (
	"encoding/json"
	"testing"
)

func TestNewNsq(t *testing.T) {
	nsq, err := NewNsq(&NsqConfig{
		Addr:  "http://127.0.0.1:4161",
		Topic: "servicelog",
	})
	if err != nil {
		panic(err)
	}
	tmp := struct {
		Value string `json:"value"`
	}{
		Value: "hello 2020!",
	}
	bs, _ := json.Marshal(tmp)
	if _, err = nsq.Write(bs); err != nil {
		panic(err)
	}
}
