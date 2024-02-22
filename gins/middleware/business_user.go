package middleware

import (
	"context"
	"errors"
	"net/http"

	"github.com/iooikaak/frame/utils"

	"github.com/iooikaak/frame/gin/render"
	"github.com/iooikaak/frame/gins"
	"github.com/iooikaak/frame/jwt"
	"github.com/iooikaak/frame/metadata"
)

var ErrTokenValidateBusinessUserAuth = errors.New("validate business user auth error")
var ErrTokenValidateBusinessUserInterfaceAuth = errors.New("validate business user auth error")

func ValidateBusinessUserAuth(mySigningKey interface{}, roles map[int64][]string, interfaceCheck bool) gins.HandlerFunc {
	return func(ctx *gins.Context) {
		claims, err := jwt.ParseJwtToken(mySigningKey, ctx.Request.Header.Get(metadata.HeaderBeToken))
		if err != nil {
			ctx.Render(http.StatusOK, render.JSON{Data: struct {
				Status int         `json:"status"`
				Msg    string      `json:"msg"`
				Data   interface{} `json:"data"`
			}{
				Status: 90400,
				Msg:    ErrTokenValidateBusinessUserAuth.Error(),
				Data:   struct{}{},
			}})
			ctx.Abort()
			return
		}
		if interfaceCheck {
			var ok bool
			var role []string
			role, ok = roles[claims.UserInfo.RoleId]
			if !ok || !utils.CheckIsExistString(role, ctx.Request.URL.Path) {
				ctx.Render(http.StatusOK, render.JSON{Data: struct {
					Status int         `json:"status"`
					Msg    string      `json:"msg"`
					Data   interface{} `json:"data"`
				}{
					Status: 90400,
					Msg:    ErrTokenValidateBusinessUserInterfaceAuth.Error(),
					Data:   struct{}{},
				}})
				ctx.Abort()
				return
			}
		}
		ctx.C = context.WithValue(ctx.C, baseInfoKeyBuUser, claims)
		ctx.Next()
	}
}
