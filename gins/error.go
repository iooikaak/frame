package gins

import (
	"fmt"
	"strconv"
	"time"

	"github.com/iooikaak/frame/xlog"

	"github.com/iooikaak/frame/apiconstant"
)

// APIError API错误对象
type APIError struct {
	code apiconstant.ResponseType
	msg  string
	data interface{}
}

// Error 实现 error 接口
func (e *APIError) Error() string {
	return e.msg
}

// Code 获取错误码
func (e *APIError) Code() apiconstant.ResponseType {
	return e.code
}

// Data 获取数据
func (e *APIError) Data() interface{} {
	return e.data
}

// NewAPIError API错误，默认code：RESPONSE_ERROR
func NewAPIError(msg string, code ...apiconstant.ResponseType) *APIError {
	c := apiconstant.RESPONSE_ERROR
	if len(code) == 1 {
		c = code[0]
	}

	return &APIError{
		code: c,
		msg:  msg,
		data: struct{}{},
	}
}

// NewAPIErrorWithData 附带data信息的API错误，默认code：RESPONSE_ERROR
func NewAPIErrorWithData(msg string, data interface{}, code ...apiconstant.ResponseType) *APIError {
	c := apiconstant.RESPONSE_ERROR
	if len(code) == 1 {
		c = code[0]
	}

	if data == nil {
		data = struct{}{}
	}

	return &APIError{
		code: c,
		msg:  msg,
		data: data,
	}
}

// NewAPIErrorWithLog 创建API错误对象，转化为友好提示并记录日志，默认code：RESPONSE_CRASH
func NewAPIErrorWithLog(title, rawMsg string, code ...apiconstant.ResponseType) *APIError {
	errCode := strconv.FormatInt(time.Now().Unix(), 10)
	xlog.Errorf("%s[%s]：\n%s", title, errCode, rawMsg)

	data := make(map[string]string, 1)
	data["errCode"] = errCode

	c := apiconstant.RESPONSE_CRASH
	if len(code) == 1 {
		c = code[0]
	}

	return &APIError{
		code: c,
		msg:  fmt.Sprintf("请求发生错误[%s]，请联系技术人员", errCode),
		data: data,
	}
}
