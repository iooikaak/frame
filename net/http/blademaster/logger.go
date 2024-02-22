package blademaster

import (
	"strconv"
	"time"

	"github.com/iooikaak/frame/ecode"
	"github.com/iooikaak/frame/net/metadata"
	log "github.com/iooikaak/frame/xlog"
)

// Logger is logger  middleware
func Logger() HandlerFunc {
	const noUser = "no_user"
	return func(c *Context) {
		now := time.Now()
		ip := metadata.String(c, metadata.RemoteIP)
		req := c.Request
		path := req.URL.Path
		params := req.Form
		var quota float64
		if deadline, ok := c.Context.Deadline(); ok {
			quota = time.Until(deadline).Seconds()
		}

		c.Next()

		err := c.Error
		cerr := ecode.Cause(err)
		dt := time.Since(now)
		caller := metadata.String(c, metadata.Caller)
		if caller == "" {
			caller = noUser
		}

		if len(c.RoutePath) > 0 {
			_metricServerReqCodeTotal.Inc(c.RoutePath[1:], caller, req.Method, strconv.FormatInt(int64(cerr.Code()), 10))
			_metricServerReqDur.Observe(int64(dt/time.Millisecond), c.RoutePath[1:], caller, req.Method)
		}

		isSlow := dt >= (time.Millisecond * 500)

		vs := make(map[string]interface{})

		vs["ts"] = dt.Seconds()
		vs["ip"] = ip
		vs["method"] = c.Request.Method
		//vs["statusCode"] = c.Writer.
		vs["path"] = path
		vs["params"] = params.Encode()
		vs["type"] = "http-access-log"
		vs["timeout_quota"] = quota
		vs["errmsg"] = ""
		vs["method"] = c.Request.Method

		if err != nil {
			vs["errmsg"] = err.Error()
			log.ErrorWithFields(vs)
		} else {
			if isSlow {
				log.WarnWithFields(vs)
			} else {
				log.InfoWithFields(vs)
			}
		}
	}
}
