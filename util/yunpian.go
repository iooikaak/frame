package util

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// 云片短信

type YunPianClient struct {
	ApiKey     string
	HttpClient *http.Client
}

type SMSResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

const (
	DEFAULT_API_URL = "https://sms.yunpian.com/v2/sms/single_send.json"
)

var (
	HttpClientTimeout = 60 * time.Second
)

func NewYunPianClient(apiKey string) *YunPianClient {
	return &YunPianClient{
		ApiKey: apiKey,
		HttpClient: &http.Client{
			Timeout: HttpClientTimeout,
		},
	}
}

func (c *YunPianClient) Send(phone, content string) error {
	if phone == "" || content == "" {
		return fmt.Errorf("手机号或短信内容为空")
	}
	params := url.Values{
		"apikey": []string{c.ApiKey},
		"text":   []string{content},
		"mobile": []string{phone},
	}
	req, _ := http.NewRequest("POST", DEFAULT_API_URL, strings.NewReader(params.Encode()))
	req.Header.Set("Accept", "application/json;charset=utf-8;")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=utf-8;")
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		var response SMSResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		if err != nil {
			return err
		}
		if response.Code == 0 {
			return nil
		}
		return fmt.Errorf("发送失败：%s", response.Msg)
	}
	return fmt.Errorf("发送失败：%d", resp.StatusCode)
}
