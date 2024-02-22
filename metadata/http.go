package metadata

const (
	HttpTraceId        = "trace_id" // HttpTraceId 链路ID
	HttpFrom           = "from"     // 来自http中的_userAgent
	HttpCircuitBreaker = "cb"       // HTTP_CIRCUIT_BREAKER 降级状态

	// 中间件拦截获取的一些参数
	HttpMiddleWareABTest      = "ABTest"      // ABTest abtest的map
	HttpMiddleWareAppPlatform = "appPlatform" // platform app平台: mf
	HttpMiddleWareClientCode  = "clientCode"
	HttpMiddleWareVersion     = "v"            // 版本
	HttpMiddleWareFrUser      = "fr_user"      // 前台用户信息
	HttpMiddleWareBuUser      = "bu_user"      // 商户用户信息
	HttpMiddleWareBeUser      = "be_user"      // 后台台用户信息
	HttpMiddleWareNetwork     = "network"      // 网络
	HttpMiddleWarePlatform    = "platform"     // 手机平台
	HttpMiddleWareMobileBrand = "mobile_brand" // 手机品牌
	HttpMiddleWareMinVersion  = "min_version"  // minVersion
	HttpColor                 = "color"        // HttpColor 流量染色的颜色
)
