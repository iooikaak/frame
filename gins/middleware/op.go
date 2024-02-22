package middleware

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/iooikaak/frame/metadata"

	"github.com/iooikaak/frame/gins"
	"github.com/iooikaak/frame/jwt"
	"github.com/iooikaak/frame/xlog"
)

var (
	_timeFormat = "2006-01-02 15:04:05"
)

type OPRecords struct {
	logger xlog.ILog
}

func (o *OPRecords) OPLog(mySigningKey interface{}) gins.HandlerFunc {
	return func(c *gins.Context) {
		claims, err := jwt.ParseJwtToken(mySigningKey, c.Request.Header.Get(metadata.HeaderBeToken))
		if err != nil {
			o.logger.Errorf("OPLog ParseJwtToken failed err: %v", err)
		}
		uri := fmt.Sprintf("%s://%s%s?%s", c.Request.URL.Scheme, c.Request.Host, c.Request.URL.Path, c.Request.URL.RawQuery)
		now := time.Now().Local().Format(_timeFormat)
		data, berr := ioutil.ReadAll(c.Request.Body)
		var nickname string
		var userId int64
		if claims != nil && claims.UserInfo != nil {
			nickname = claims.UserInfo.UserName
			userId = claims.UserInfo.UserId
		}

		fields := make(map[string]interface{})
		fields["username"] = nickname
		fields["time"] = now
		fields["url"] = c.Request.URL.Path
		fields["param1"] = uri
		fields["ip"] = c.ClientIP()
		fields["uid"] = userId
		fields["type"] = c.Request.Method
		fields["param3"] = string(data)
		o.logger.WithFields(fields).Infof("Username: %s request: %s by: %s at: %s", nickname, c.Request.URL.Path, c.Request.Method, now)

		if berr == nil {
			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(data))
		}

		c.Next()
	}
}
