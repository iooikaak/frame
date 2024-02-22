package gins

import (
	"strings"
	"time"

	"github.com/iooikaak/frame/net/metadata"
)

const (
	defaultMaxMemory = 32 << 20 // 32 MB
)

// Logger is logger  middleware
func logger() HandlerFunc {
	const noUser = "no_user"
	return func(c *Context) {
		now := time.Now()
		ip := metadata.String(c, metadata.RemoteIP)
		req := c.Request
		path := req.URL.Path

		ctype := req.Header.Get("Content-Type")
		switch {
		case strings.Contains(ctype, "multipart/form-data"):
			_ = req.ParseMultipartForm(defaultMaxMemory)
		default:
			_ = req.ParseForm()
		}

		params := req.Form
		var quota float64
		// TODO: 将 gins 内部的 Context 集成官方
		if deadline, ok := c.Deadline(); ok {
			quota = time.Until(deadline).Seconds()
		}

		c.Next()

		err := c.Err()
		dt := time.Since(now)
		caller := metadata.String(c, metadata.Caller)
		if caller == "" {
			caller = noUser
		}

		isSlow := dt >= (time.Millisecond * 500)

		// 健康检查不打印日志
		// xlog.Info("path:" + path)
		if path == "/health" {
			return
		}

		vs := make(map[string]interface{})

		vs["timestamp"] = dt.Milliseconds()
		vs["ip"] = ip
		vs["method"] = c.Request.Method
		vs["statusCode"] = c.Writer.Status()
		vs["path"] = path
		vs["params"] = params.Encode()
		vs["type"] = "accesslog"
		vs["timeout_quota"] = quota
		vs["method"] = c.Request.Method

		if err != nil {
			c.Logger().WithFields(vs).Error(err)
		} else {
			if isSlow {
				c.Logger().WithFields(vs).Warn("slow")
			} else {
				c.Logger().WithFields(vs).Info("")
			}
		}
	}
}
