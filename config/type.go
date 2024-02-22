package config

//ConfigType 配置对象类型
type ConfigType int

const (
	CONFIG_TYPE_FILE ConfigType = iota
	CONFIG_TYPE_CONSUL
)

func (t ConfigType) String() string {
	switch t {
	case CONFIG_TYPE_FILE:
		return "File"
	case CONFIG_TYPE_CONSUL:
		return "Consul"
	}

	return "unknow"
}
