package middleware

import (
	"context"
	"errors"
	"net/http"

	"github.com/iooikaak/frame/gin/render"
	"github.com/iooikaak/frame/gins"
	"github.com/iooikaak/frame/jwt"
	"github.com/iooikaak/frame/metadata"
)

var ErrTokenValidateFrontUserAuth = errors.New("validate front user auth error")

func ValidateFrontUserAuth(mySigningKey interface{}) gins.HandlerFunc {
	return func(ctx *gins.Context) {
		claims, err := jwt.ParseJwtToken(mySigningKey, ctx.Request.Header.Get(metadata.HeaderFrToken))
		if err != nil {
			ctx.Render(http.StatusOK, render.JSON{Data: struct {
				Status int         `json:"status"`
				Msg    string      `json:"msg"`
				Data   interface{} `json:"data"`
			}{
				Status: 90400,
				Msg:    ErrTokenValidateFrontUserAuth.Error(),
				Data:   struct{}{},
			}})
			ctx.Abort()
			return
		}
		ctx.C = context.WithValue(ctx.C, baseInfoKeyFrUser, claims)
		ctx.Next()
	}
}
