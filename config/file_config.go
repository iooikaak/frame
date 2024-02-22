package config

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
)

//FileConfig
type FileConfig struct {
	BaseConfig
	configPath string
}

//NewFileConfig 创建File配置对象
func NewFileConfig(configStr string) (*FileConfig, error) {
	if configStr == "" {
		return nil, errors.New("configStr不能为空")
	}

	appPath, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	//将平台路径分隔符转换为'/'
	//FromSlash方法，将'/'转换为系统相关的路径分隔符
	configPath := filepath.Join(appPath, filepath.ToSlash(configStr))

	c := &FileConfig{
		configPath: configPath,
	}
	c.typ = CONFIG_TYPE_FILE

	return c, nil
}

//Body 获取配置文件内容
func (c *FileConfig) Body() ([]byte, error) {
	return ioutil.ReadFile(c.configPath)
}
