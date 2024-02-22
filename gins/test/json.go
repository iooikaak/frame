package test

import (
	"encoding/json"
)

// TestJSON JSON参数对象
type TestJSON struct {
	data map[string]interface{}
}

// NewTestJSON 新JSON参数对象
func NewTestJSON() *TestJSON {
	return &TestJSON{
		data: make(map[string]interface{}),
	}
}

// Add 添加参数
func (j *TestJSON) Add(key string, value interface{}) {
	j.data[key] = value
}

// Body 获取JSON字符串
func (j *TestJSON) Body() []byte {
	body, _ := json.Marshal(j.data)
	return body
}
