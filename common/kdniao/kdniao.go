package kdniao

import (
	"context"
	"net/http"
	"net/url"

	"github.com/iooikaak/frame/common/helper"
)

// 接口文档地址
// https://netmarket.oss.aliyuncs.com/4f693fa1-6c21-417b-a577-c31919f155b1.pdf?spm=5176.2020520132.101.4.36487218zXZC66&file=4f693fa1-6c21-417b-a577-c31919f155b1.pdf

const _tracingUrl = "http://www.cha.kdniao.com/getTracking"

type ExpressTracing struct {
	AppCode string
}

type Config struct {
	AppCode string `yaml:"appCode" json:"appCode"`
}

type GetTraceResp struct {
	StateEx      string   `json:"StateEx"`
	LogisticCode string   `json:"LogisticCode"`
	ShipperCode  string   `json:"ShipperCode"`
	Traces       []Traces `json:"Traces"`
	State        string   `json:"State"`
	NextCity     string   `json:"NextCity"`
	EBusinessID  string   `json:"EBusinessID"`
	Success      bool     `json:"Success"`
	Location     string   `json:"Location"`
}
type Traces struct {
	Action        string `json:"Action"`
	AcceptStation string `json:"AcceptStation"`
	AcceptTime    string `json:"AcceptTime"`
	Location      string `json:"Location"`
}

var defaultExpressTracing = &ExpressTracing{}

func Init(c *Config) {
	defaultExpressTracing.AppCode = c.AppCode
}

// 如果是JD， custInfo是京东的青龙配送编码，也叫商家编码，
// 如果是SF， custInfo是手机号后四位
// 其他为空
func GetTracing(ctx context.Context, shipperCode, logisticCode, custInfo string) (res *GetTraceResp, err error) {
	params := url.Values{
		"ShipperCode":  []string{shipperCode},
		"LogisticCode": []string{logisticCode},
		"CustInfo":     []string{custInfo},
	}
	res = &GetTraceResp{}
	req, err := helper.HTTPClient.NewRequest(http.MethodPost, _tracingUrl, "", params)
	if err != nil {
		return
	}
	req.Header.Set("Authorization", "APPCODE "+defaultExpressTracing.AppCode)
	err = helper.HTTPClient.Do(ctx, req, res)
	return
}
