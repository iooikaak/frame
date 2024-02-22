package gorm

import (
	"testing"
	"time"
)

func TestGormNewInstance(t *testing.T) {
	db := New(&Config{
		Alias:        "test",
		Type:         "mysql",
		Server:       "10.1.2.13",
		Port:         3306,
		Database:     "channel_center",
		User:         "dev_niumowang",
		Password:     "1QAZ2wsx",
		MaxIdleConns: 200,
		MaxOpenConns: 500,
		Charset:      "utf8mb4",
		MaxLeftTime:  time.Second * 10,
	})
	if err := db.gorm.DB().Ping(); err != nil {
		t.Error("数据库连接失败")
	}
	t.Log("数据库连接成功")
}
