package middleware

import (
	"context"
	"strings"

	"github.com/iooikaak/frame/gins"
)

const _abTestControl = "Abtest-Control"

func ABTest() gins.HandlerFunc {
	return func(ctx *gins.Context) {
		abTestControl := ctx.Request.Header.Get(_abTestControl)
		if abTestControl == "" {
			ctx.Next()
			return
		}

		testList := strings.Split(abTestControl, ";")
		if len(testList) == 0 {
			ctx.Next()
			return
		}

		data := map[string]string{}
		for _, comma := range testList {
			arr := strings.Split(comma, "=")
			if len(arr) != 2 {
				continue
			}
			data[strings.TrimSpace(arr[0])] = strings.TrimSpace(arr[1])
		}
		ctx.C = context.WithValue(ctx.C, abTestKey, data)

		ctx.Next()
	}

}
