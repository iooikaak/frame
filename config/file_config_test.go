package config

import (
	"testing"
)

func TestFileConfig(t *testing.T) {
	c, err := NewFileConfig("type.go")
	if err != nil {
		t.Error(err)
		return
	}

	body, err := c.Body()
	if err != nil {
		t.Error(err)
		return
	}

	t.Log("File配置内容：", string(body))
}
