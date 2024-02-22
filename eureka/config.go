package eureka

type Config struct {
	DefaultZone                    []string          // eureka服务端地址，逗号分割
	RenewalIntervalInSecs          int               // 心跳间隔，默认30s
	RegistryFetchIntervalSeconds   int               // 获取服务列表间隔，默认15s
	RollDiscoveriesIntervalSeconds int               // 滚动发现地址，默认60s
	DurationInSecs                 int               // 过期间隔，默认90s
	App                            string            // 应用名称
	Port                           int               // 端口
	Metadata                       map[string]string // 元数据
	DataCenterName                 string            // 注册中心名称，eureka把自己也注册到了实例列表里面，用于处理集群健康检查
	StatusUrl                      string            // status接口
	HealthUrl                      string            // health接口
	Instance                       *Instance         `json:"instance"` // 服务实例信息
	ServerOnly                     bool
}
