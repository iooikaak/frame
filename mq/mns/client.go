package mns

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	log "github.com/iooikaak/frame/log"

	"golang.org/x/time/rate"
)

const (
	mnsVersion     = "2015-06-06"
	defaultTimeout = int64(35)
)

// Client client
type Client struct {
	c *http.Client
	// 消费Handler
	handler Handler
	// 删除失败重试次数
	retry int

	// QPS
	qps int

	// 最大token数
	burst int

	// 控制消费协程退出
	ctx context.Context

	queueName string
	// 认证
	cre Credential
	// 限速
	limiter *rate.Limiter
}

// NewClient
func NewClient(ops ...ClientOptionFunc) *Client {
	client := new(Client)
	client.retry = 3
	client.ctx = context.Background()
	client.qps = 100
	client.burst = 5000
	for _, o := range ops {
		if err := o(client); err != nil {
			panic(err.Error())
		}
	}

	// 1000 毫秒 生成 client.cfg.QPS 个token
	r := rate.Every(time.Millisecond * 1000 / time.Duration(client.qps))
	client.limiter = rate.NewLimiter(r, client.burst)

	client.cre = NewCredential()
	// set request timeout
	transport := &http.Transport{
		Dial: func(network, addr string) (net.Conn, error) {
			c, err := net.DialTimeout(network, addr, time.Second*5)
			if err != nil {
				return nil, err
			}
			return c, nil
		},
		ResponseHeaderTimeout: time.Second * 35,
	}

	client.c = &http.Client{
		Transport: transport,
	}
	return client
}

// Send send mns request
func (client *Client) Send(
	method string, h http.Header,
	message interface{}, resource string, cfg QueueNode) (*http.Response, error) {
	// 限流器
	if err := client.limiter.Wait(context.Background()); err != nil {
		log.Error(err.Error())
	}
	var err error
	var body []byte

	if message == nil {
		body = nil
	} else {
		switch m := message.(type) {
		case []byte:
			body = m
		default:
			if body, err = xml.Marshal(message); err != nil {
				return nil, err
			}
		}
	}

	bodyMd5 := md5.Sum(body)
	bodyMd5Str := fmt.Sprintf("%x", bodyMd5)
	if h == nil {
		h = make(http.Header)
	}

	h.Add("x-mns-version", mnsVersion)
	h.Add("Content-Type", "application/xml")
	h.Add("Content-MD5", base64.StdEncoding.EncodeToString([]byte(bodyMd5Str)))
	h.Add("Date", time.Now().UTC().Format(http.TimeFormat))

	return client.sendMsg(h, method, resource, cfg, body)
}

func (client *Client) sendMsg(h http.Header, method, resource string,
	nodeCfg QueueNode, body []byte) (*http.Response, error) {

	signStr, err := client.cre.Sign(h, method, resource, nodeCfg.AccessSecret)
	if err != nil {
		return nil, err
	}

	authSignStr := fmt.Sprintf("MNS %s:%s", nodeCfg.AccessID, signStr)
	h.Add("Authorization", authSignStr)

	url := nodeCfg.Host + "/" + resource
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header = h

	// 重试 retry 次
	var dResp *http.Response
	for i := 0; i < client.retry; i++ {
		if dResp, err = client.c.Do(req); err == nil {
			if err = toErr(method, dResp); err == nil {
				return dResp, nil
			}
		}
	}
	return nil, err
}

// deleteMsg 删除消息
func (client *Client) deleteMsg(receiptHandle string, nodeCfg QueueNode) error {
	resource := fmt.Sprintf("queues/%s/%s?ReceiptHandle=%s", client.queueName, "messages", receiptHandle)
	_, err := client.Send("DELETE", nil, nil, resource, nodeCfg)
	return err
}

// recvMessages 批量接受消息接受消息，最大限制16个消息
func (client *Client) recvMessages(nodeCfg QueueNode) {
	if !nodeCfg.OK() {
		return
	}

	log.Infof("开始从节点[%s]地址[%s]接受消息", nodeCfg.Type, nodeCfg.Host)
	resource := fmt.Sprintf("queues/%s/%s?numOfMessages=%d&waitseconds=%d", client.queueName, "messages", 16, 30)
	for {
		sResp, err := client.Send("GET", nil, nil, resource, nodeCfg)
		if err != nil {
			if err != errMnsMessageNotExsit {
				client.handler.Error(err)
			}
			continue
		}

		var messages BatchMessageResp
		decoder := xml.NewDecoder(sResp.Body)
		if err := decoder.Decode(&messages); err != nil {
			client.handler.Error(err)
			continue
		}

		for _, message := range messages.Messages {
			if err := client.handler.HandleMsg(message); err == nil {
				if err := client.deleteMsg(message.ReceiptHandle, nodeCfg); err != nil {
					log.Error(err.Error())
				}
			}
		}

		select {
		case <-client.ctx.Done():
			log.Infof("停止从节点[%s]地址[%s]接受消息", nodeCfg.Type, nodeCfg.Host)
			return
		default:
		}
	}
}

func toErr(method string, resp *http.Response) error {
	switch strings.ToUpper(method) {
	case "GET":
		if resp.StatusCode == http.StatusOK {
			return nil
		}
	case "DELETE":
		if resp.StatusCode == http.StatusNoContent {
			return nil
		}
	case "POST":
		if resp.StatusCode == http.StatusCreated {
			return nil
		}
	default:
		return fmt.Errorf("bad method:%s", method)
	}

	decoder := xml.NewDecoder(resp.Body)
	var errMsg ErrorMessage
	if err := decoder.Decode(&errMsg); err != nil {
		return err
	}

	// GET 没有获取到消息
	if strings.ToUpper(method) == "GET" && strings.Contains(errMsg.Code, "MessageNotExist") {
		return errMnsMessageNotExsit
	}

	log.Errorf("METHOD:%s ERR:%s", method, errMsg.Error())
	return errMsg.Error()
}
