package redisLock

import (
	"context"
	"reflect"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNew(t *testing.T) {
	Convey("Test GET", t, func() {
		lockInst, err := New(redisConf) //实例化redis
		var wantR RedisLock
		So(err, ShouldBeNil)
		So(reflect.DeepEqual(reflect.TypeOf(lockInst), reflect.TypeOf(&wantR)), ShouldEqual, true)
	})

}
func TestRedisLock_Lock(t *testing.T) {
	Convey("Test GET", t, func() {
		lockInst, err := New(redisConf) //实例化redis
		So(err, ShouldBeNil)
		ctx := context.Background()
		res, err := lockInst.UnLock(ctx, "myLock")
		t.Log(res, err)

		res2, _ := lockInst.Lock(ctx, "myLock", 2)
		res3, _ := lockInst.Lock(ctx, "myLock", 1)
		time.Sleep(2 * time.Second)
		res4, err := lockInst.Lock(ctx, "myLock", 1)
		So(err, ShouldBeNil)
		So(res2, ShouldEqual, true)
		So(res3, ShouldEqual, false)
		So(res4, ShouldEqual, true)
	})
}

func BenchmarkRedisLock_UnLock(b *testing.B) {
	lockInst, err := New(redisConf) //实例化redis
	if err != nil {
		b.Log("发生错误", err)
	}
	ctx := context.Background()

	for i := 0; i < b.N; i++ {
		_, _ = lockInst.UnLock(ctx, "myLock")
	}
}

func BenchmarkRedisLock_Lock(b *testing.B) {
	lockInst, err := New(redisConf) //实例化redis
	if err != nil {
		b.Log("发生错误", err)
	}
	ctx := context.Background()

	for i := 0; i < b.N; i++ {
		b.Log(lockInst.Lock(ctx, "myLock", 1))
	}
}
