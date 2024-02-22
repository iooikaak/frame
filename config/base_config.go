package config

type BaseConfig struct {
	typ ConfigType
}

//Type 获取配置对象的类型
func (c *BaseConfig) Type() ConfigType {
	return c.typ
}
