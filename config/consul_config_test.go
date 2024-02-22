package config

import "testing"

func TestConsulConfigSync(t *testing.T) {
	configStr := `{"Server":"120.25.94.79", "Port": 8500, "DefaultKey": "testConfig"}`

	c, err := NewConsulConfig(configStr)
	if err != nil {
		t.Error(err)
		return
	}

	body, err := c.Body()
	if err != nil {
		t.Error(err)
		return
	}

	t.Log("Consul配置内容：", string(body))
}

func TestConsulConfigAsync(t *testing.T) {
	configStr := `{"Server":"120.25.94.79", "Port": 8500}`

	c, err := NewConsulConfig(configStr)
	if err != nil {
		t.Error(err)
		return
	}

	key := "testConfig"

	c.Watch(key)

	body, err := c.Get(key)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log("Consul配置内容：", string(body))
}
