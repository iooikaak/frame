package middleware

import (
	"errors"
	"net/http"

	"github.com/iooikaak/frame/metadata"

	"github.com/iooikaak/frame/sig"

	"github.com/iooikaak/frame/gin/render"
	"github.com/iooikaak/frame/gins"
)

var ErrTokenInvalid = errors.New("token error")

func VerifyToken(CommonToken, salt string, force bool) gins.HandlerFunc {
	return func(ctx *gins.Context) {
		token := ctx.Request.Header.Get(metadata.VerifyToken)
		if token != CommonToken && force {
			r, err := sig.GetSign(salt, ctx.Request)

			if err != nil || r != token {
				ctx.Render(http.StatusOK, render.JSON{Data: struct {
					Status int         `json:"status"`
					Msg    string      `json:"msg"`
					Data   interface{} `json:"data"`
				}{
					Status: 90400,
					Msg:    ErrTokenInvalid.Error(),
					Data:   struct{}{},
				}})
				ctx.Abort()
				return
			}
		}
		ctx.Next()
	}
}
