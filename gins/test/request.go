package test

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
)

type TestRequest struct {
	// svr  *httptest.Server
	pm     *TestParam
	json   *TestJSON
	header http.Header
	isJSON bool
}

func NewTestRequest() *TestRequest {
	return &TestRequest{
		// svr: httptest.NewServer(http.HandlerFunc(engine.ServeHTTP)),
		pm:   NewTestParam(),
		json: NewTestJSON(),
	}
}

func (c *TestRequest) SetJSON() {
	c.isJSON = true
}

func (c *TestRequest) AddParam(key, value string) {
	c.pm.Add(key, value)
}

func (c *TestRequest) AddJSON(key string, value interface{}) {
	c.json.Add(key, value)
}

func (c *TestRequest) SetHeader(key string, value string) {
	c.header.Set(key, value)
}

func (c *TestRequest) AddHeader(key string, value string) {
	c.header.Add(key, value)
}

func (c *TestRequest) Call(method, url string) (*TestResponse, error) {
	method = strings.ToUpper(method)
	var (
		body []byte
		err  error
	)

	if !c.isJSON {
		if c.pm != nil {
			if method == "GET" {
				if strings.Contains(url, "?") {
					url += "&"
				} else {
					url += "?"
				}
				url += c.pm.Encode()
			} else {
				body = []byte(c.pm.Encode())
			}
		}
	} else {
		body = c.json.Body()
	}

	var req *http.Request
	if body != nil {
		req, err = http.NewRequest(method, url, bytes.NewReader(body))
	} else {
		req, err = http.NewRequest(method, url, nil)
	}

	if len(c.header) > 0 {
		req.Header = c.header
	}
	if c.isJSON {
		req.Header.Set("Content-Type", "application/json")
	} else {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	w := httptest.NewRecorder()
	Engine.ServeHTTP(w, req)

	if err != nil {
		err = errors.New("请求错误：" + err.Error())
		return nil, err
	}

	if w.Result().StatusCode != 200 {
		return nil, fmt.Errorf("请求错误：%d", w.Result().StatusCode)
	}

	// return getResponse(w.Result().Body)
	return NewTestResponse(w.Body.Bytes())
}

func getResponse(r io.ReadCloser) (*TestResponse, error) {
	body, err := ioutil.ReadAll(r)
	r.Close()
	if err != nil {
		err = errors.New("读取结果错误：" + err.Error())
		return nil, err
	}

	return NewTestResponse(body)
}
