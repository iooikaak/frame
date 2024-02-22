package xorm

import (
	"context"
	"testing"
)

var _testConfig = &Config{
	Type:     "mysql",
	Server:   "10.1.2.13",
	Port:     3306,
	Database: "channel_center",
	User:     "dev_niumowang",
	Password: "1QAZ2wsx",
}

func TestDB(t *testing.T) {
	var err error
	//step 1: 实例化db
	db, err := New(_testConfig)
	if err != nil {
		panic(err)
	}

	//step 2: 任何地方使用db链接
	t.Log(db.Ping())

	sql := "select * from  trd_supplier_activity limit 2"
	sess := db.NewSession(context.Background())
	defer sess.Close()
	var am = make([]map[string]interface{}, 0)
	//纯sql 查询多条
	err = sess.SQL(sql).Find(&am)
	if err != nil {
		panic(err)
	}
	t.Log(am)
}

func TestDBWhereQuery(t *testing.T) {
	var err error
	db, err := New(_testConfig)
	if err != nil {
		panic(err)
	}
	sess := db.NewSession(context.Background())
	var bean interface{}
	_, err = sess.Table("trd_shop_info").Where("id = 1 or name = ?", "default").Get(&bean)
	if err != nil {
		t.Error(err)
	}
	t.Log(sess.LastSQL())
}

func BenchmarkSelect(b *testing.B) {
	var err error
	//step 1: 实例化db
	db, err := New(&Config{
		Type:     "mysql",
		Server:   "10.1.2.13",
		Port:     3306,
		Database: "channel_center",
		User:     "dev_niumowang",
		Password: "1QAZ2wsx",
	})

	if err != nil {
		panic(err)
	}

	b.ResetTimer()
	sql := "select goods_id from  trd_supplier_activity limit 2"
	var am = make([]map[string]interface{}, 0)
	sess := db.NewSession(context.Background())
	defer sess.Close()
	for i := 0; i < b.N; i++ {
		//纯sql
		err = sess.SQL(sql).Find(&am)
		if err != nil {
			b.Error(err)
		}
	}

}
