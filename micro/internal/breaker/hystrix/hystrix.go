package hystrix

import (
	"context"

	"github.com/iooikaak/frame/pproxy/breaker"
	"github.com/micro/go-micro/v2/client"
)

type clientWrapper struct {
	client.Client
}

func (c *clientWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	return breaker.Do(req.Service()+"."+req.Endpoint(), func() error {
		return c.Client.Call(ctx, req, rsp, opts...)
	}, nil)
}

// NewClientWrapper returns a hystrix client Wrapper.
func NewClientWrapper(enabled bool) client.Wrapper {
	return func(c client.Client) client.Client {
		if enabled {
			return &clientWrapper{c}
		}
		return c
	}
}
