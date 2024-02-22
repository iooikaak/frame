package core

import (
	goctx "context"
	"encoding/json"
	"encoding/xml"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	log "github.com/iooikaak/frame/log"
	"github.com/iooikaak/frame/protocol"

	"github.com/golang/protobuf/proto"
)

// Context 请求上下文
// Context包含了请求的所有信息，并封装了一系列所需的操作
type Context interface {
	// GetClientInfo 根据请求头Client-Info获取appID，appVersion，platform
	GetClientInfo() (string, int64, int64)

	// CurrentUserID
	CurrentUserID() int64

	// Context
	Ctx() goctx.Context

	// Bizid 服务全局ID
	Bizid() string

	// Reset Reset
	Reset(r *protocol.Proto, w *protocol.Proto)

	// Request 原始请求信息
	Request() *protocol.Proto

	// Response 响应信息
	Response() *protocol.Proto

	// Bind Parse Request data
	Bind(i interface{}) error

	// HTTP Header
	Header() http.Header

	// 请求数据序列化格式
	ReqFormat() protocol.RestfulFormat

	// 响应数据序列化格式
	RespFormat() protocol.RestfulFormat

	// Client Request RealIP
	RealIP() string

	// http 表单数据和raw query都使用此结构获取看k,v对
	FormValue(name string) string

	// FormValues FormValues
	FormValues() url.Values

	// GetString 获取FormValue的值
	GetString(name string, defaultValue string) string

	// GetInt 获取FormValue的值
	GetInt(name string, defaultValue int) int

	GetInt64(name string, defaultValue int64) int64

	// GetUint 获取FormValue值并转化为 uint64
	GetUint(name string, defaultValue uint64) uint64

	// GetFloat 获取FormValue值并转化为 float64
	GetFloat(name string, defaultValue float64) float64

	// GetBool 获取FormValue值并转化为 bool
	GetBool(name string, defaultValue bool) bool

	// Ctx Get
	Get(key string) interface{}

	// Ctx Set
	Set(key string, val interface{})

	// JSON 响应JSON数据
	JSON(i interface{}) error

	// XML 响应XML数据
	XML(i interface{}) error

	// PROTOBUF 响应PROTOBUF数据
	Protobuf(i proto.Message) error

	// Bytes 响应数据
	Bytes(i []byte, format protocol.RestfulFormat) error

	// ImagePng 响应PNG图片
	ImagePng(i []byte) error

	// String 响应数据
	String(s string) error

	// JSON2 返回code和msg的json数据格式
	JSON2(code int, msg string, data interface{}) error

	RetBadRequestError(msg string) error
	RetForbiddenError(msg string) error
	RetNotFoundError(msg string) error
	RetInternalError(msg string) error
	RetDisplayError(msg string) error
	RetCustomError(code int, msg string) error
	RetCustomErrorWithData(code int, msg string, data interface{}) error
	RetSuccess(msg string, data interface{}) error
	CrossDomain()

	// XML2 返回code和msg的xml数据格式
	XML2(code int, msg string, data interface{}) error

	HTTP(code int, msg string) error

	// GetHeaderString 从请求头中获取键值对并转换为string
	GetHeaderString(key string, defaultValue string) string

	// GetHeaderInt64 从请求头中获取键值对并转换为int64
	GetHeaderInt64(key string, defaultValue int64) int64

	// GetHeaderFloat64 从请求头中获取键值对并转换为float64
	GetHeaderFloat64(key string, defaultValue float64) float64

	Debug(args ...interface{})

	Warn(args ...interface{})

	Info(args ...interface{})

	Error(args ...interface{})

	Fatal(args ...interface{})

	Debugf(fmt string, args ...interface{})

	Warnf(fmt string, args ...interface{})

	Infof(fmt string, args ...interface{})

	Errorf(fmt string, args ...interface{})

	Fatalf(fmt string, args ...interface{})
}

// newContext new context
func newContext() Context {
	return &icecontext{
		req:    nil,
		resp:   nil,
		header: nil,
		form:   nil}
}

type icecontext struct {
	req       *protocol.Proto
	resp      *protocol.Proto
	header    http.Header
	srcFormat protocol.RestfulFormat
	dstFormat protocol.RestfulFormat
	form      url.Values
	clientip  string
	ctx       goctx.Context
}

// CurrentUserID 获取当前用户ID，仅针对APP的api有效
func (c *icecontext) CurrentUserID() int64 {
	uid, _ := strconv.ParseInt(c.Header().Get("userId"), 10, 64)
	return uid
}

// Ctx 将一些需要的参数传递给下一个请求的Context
// 多层级调用链，使用go context将bizid传递下去
func (c *icecontext) Ctx() goctx.Context {
	if c.ctx == nil {
		c.ctx = goctx.TODO()
	}
	return goctx.WithValue(c.ctx, "bizid", c.Bizid())
}

// Bizid 用户追踪ID
func (c *icecontext) Bizid() string {
	if c.req != nil {
		return c.req.GetBizid()
	}
	return ""
}

// Reset 调用其他方法前已在框架中Reset，故其他方法获取参数是安全的
func (c *icecontext) Reset(r *protocol.Proto, w *protocol.Proto) {
	c.req = r
	c.resp = w
	c.srcFormat = r.GetFormat()
	c.dstFormat = protocol.RestfulFormat_FORMATNULL

	c.header = make(http.Header)
	for k, v := range r.GetHeader() {
		c.header.Set(k, v)
	}

	c.form = make(url.Values)
	for k, v := range r.GetForm() {
		c.form.Set(k, v)
	}

	c.clientip = ""
	c.ctx = goctx.TODO()
}

// Header HTTP header
func (c *icecontext) Header() http.Header {
	if c.header == nil {
		c.header = make(http.Header)
		for k, v := range c.req.GetHeader() {
			c.header.Set(k, v)
		}
	}
	return c.header
}

// Request 原始请求信息
func (c *icecontext) Request() *protocol.Proto {
	return c.req
}

// Response 响应信息
func (c *icecontext) Response() *protocol.Proto {
	return c.resp
}

// Bind Parse Request data
func (c *icecontext) Bind(i interface{}) error {
	return protocol.Unpack(c.Request().GetFormat(),
		c.Request().GetBody(), i)
}

// ReqFormat 请求数据序列化格式
func (c *icecontext) ReqFormat() protocol.RestfulFormat {
	if c.req != nil {
		return c.req.GetFormat()
	}
	return protocol.RestfulFormat_FORMATNULL
}

// RespFormat 响应数据序列化格式
func (c *icecontext) RespFormat() protocol.RestfulFormat {
	if c.resp != nil {
		return c.resp.GetFormat()
	}
	return protocol.RestfulFormat_FORMATNULL
}

// URL RAW Query data
func (c *icecontext) FormValue(name string) string {
	return c.form.Get(name)
}

// GetString 获取FormValue中的参数值
func (c *icecontext) GetString(name string, defaultValue string) string {
	if v := c.form.Get(name); v != "" {
		return v
	}
	return defaultValue
}

// GetInt 获取FormValue中的参数值
func (c *icecontext) GetInt(name string, defaultValue int) int {
	if v := c.form.Get(name); v != "" {
		number, err := strconv.Atoi(v)
		if err != nil {
			return defaultValue
		}
		return number
	}
	return defaultValue
}

func (c *icecontext) GetInt64(name string, defaultValue int64) int64 {
	if v := c.form.Get(name); v != "" {
		number, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return defaultValue
		}
		return number
	}
	return defaultValue
}

// GetUint 获取FormValue中的参数值
func (c *icecontext) GetUint(name string, defaultValue uint64) uint64 {
	if v := c.form.Get(name); v != "" {
		number, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			return defaultValue
		}
		return number
	}
	return defaultValue
}

// GetFloat 获取FormValue中的参数值
func (c *icecontext) GetFloat(name string, defaultValue float64) float64 {
	if v := c.form.Get(name); v != "" {
		number, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return defaultValue
		}
		return number
	}
	return defaultValue
}

// GetBool 获取FormValue中的参数值
func (c *icecontext) GetBool(name string, defaultValue bool) bool {
	if v := c.form.Get(name); v != "" {
		if strings.ToLower(v) == "true" {
			return true
		}
		return false
	}
	return defaultValue
}

// FormValues FormValues
func (c *icecontext) FormValues() url.Values {
	return c.form
}

// RealIP Client Request RealIP
func (c *icecontext) RealIP() string {
	if c.clientip == "" {
		ra := c.Request().GetRemoteAddr()
		if ip := c.Header().Get(protocol.HeaderXForwardedFor); ip != "" {
			ra = strings.Split(ip, ", ")[0]
		} else if ip := c.Header().Get(protocol.HeaderXRealIP); ip != "" {
			ra = ip
		} else {
			ra, _, _ = net.SplitHostPort(ra)
		}
		c.clientip = ra
		return ra
	}
	return c.clientip
}

// Get Ctx Get
func (c *icecontext) Get(key string) interface{} {
	if c.ctx == nil {
		return nil
	}
	return c.ctx.Value(key)
}

// Set Ctx Set
func (c *icecontext) Set(key string, val interface{}) {
	if c.ctx == nil {
		c.ctx = goctx.TODO()
	}
	goctx.WithValue(c.ctx, key, val)
}

// JSON 响应JSON数据
func (c *icecontext) JSON(i interface{}) error {
	c.dstFormat = protocol.RestfulFormat_JSON
	b, err := json.Marshal(i)
	if err != nil {
		return err
	}
	if c.resp == nil {
		s := c.Request().Shadow()
		c.resp = &s
	}
	c.resp.Format = protocol.RestfulFormat_JSON
	c.resp.Body = b
	c.resp.SetHeader(protocol.HeaderContentType, protocol.MIMEApplicationJSONCharsetUTF8)
	return nil
}

// XML 响应XML数据
func (c *icecontext) XML(i interface{}) error {
	c.dstFormat = protocol.RestfulFormat_XML
	b, err := xml.Marshal(i)
	if err != nil {
		return err
	}
	if c.resp == nil {
		s := c.Request().Shadow()
		c.resp = &s
	}
	c.resp.Format = protocol.RestfulFormat_XML
	c.resp.Body = b
	c.resp.SetHeader(protocol.HeaderContentType, protocol.MIMEApplicationXML)
	return nil
}

// PROTOBUF 响应PROTOBUF数据
func (c *icecontext) Protobuf(i proto.Message) error {
	c.dstFormat = protocol.RestfulFormat_PROTOBUF
	b, err := proto.Marshal(i)
	if err != nil {
		return err
	}
	if c.resp == nil {
		s := c.Request().Shadow()
		c.resp = &s
	}
	c.resp.Format = protocol.RestfulFormat_PROTOBUF
	c.resp.Body = b
	c.resp.SetHeader(protocol.HeaderContentType, protocol.MIMEApplicationProtobuf)
	return nil
}

// Bytes 响应Bytes数据
func (c *icecontext) Bytes(i []byte, format protocol.RestfulFormat) error {
	c.dstFormat = format
	if c.resp == nil {
		s := c.Request().Shadow()
		c.resp = &s
	}
	c.resp.Format = format
	c.resp.Body = i
	switch format {
	case protocol.RestfulFormat_JSON:
		c.resp.SetHeader(protocol.HeaderContentType, protocol.MIMEApplicationJSONCharsetUTF8)
	case protocol.RestfulFormat_PROTOBUF:
		c.resp.SetHeader(protocol.HeaderContentType, protocol.MIMEApplicationProtobuf)
	case protocol.RestfulFormat_XML:
		c.resp.SetHeader(protocol.HeaderContentType, protocol.MIMETextXMLCharsetUTF8)
	default:
		c.resp.SetHeader(protocol.HeaderContentType, protocol.MIMETextPlainCharsetUTF8)
	}
	return nil
}

// ImagePng 响应Png数据
func (c *icecontext) ImagePng(i []byte) error {
	c.resp.Body = i
	c.resp.SetHeader(protocol.HeaderContentType, protocol.MIMEImagePNG)
	return nil
}

// String 响应数据
func (c *icecontext) String(s string) error {
	c.dstFormat = protocol.RestfulFormat_RAWQUERY
	if c.resp == nil {
		s := c.Request().Shadow()
		c.resp = &s
	}
	c.resp.Format = protocol.RestfulFormat_RAWQUERY
	c.resp.Body = []byte(s)
	c.resp.SetHeader(protocol.HeaderContentType, protocol.MIMETextPlainCharsetUTF8)
	return nil
}

// JSON2 公共JSON响应
func (c *icecontext) JSON2(code int, msg string, data interface{}) error {
	return c.JSON(&protocol.Message{
		Errcode: code,
		Errmsg:  msg,
		Data:    data,
	})
}

// RetBadRequestError(msg string) error
// 	RetForbiddenError(msg string) error
// 	RetNotFoundError(msg string) error
// 	RetInternalError(msg string) error
// 	RetDisplayError(msg string) error
// 	RetCustomError(code int, msg string) error
// 	RetSuccess(msg string, data interface{}) error

func (c *icecontext) RetBadRequestError(msg string) error {
	if msg == "" {
		msg = "bad_request"
	}
	return c.JSON(&protocol.Message{
		Errcode: http.StatusBadRequest,
		Errmsg:  msg,
		Data:    nil,
	})
}

func (c *icecontext) RetForbiddenError(msg string) error {
	if msg == "" {
		msg = "forbidden"
	}
	return c.JSON(&protocol.Message{
		Errcode: http.StatusForbidden,
		Errmsg:  msg,
		Data:    nil,
	})
}

func (c *icecontext) RetNotFoundError(msg string) error {
	if msg == "" {
		msg = "not_found"
	}
	return c.JSON(&protocol.Message{
		Errcode: http.StatusNotFound,
		Errmsg:  msg,
		Data:    nil,
	})
}

func (c *icecontext) RetInternalError(msg string) error {
	if msg == "" {
		return c.RetDisplayError("服务器开了点小差，请稍后再试～")
	}
	return c.JSON(&protocol.Message{
		Errcode: http.StatusInternalServerError,
		Errmsg:  msg,
		Data:    nil,
	})
}

func (c *icecontext) RetDisplayError(msg string) error {
	if msg == "" {
		msg = "unkown_error"
	}
	return c.JSON(&protocol.Message{
		Errcode: 1,
		Errmsg:  msg,
		Data:    nil,
	})
}

func (c *icecontext) RetCustomError(code int, msg string) error {
	if msg == "" {
		msg = "unkown_error"
	}
	return c.JSON(&protocol.Message{
		Errcode: code,
		Errmsg:  msg,
		Data:    nil,
	})
}

func (c *icecontext) RetCustomErrorWithData(code int, msg string, data interface{}) error {
	if msg == "" {
		msg = "unkown_error"
	}
	return c.JSON(&protocol.Message{
		Errcode: code,
		Errmsg:  msg,
		Data:    data,
	})
}

func (c *icecontext) RetSuccess(msg string, data interface{}) error {
	return c.JSON(&protocol.Message{
		Errcode: 0,
		Errmsg:  msg,
		Data:    data,
	})
}

// XML2 公共XML响应
func (c *icecontext) XML2(code int, msg string, data interface{}) error {
	return c.XML(&protocol.Message{
		Errcode: code,
		Errmsg:  msg,
		Data:    data,
	})
}

// XML2 公共XML响应
func (c *icecontext) HTTP(code int, msg string) error {
	return c.Response().FillErr(code, msg)
}

// GetHeaderString 从请求头中获取键值对并转换为string
func (c *icecontext) GetHeaderString(key string, defaultValue string) string {
	if v := c.Header().Get(key); v != "" {
		return v
	}
	return defaultValue
}

// GetHeaderInt64 从请求头中获取键值对并转换为int64
func (c *icecontext) GetHeaderInt64(key string, defaultValue int64) int64 {
	v := c.Header().Get(key)
	i, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return defaultValue
	}
	return i
}

// GetHeaderFloat64 从请求头中获取键值对并转换为float64
func (c *icecontext) GetHeaderFloat64(key string, defaultValue float64) float64 {
	v := c.Header().Get(key)
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return defaultValue
	}
	return f
}

func (c *icecontext) appendBiz(args []interface{}) []interface{} {
	var newargs = make([]interface{}, len(args)+1)
	newargs[0] = "Bizid:" + c.Bizid() + " "
	for i := range args {
		newargs[i+1] = args[i]
	}
	return newargs
}

func (c *icecontext) appendf(fmt string, args []interface{}) (string, []interface{}) {
	return "%s " + fmt, c.appendBiz(args)
}

// Debug global debug
func (c *icecontext) Debug(args ...interface{}) {
	log.Debug(c.appendBiz(args)...)
}

// Warn defalut warn
func (c *icecontext) Warn(args ...interface{}) {
	log.Warn(c.appendBiz(args)...)
}

// Info default info
func (c *icecontext) Info(args ...interface{}) {
	log.Info(c.appendBiz(args)...)
}

// Error default error
func (c *icecontext) Error(args ...interface{}) {
	log.Error(c.appendBiz(args)...)
}

// Fatal default fatal
func (c *icecontext) Fatal(args ...interface{}) {
	log.Fatal(c.appendBiz(args)...)
}

// Debugf global debug
func (c *icecontext) Debugf(fmt string, args ...interface{}) {
	fmt, args = c.appendf(fmt, args)
	log.Debugf(fmt, args...)
}

// Warnf defalut wawrn
func (c *icecontext) Warnf(fmt string, args ...interface{}) {
	fmt, args = c.appendf(fmt, args)
	log.Warnf(fmt, args...)
}

// Infof default info
func (c *icecontext) Infof(fmt string, args ...interface{}) {
	fmt, args = c.appendf(fmt, args)
	log.Infof(fmt, args...)
}

// Errorf default error
func (c *icecontext) Errorf(fmt string, args ...interface{}) {
	fmt, args = c.appendf(fmt, args)
	log.Errorf(fmt, args...)
}

// Fatalf default fatal
func (c *icecontext) Fatalf(fmt string, args ...interface{}) {
	fmt, args = c.appendf(fmt, args)
	log.Fatalf(fmt, args...)
}

func (c *icecontext) GetClientInfo() (appID string, version, platform int64) {
	clientInfo := c.GetHeaderString("Client-Info", "")
	if clientInfo == "" {
		return
	}
	var v = make(map[string]string)
	for _, s := range strings.Split(strings.TrimRight(strings.TrimLeft(clientInfo, "("), ")"), " ") {
		if ss := strings.Split(s, "/"); len(ss) == 2 {
			v[ss[0]] = ss[1]
		}
	}
	version, _ = strconv.ParseInt(strings.Replace(v["v"], ".", "", -1), 10, 64)
	appID = v["id"]
	platform, _ = strconv.ParseInt(v["p"], 10, 64)
	return
}

func (c *icecontext) CrossDomain() {
	if c.resp == nil {
		s := c.Request().Shadow()
		c.resp = &s
	}
	c.Response().SetHeader(protocol.HeaderAccessControlAllowOrigin, "*")
	c.Response().SetHeader(protocol.HeaderAccessControlAllowCredentials, "true")
	c.Response().SetHeader(protocol.HeaderAccessControlAllowHeaders, "X-Requested-With, Content-Type")
	c.Response().SetHeader(protocol.HeaderAccessControlAllowMethods, "PUT,POST,GET,DELETE,OPTIONS")
}

// TODO 默认
func TODO() Context {
	return newContext()
}
