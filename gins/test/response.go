package test

import (
	"encoding/json"
	"fmt"

	"github.com/iooikaak/frame/apiconstant"
)

type TestResponse struct {
	Code apiconstant.ResponseType `json:"code"`
	Msg  string                   `json:"msg"`
	Data interface{}              `json:"data"`
	body []byte
}

func NewTestResponse(body []byte) (*TestResponse, error) {
	r := &TestResponse{body: body}
	err := json.Unmarshal(body, r)
	if err != nil {
		r = nil
		err = fmt.Errorf("响应结果错误：%s，响应原文：%s", err.Error(), string(body))
	}

	return r, err
}

//响应结果
func (r *TestResponse) Content() string {
	return string(r.body)
}
