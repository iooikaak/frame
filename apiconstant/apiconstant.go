package apiconstant

// ResponseType 响应类型
type ResponseType int

// 常量的定义以客户端可能需要进行特殊处理为建立标准
const (
	// 空响应
	RESPONSE_UNKNOW ResponseType = -1

	// 正常响应
	RESPONSE_OK ResponseType = 0

	// 正常处理，但没有找到对应的数据
	RESPONSE_NO_DATA ResponseType = 10000

	// 用户参数校验失败
	RESPONSE_PARAM_INVALID ResponseType = 90400

	// 常规错误
	RESPONSE_ERROR ResponseType = 90000

	// 用户登录令牌无效，含过期
	RESPONSE_TOKEN_INVALID ResponseType = 90401

	// 没有接口权限
	RESPONSE_RBAC_INVALID ResponseType = 90403

	// Action无效
	RESPONSE_ACTION_INVALID ResponseType = 90404

	// AccessToken无效
	RESPONSE_ACCESS_TOKEN_INVALID ResponseType = 90407

	// 请求数超限被拒绝
	RESPONSE_REJECT ResponseType = 90429

	// 响应异常
	RESPONSE_CRASH ResponseType = 90500

	// 服务不可用
	RESPONSE_SERVICE_INVALID ResponseType = 90503
)
