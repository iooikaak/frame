package gins

import (
	"context"
	"net/http"
	"time"

	"github.com/iooikaak/frame/apiconstant"
	"github.com/iooikaak/frame/ecode"
	"github.com/iooikaak/frame/gin"
	"github.com/iooikaak/frame/metadata"
	"github.com/iooikaak/frame/xlog"
)

// Context 请求上下文
type Context struct {
	*gin.Context // gin context
	C            context.Context
	// routerCtx    context.Context
	// routerCancel context.CancelFunc
	doneChan chan struct{}

	stack string
	isAPI bool
	API   api
	Web   web
}

// reset 重置Context
func (ctx *Context) reset(ginCtx *gin.Context) {
	ctx.Context = ginCtx
	ctx.C = ctx.Request.Context()
	// ctx.routerCtx, ctx.routerCancel = context.WithCancel(Instance.rootCtx)
	ctx.doneChan = make(chan struct{})

	ctx.stack = ""
	ctx.isAPI = false

	ctx.API.ctx = ctx
	ctx.API.result.Code = apiconstant.RESPONSE_UNKNOW
	ctx.API.result.Msg = ""
	ctx.API.result.Data = nil
	ctx.API.result.dataKV = nil
	ctx.API.rawResult = nil

	ctx.Web.ctx = ctx
}

// SetIsAPI 设置是否 API 请求标记
func (ctx *Context) SetIsAPI() {
	ctx.isAPI = true
}

// IsAPI 是否 API 请求
func (ctx *Context) IsAPI() bool {
	return ctx.isAPI
}

// setPanic 设置异常信息
func (ctx *Context) setPanic(stack string) {
	ctx.stack = stack
}

// Panic 异常堆栈信息
// gins.On500 里可获取
func (ctx *Context) Panic() (stack string) {
	return ctx.stack
}

// gins.JSON 劫持原有 JSON
func (ctx *Context) JSON(data interface{}, err error) {
	ctx.isAPI = true
	if err != nil {
		ctx.API.SetError(err)
	}
	ctx.API.SetData(data)
}

func (ctx *Context) Logger() *xlog.Entry {
	tid := ""
	if ctx.Context == nil {
		// cli or test mode
		return xlog.WithField(metadata.HttpTraceId, tid)
	}
	if t := ctx.Value(metadata.HttpTraceId); t != nil {
		tid = t.(string)
	}
	return xlog.WithField(metadata.HttpTraceId, tid)
}

// 如果非 Ready 状态为降级
func (ctx *Context) Ready() bool {
	return ctx.GetHeader(metadata.HttpCircuitBreaker) != "open"
}

// Value   ============================= 黑魔法高能预警, 原生gin.Context不支持非string的Value查找, 让gin可以拥有原生上下文的查找功能 ===================/
func (ctx *Context) Value(key interface{}) interface{} {
	if key == 0 {
		return ctx.Request
	}
	if keyAsString, ok := key.(string); ok {
		val, _ := ctx.Get(keyAsString)
		return val
	}
	return ctx.C.Value(key)
}

// Get returns the value for the given key, ie: (value, true).
// If the value does not exists it returns (nil, false)
func (c *Context) Get(key string) (value interface{}, exists bool) {
	c.KeysMutex.RLock()
	value, exists = c.Keys[key]
	c.KeysMutex.RUnlock()

	if exists {
		return
	}

	value = c.C.Value(key)
	return value, value != nil
}

// MustGet returns the value for the given key if it exists, otherwise it panics.
func (c *Context) MustGet(key string) interface{} {
	if value, exists := c.Get(key); exists {
		return value
	}
	if value := c.C.Value(key); value != nil {
		return value
	}
	panic("Key \"" + key + "\" does not exist")
}

// GetString returns the value associated with the key as a string.
func (c *Context) GetString(key string) (s string) {
	if val, ok := c.Get(key); ok && val != nil {
		s, _ = val.(string)
	}
	return
}

// GetBool returns the value associated with the key as a boolean.
func (c *Context) GetBool(key string) (b bool) {
	if val, ok := c.Get(key); ok && val != nil {
		b, _ = val.(bool)
	}
	return
}

// GetInt returns the value associated with the key as an integer.
func (c *Context) GetInt(key string) (i int) {
	if val, ok := c.Get(key); ok && val != nil {
		i, _ = val.(int)
	}
	return
}

// GetInt64 returns the value associated with the key as an integer.
func (c *Context) GetInt64(key string) (i64 int64) {
	if val, ok := c.Get(key); ok && val != nil {
		i64, _ = val.(int64)
	}
	return
}

// GetFloat64 returns the value associated with the key as a float64.
func (c *Context) GetFloat64(key string) (f64 float64) {
	if val, ok := c.Get(key); ok && val != nil {
		f64, _ = val.(float64)
	}
	return
}

// GetTime returns the value associated with the key as time.
func (c *Context) GetTime(key string) (t time.Time) {
	if val, ok := c.Get(key); ok && val != nil {
		t, _ = val.(time.Time)
	}
	return
}

// GetDuration returns the value associated with the key as a duration.
func (c *Context) GetDuration(key string) (d time.Duration) {
	if val, ok := c.Get(key); ok && val != nil {
		d, _ = val.(time.Duration)
	}
	return
}

// GetStringSlice returns the value associated with the key as a slice of strings.
func (c *Context) GetStringSlice(key string) (ss []string) {
	if val, ok := c.Get(key); ok && val != nil {
		ss, _ = val.([]string)
	}
	return
}

// GetStringMap returns the value associated with the key as a map of interfaces.
func (c *Context) GetStringMap(key string) (sm map[string]interface{}) {
	if val, ok := c.Get(key); ok && val != nil {
		sm, _ = val.(map[string]interface{})
	}
	return
}

// GetStringMapString returns the value associated with the key as a map of strings.
func (c *Context) GetStringMapString(key string) (sms map[string]string) {
	if val, ok := c.Get(key); ok && val != nil {
		sms, _ = val.(map[string]string)
	}
	return
}

// GetStringMapStringSlice returns the value associated with the key as a map to a slice of strings.
func (c *Context) GetStringMapStringSlice(key string) (smss map[string][]string) {
	if val, ok := c.Get(key); ok && val != nil {
		smss, _ = val.(map[string][]string)
	}
	return
}

// ============================= 黑魔法高能预警结束 ===================/

type api struct {
	ctx       *Context
	result    apiResult
	rawResult []byte
}

type apiResult struct {
	Code   apiconstant.ResponseType `json:"status"`
	Msg    string                   `json:"msg"`
	Data   interface{}              `json:"data"`
	dataKV map[string]interface{}
}

type web struct {
	ctx *Context
}

// SetError 设置错误信息
func (a *api) SetError(err error) {
	a.result.Msg = err.Error()

	if e, ok := err.(*APIError); ok {
		a.result.Code = e.code
		a.result.Data = e.data
		return
	}

	if e, ok := err.(ecode.Code); ok {
		a.result.Code = apiconstant.ResponseType(e.Code())
		a.result.Msg = e.Message()
		return
	}

	a.result.Code = apiconstant.RESPONSE_ERROR
}

// SetMsg 设置信息，code默认 RESPONSE_ERROR
func (a *api) SetMsg(msg string, code ...apiconstant.ResponseType) {
	a.result.Msg = msg
	if len(code) == 1 {
		a.result.Code = code[0]
		return
	}

	a.result.Code = apiconstant.RESPONSE_ERROR
}

// SetData 设置输出的model
func (a *api) SetData(data interface{}) {
	if a.result.Code == apiconstant.RESPONSE_UNKNOW {
		a.result.Code = apiconstant.RESPONSE_OK
	}
	a.result.Data = data
}

// SetDataKV 设置KV，会覆盖掉 SetData
func (a *api) SetDataKV(key string, value interface{}) {
	a.result.Code = apiconstant.RESPONSE_OK
	if a.result.dataKV == nil {
		a.result.dataKV = make(map[string]interface{})
	}

	a.result.dataKV[key] = value
}

// SetRawResult 设置原始内容输出，Content-Type为application/json，优先响应
func (a *api) SetRawResult(rawResult []byte) {
	a.rawResult = rawResult
}

func (a *api) json() {
	if a.rawResult != nil {
		a.ctx.Context.Data(http.StatusOK, "application/json", a.rawResult)
		return
	}

	if a.result.dataKV != nil {
		a.result.Data = a.result.dataKV
	}

	if a.result.Data == nil {
		a.result.Data = struct{}{}
	}

	a.ctx.Context.JSON(http.StatusOK, a.result)
}

// Render 立即渲染API
func (a *api) Render() {
	a.json()
}

// GetCode 设置输出的业务code
func (a *api) GetCode() int {
	return int(a.result.Code)
}

// Render 立即渲染Web
func (w *web) Render(name string, obj interface{}) {
	w.ctx.Context.HTML(http.StatusOK, name, obj)
}
