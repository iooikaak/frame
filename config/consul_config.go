package config

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"

	"encoding/json"

	"time"

	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/go-cleanhttp"
)

//ConsulConfig 基于Consul的配置对象
type ConsulConfig struct {
	BaseConfig
	defaultKey   string
	consulClient *api.Client
	confs        map[string]*consulConfiger
}

type consulConfiger struct {
	value              atomic.Value
	mu                 sync.Mutex
	err                error
	consulQueryOptions *api.QueryOptions
}

//NewConsulConfig 创建Consul配置对象
func NewConsulConfig(configStr string) (*ConsulConfig, error) {
	if configStr == "" {
		return nil, errors.New("configStr不能为空")
	}

	configJSON := &struct {
		Scheme     string
		Server     string
		Port       uint
		DefaultKey string
	}{}

	err := json.Unmarshal([]byte(configStr), configJSON)
	if err != nil {
		return nil, errors.New("configStr参数格式错误")
	}

	if configJSON.Server == "" {
		return nil, errors.New("Server参数不能为空")
	}

	if configJSON.Port == 0 {
		return nil, errors.New("Port参数格式错误")
	}

	if configJSON.Scheme == "" {
		configJSON.Scheme = "http"
	}

	consulConf := &api.Config{
		Address:   fmt.Sprintf("%s:%d", configJSON.Server, configJSON.Port),
		Scheme:    configJSON.Scheme,
		Transport: cleanhttp.DefaultPooledTransport(),
	}

	consulClient, err := api.NewClient(consulConf)
	if err != nil {
		return nil, err
	}

	return NewConsulConfigFast(consulClient, configJSON.DefaultKey)
}

//NewConsulConfigFast 创建Consul配置对象（使用指定的ConsulClient）
func NewConsulConfigFast(consulClient *api.Client, defaultKey ...string) (*ConsulConfig, error) {
	if consulClient == nil {
		return nil, errors.New("ConsulClient对象不能为空")
	}

	defaultKeyLen := len(defaultKey)
	if defaultKeyLen > 1 {
		return nil, errors.New("DefaultKey参数只能设置一个")
	}

	c := &ConsulConfig{
		consulClient: consulClient,
		confs:        make(map[string]*consulConfiger),
	}
	c.typ = CONFIG_TYPE_CONSUL

	if defaultKeyLen == 1 {
		c.defaultKey = defaultKey[0]
	}

	return c, nil
}

//Body 返回DefaultKey的配置内容
func (c *ConsulConfig) Body() ([]byte, error) {
	if c.defaultKey == "" {
		return nil, errors.New("没有设置DefaultKey")
	}

	return c.Get(c.defaultKey)
}

func (c *ConsulConfig) getConf(key string) (conf *consulConfiger, has bool) {
	conf, has = c.confs[key]
	return
}

func (c *ConsulConfig) setConf(key string) {
	_, has := c.getConf(key)
	if has {
		return
	}

	c.confs[key] = &consulConfiger{
		// value:        atomic.Value{},
		// mu:           sync.Mutex{},
		consulQueryOptions: &api.QueryOptions{RequireConsistent: true, WaitIndex: 0},
	}
}

//Get 获取指定的配置内容
func (c *ConsulConfig) Get(key string) (body []byte, err error) {
	conf, has := c.getConf(key)
	if has {
		//阻塞确保Watch事件已经开始
		for {
			if conf.err != nil {
				err = conf.err
				return
			}

			if conf.consulQueryOptions.WaitIndex > 0 {
				//有数据
				break
			}

			//无数据，首次运行，阻塞等待，每10微秒一次
			time.Sleep(10 * time.Millisecond)
		}

		val := conf.value.Load()
		if val == nil {
			err = fmt.Errorf("[%s]配置不存在", key)
			return
		}

		body = val.([]byte)
	} else {
		body, err = c.getKV(key, &api.QueryOptions{RequireConsistent: true, WaitIndex: 0})
	}

	return
}

func (c *ConsulConfig) set(key string, body []byte, err error) {
	conf, has := c.getConf(key)
	if !has {
		return
	}

	if err != nil {
		conf.err = err
		return
	}

	if len(body) == 0 {
		return
	}

	conf.mu.Lock()
	conf.value.Store(body)
	conf.mu.Unlock()
}

//Watch 监控配置变化，后台执行
func (c *ConsulConfig) Watch(keys ...string) {
	for i, l := 0, len(keys); i < l; i++ {
		key := strings.Trim(keys[i], " ")
		c.setConf(key)
	}

	for k, v := range c.confs {
		go func(key string, conf *consulConfiger) {
			for {
				body, err := c.getKV(key, conf.consulQueryOptions)
				c.set(key, body, err)
				if err != nil {
					break
				}
			}
		}(k, v)
	}
}

func (c *ConsulConfig) getKV(key string, consulQueryOptions *api.QueryOptions) (body []byte, err error) {
	kv, meta, err := c.consulClient.KV().Get(key, consulQueryOptions)
	if err != nil {
		err = fmt.Errorf("检查[%s]配置出错：%s", key, err)
		return
	}

	if kv == nil {
		err = fmt.Errorf("[%s]配置不存在", key)
		return
	}

	body = kv.Value
	consulQueryOptions.WaitIndex = meta.LastIndex
	return
}
