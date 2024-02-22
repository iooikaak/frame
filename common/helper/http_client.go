package helper

import (
	"time"

	bm "github.com/iooikaak/frame/net/http/blademaster"
)

var HTTPClient *bm.Client

func init() {
	c := &bm.ClientConfig{
		Dial:      200 * time.Millisecond,
		Timeout:   10 * time.Second,
		KeepAlive: 60 * time.Second,
	}
	HTTPClient = bm.NewClient(c)
}
