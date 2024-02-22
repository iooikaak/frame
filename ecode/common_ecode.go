package ecode

// All common ecode
var (
	OK = add(0) // 正确

	NotModified        = add(10000) // 木有改动
	TemporaryRedirect  = add(90307) // 撞车跳转
	RequestErr         = add(90400) // 请求错误
	Unauthorized       = add(90401) // 未认证
	AccessDenied       = add(90403) // 访问权限不足
	NothingFound       = add(90404) // 啥都木有
	MethodNotAllowed   = add(90405) // 不支持该方法
	Conflict           = add(90409) // 冲突
	Canceled           = add(90498) // 客户端取消请求
	ServerErr          = add(90500) // 服务器错误
	ServiceUnavailable = add(90503) // 过载保护,服务暂不可用
	Deadline           = add(90504) // 服务调用超时
	LimitExceed        = add(90509) // 超出限制
)
