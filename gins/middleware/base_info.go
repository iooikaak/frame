package middleware

import (
	"context"
	"encoding/hex"
	"regexp"
	"strings"

	"github.com/iooikaak/frame/gins"
	"github.com/iooikaak/frame/metadata"
	"github.com/iooikaak/frame/utils"
)

const (
	_appPlatformMf  = "mf"
	_appPlatformCyy = "cyy"
)

var (
	appPlatformReg = regexp.MustCompile("(?i)cyy/([0-9.]+)")
	versionReg     = regexp.MustCompile("(?i)mf/([0-9.]+)")
	clientCodeReg  = regexp.MustCompile(`(?i)sc\((.+?)\)`)
	networkReg     = regexp.MustCompile(`(?i)network/([\w]+)`)
	minVersionReg  = regexp.MustCompile(`(?i)minVersion\(([0-9.]+)\)`)
	mfDvSignReg    = regexp.MustCompile(`mf-dv-sign[[\w]+\|(\w+)]`)
	brandReg       = regexp.MustCompile(`^\w+ [\d.]+ (\w+) `)
)

type GomsUser struct {
	Uid      int    `json:"uid"`
	UserName string `json:"username"`
	NickName string `json:"nickname"`
	Uex      string `json:"uex"`
	U        string `json:"u"`
	Uage     int    `json:"uage"`
}

type contextInfoKey string

const (
	baseInfoKeyAppPlatform contextInfoKey = metadata.HttpMiddleWareAppPlatform
	baseInfoKeyClientCode  contextInfoKey = metadata.HttpMiddleWareClientCode
	baseInfoKeyVersion     contextInfoKey = metadata.HttpMiddleWareVersion
	baseInfoKeyNetwork     contextInfoKey = metadata.HttpMiddleWareNetwork
	baseInfoKeyPlatform    contextInfoKey = metadata.HttpMiddleWarePlatform
	baseInfoKeyMinVersion  contextInfoKey = metadata.HttpMiddleWareMinVersion
	baseInfoKeyMobileBrand contextInfoKey = metadata.HttpMiddleWareMobileBrand
	baseInfoKeyBeUser      contextInfoKey = metadata.HttpMiddleWareBeUser
	baseInfoKeyFrUser      contextInfoKey = metadata.HttpMiddleWareFrUser
	baseInfoKeyBuUser      contextInfoKey = metadata.HttpMiddleWareBuUser
	abTestKey              contextInfoKey = metadata.HttpMiddleWareABTest
)

// 注入app请求的参数, 设备号,版本,平台 大礼包
func ClientBaseInfo() gins.HandlerFunc {
	return func(ctx *gins.Context) {
		agent := ctx.Request.UserAgent()

		ok := appPlatformReg.MatchString(agent)
		if ok {
			ctx.C = context.WithValue(ctx.C, baseInfoKeyAppPlatform, _appPlatformCyy)
		} else {
			ctx.C = context.WithValue(ctx.C, baseInfoKeyAppPlatform, _appPlatformMf)
		}

		// 注入clientCode
		res := clientCodeReg.FindAllStringSubmatch(agent, -1)
		if len(res) > 0 && len(res[0]) > 1 {
			// c58f4f4d242f42b3,oppo
			str := strings.Split(res[0][1], ",")
			if str[0] == "{holder}" {
				sign := mfDvSignReg.FindStringSubmatch(agent)
				if len(sign) == 2 {
					data := []byte(sign[1])
					src, _ := hex.DecodeString(string(data))
					key := []byte("mf")
					str, err := utils.DesECBDecrypt(src, key, utils.PKCS7_PADDING)
					if len(str) > 0 && err == nil {
						ctx.C = context.WithValue(ctx.C, baseInfoKeyClientCode, string(str))
					}
				}
			} else {
				ctx.C = context.WithValue(ctx.C, baseInfoKeyClientCode, strings.TrimSpace(str[0]))
			}
		}

		// 注入版本v
		v := ctx.Query("v")
		if v == "" {
			res := versionReg.FindAllStringSubmatch(agent, -1)

			if len(res) > 0 && len(res[0]) > 1 {
				v = res[0][1]
			}
		}
		// 注入版本号v
		ctx.C = context.WithValue(ctx.C, baseInfoKeyVersion, v)

		ctx.Next()
	}
}

func ClientInfo() gins.HandlerFunc {
	return func(ctx *gins.Context) {
		agent := ctx.Request.UserAgent()

		ok := appPlatformReg.MatchString(agent)
		if ok {
			ctx.C = context.WithValue(ctx.C, baseInfoKeyAppPlatform, _appPlatformCyy)
		} else {
			ctx.C = context.WithValue(ctx.C, baseInfoKeyAppPlatform, _appPlatformMf)
		}

		// 注入clientCode
		res := clientCodeReg.FindAllStringSubmatch(agent, -1)
		if len(res) > 0 && len(res[0]) > 1 {
			// c58f4f4d242f42b3,oppo
			str := strings.Split(res[0][1], ",")
			if str[0] == "{holder}" {
				sign := mfDvSignReg.FindStringSubmatch(agent)
				if len(sign) == 2 {
					data := []byte(sign[1])
					src, _ := hex.DecodeString(string(data))
					key := []byte("mf")
					str, err := utils.DesECBDecrypt(src, key, utils.PKCS7_PADDING)
					if len(str) > 0 && err == nil {
						ctx.C = context.WithValue(ctx.C, baseInfoKeyClientCode, string(str))
					}
				}
			} else {
				ctx.C = context.WithValue(ctx.C, baseInfoKeyClientCode, strings.TrimSpace(str[0]))
			}
		}

		// 注入版本v
		v := ctx.Query("v")
		if v == "" {
			res := versionReg.FindAllStringSubmatch(agent, -1)

			if len(res) > 0 && len(res[0]) > 1 {
				v = res[0][1]
			}
		}
		// 注入版本号v
		ctx.C = context.WithValue(ctx.C, baseInfoKeyVersion, v)

		//注入网络类型
		net := networkReg.FindStringSubmatch(agent)
		if len(net) == 2 {
			ctx.C = context.WithValue(ctx.C, baseInfoKeyNetwork, net[1])
		}

		//注入平台platform
		platform := ctx.GetHeader("platform")
		if strings.TrimSpace(platform) == "" {
			agent := ctx.Request.UserAgent()
			if strings.Contains(agent, "iphone") || strings.Contains(agent, "mac os") {
				platform = "ios"
			} else {
				platform = "android"
			}
		}
		if platform != "" {
			ctx.C = context.WithValue(ctx.C, baseInfoKeyPlatform, platform)
		}

		//注入最小版本号
		minV := minVersionReg.FindStringSubmatch(agent)
		if len(minV) > 1 {
			ctx.C = context.WithValue(ctx.C, baseInfoKeyMinVersion, minV[1])
		}

		//注入手机品牌
		brand := brandReg.FindStringSubmatch(agent)
		if len(brand) == 2 {
			ctx.C = context.WithValue(ctx.C, baseInfoKeyMobileBrand, brand[1])
		}

		ctx.Next()
	}
}
