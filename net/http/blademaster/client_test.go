package blademaster

import (
	"bytes"
	"context"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/iooikaak/frame/json"
)

func TestClientDoGatewayGetRequest(t *testing.T) {
	_gomsAddr = "10.180.18.20:8080"
	defer func() { _gomsAddr = "" }()
	c := NewClient(&ClientConfig{
		Dial:      time.Second,
		Timeout:   time.Second,
		KeepAlive: 100,
	})
	var resp interface{}
	param := url.Values{}
	param.Set("ids", "1")
	err := c.Get(context.Background(), "http://gateway.juqitech.com/v4/services/checkapi/status", "", param, &resp)
	if err != nil {
		t.Error(err)
	}
	t.Log(resp)
	_gomsAddr = ""
}

func TestClientDoGatewayPOSTRequest(t *testing.T) {
	_gomsAddr = "10.180.18.20:8080"
	defer func() { _gomsAddr = "" }()
	c := NewClient(&ClientConfig{
		Dial:      time.Second,
		Timeout:   time.Second,
		KeepAlive: 100,
	})
	var resp interface{}
	param := url.Values{}
	param.Set("click_url", "http://baidu.com")
	param.Set("supplier_url", "http://baidu.com")
	param.Set("item_id", "1")
	err := c.Post(context.Background(), "http://gateway.juqitech.com/v4/services/couponsapi/app/coupon", "", param, &resp)
	if err != nil {
		t.Error(err)
	}
	t.Log(resp)
	_gomsAddr = ""
}

func TestClientGetRequest(t *testing.T) {
	_gomsAddr = "10.180.18.20:8080"
	defer func() { _gomsAddr = "" }()
	c := NewClient(&ClientConfig{
		Dial:      time.Second,
		Timeout:   time.Second,
		KeepAlive: 100,
	})
	var resp interface{}
	err := c.Get(context.Background(), "http://www.juqitech.com/getIndexInfo?token=hellomf2018love", "", nil, &resp)
	if err != nil {
		t.Error(err)
	}
	t.Log(resp)
	_gomsAddr = ""
}

func TestClientGetRequestOnLocal(t *testing.T) {
	c := NewClient(&ClientConfig{
		Dial:      time.Second,
		Timeout:   time.Second,
		KeepAlive: 100,
	})
	var resp interface{}
	params := make(url.Values, 1)
	params.Add("access_token", "1")
	err := c.Get(context.Background(), "http://gateway.juqitech.com/v4/services/unknow/list", "", params, &resp)
	if err != nil {
		t.Error(err)
	}
	_, err = json.Marshal(&resp)
	if err != nil {
		t.Error("response is invalid")
	}
	t.Log(resp)
	_gomsAddr = ""
}

func TestClientOriginRequestOnExternal(t *testing.T) {
	_gomsAddr = "10.180.18.20:8080"
	c := NewClient(&ClientConfig{
		Dial:      time.Second,
		Timeout:   time.Second,
		KeepAlive: 100,
	})
	req, err := http.NewRequest(http.MethodGet, "https://www.juqitech.com", nil)
	if err != nil {
		t.Error(err)
	}
	_, err = c.Raw(context.Background(), req)
	if err != nil {
		t.Error(err)
	}
	_gomsAddr = ""
}

func TestClientOriginRequestOnInternal(t *testing.T) {
	_gomsAddr = "10.180.18.20:8080"
	c := NewClient(&ClientConfig{
		Dial:      time.Second,
		Timeout:   time.Second,
		KeepAlive: 100,
	})
	req, err := http.NewRequest(http.MethodGet, "http://www.juqitech.com/v4/services/demo/test/testy", nil)
	if err != nil {
		t.Error(err)
	}
	b, err := c.Raw(context.Background(), req)
	if err != nil {
		t.Error(err)
	}
	t.Log(string(b))
	_gomsAddr = ""
}

func TestClientOriginPostRequestOnInternal(t *testing.T) {
	_gomsAddr = "10.180.18.20:8080"
	defer func() { _gomsAddr = "" }()
	c := NewClient(&ClientConfig{
		Dial:      time.Second,
		Timeout:   time.Second,
		KeepAlive: 100,
	})

	param := url.Values{}
	param.Set("click_url", "http://baidu.com")
	param.Set("supplier_url", "http://baidu.com")
	param.Set("item_id", "1")

	req, err := http.NewRequest(http.MethodPost, "http://gateway.juqitech.com/v4/services/api/app/coupon", strings.NewReader(param.Encode()))
	if err != nil {
		t.Error(err)
	}

	b, err := c.Raw(context.Background(), req)
	if err != nil {
		t.Error(err)
	}

	t.Log(string(b))
	_gomsAddr = ""
}

func TestClientJson(t *testing.T) {
	url := "http://10.1.9.181:9012/endpoint/v1/notification_sign_list/find_by_id"
	c := NewClient(&ClientConfig{
		Dial:      time.Second,
		Timeout:   time.Second,
		KeepAlive: 100,
	})
	id := struct {
		Id string `json:"id"`
	}{
		"5e15708d7b42e30d4ce20283",
	}
	type Result struct {
	}
	type Resp struct {
		StatusCode int         `json:"statusCode"`
		ErrorCode  interface{} `json:"errorCode"`
		Comments   string      `json:"comments"`
		Result     Result      `json:"result"`
		Data       interface{} `json:"data"`
	}

	var resp Resp
	b, _ := json.Marshal(id)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(b))
	if err != nil {
		t.Error(err)
	}

	err = c.JSON(context.Background(), req, &resp)
	if err != nil {
		t.Error(err)
	}
	t.Log(resp)
	_gomsAddr = ""
}

func TestClientPostJson(t *testing.T) {
	url := "http://10.1.9.181:9012/endpoint/v1/notification_sign_list/find_by_id"
	c := NewClient(&ClientConfig{
		Dial:      time.Second,
		Timeout:   time.Second,
		KeepAlive: 100,
	})
	id := struct {
		Id string `json:"id"`
	}{
		"5e15708d7b42e30d4ce20283",
	}
	type Result struct {
	}
	type Resp struct {
		StatusCode int         `json:"statusCode"`
		ErrorCode  interface{} `json:"errorCode"`
		Comments   string      `json:"comments"`
		Result     Result      `json:"result"`
		Data       interface{} `json:"data"`
	}

	var resp Resp
	err := c.PostJson(context.Background(), url, "", id, &resp)
	if err != nil {
		t.Error(err)
	}
	t.Log(resp)
}
