package tracer

import (
	"context"
	"fmt"

	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/metadata"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/server"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

func StartSpanFromContext(ctx context.Context, tracer opentracing.Tracer, name string, opts ...opentracing.StartSpanOption) (context.Context, opentracing.Span, error) {
	md, ok := metadata.FromContext(ctx)
	if !ok {
		md = make(map[string]string)
	}

	//防止数据竞争
	md = metadata.Copy(md)

	if parentSpan := opentracing.SpanFromContext(ctx); parentSpan != nil {
		opts = append(opts, opentracing.ChildOf(parentSpan.Context()))
	} else if spanCtx, err := tracer.Extract(opentracing.TextMap, opentracing.TextMapCarrier(md)); err == nil {
		opts = append(opts, opentracing.ChildOf(spanCtx))
	}

	sp := tracer.StartSpan(name, opts...)

	if err := sp.Tracer().Inject(sp.Context(), opentracing.TextMap, opentracing.TextMapCarrier(md)); err != nil {
		return nil, nil, err
	}

	ctx = opentracing.ContextWithSpan(ctx, sp)
	ctx = metadata.NewContext(ctx, md)
	return ctx, sp, nil
}

// server form go-micro
func NewHandlerWrapper(t opentracing.Tracer) server.HandlerWrapper {
	return func(h server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			name := fmt.Sprintf("%s.%s", req.Service(), req.Endpoint())
			ctx, span, err := StartSpanFromContext(ctx, t, name)
			if err != nil {
				return err
			}
			ext.SpanKindRPCServer.Set(span)

			err = h(ctx, req, rsp)
			if err != nil {
				span.SetTag("handlerError", err)
			}
			span.Finish()
			return err
		}
	}
}

//client form go-micro
func NewCallWrapper(t opentracing.Tracer, peerService func() string) client.CallWrapper {
	return func(cf client.CallFunc) client.CallFunc {
		return func(ctx context.Context, node *registry.Node, req client.Request, rsp interface{}, opts client.CallOptions) error {
			name := fmt.Sprintf("%s.%s", req.Service(), req.Endpoint())
			ctx, span, err := StartSpanFromContext(ctx, t, name)
			if err != nil {
				return err
			}
			ext.SpanKindRPCClient.Set(span)
			ext.PeerAddress.Set(span, node.Address)
			ext.PeerService.Set(span, peerService())

			err = cf(ctx, node, req, rsp, opts)
			if err != nil {
				span.SetTag("handlerError", err)
			}
			span.Finish()
			return err
		}
	}
}
