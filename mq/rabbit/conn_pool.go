package rabbit

import (
	"context"

	pool "github.com/jolestar/go-commons-pool/v2"
)

// RabbitPoolFactory 创建Rabbit连接池工厂类
type connFactory struct {
	cfg *Config
}

// MakeObject 创建对象方法
func (f *connFactory) MakeObject(ctx context.Context) (*pool.PooledObject, error) {
	conn, err := Dial(f.cfg.DSN, f.cfg)
	return pool.NewPooledObject(conn), err
}

// DestroyObject 回收对象方法
func (f *connFactory) DestroyObject(ctx context.Context, object *pool.PooledObject) error {
	v := object.Object.(*Connection)
	if !v.IsClosed() {
		return v.Close()
	}
	return nil
}

// ValidateObject 验证对象有效性
func (f *connFactory) ValidateObject(ctx context.Context, object *pool.PooledObject) bool {
	return object.Object.(*Connection).IsClosed()
}

// ActivateObject 激活对象
func (f *connFactory) ActivateObject(ctx context.Context, object *pool.PooledObject) error {
	// v := object.Object.(*Connection)
	return nil
}

// PassivateObject 钝化对象
func (f *connFactory) PassivateObject(ctx context.Context, object *pool.PooledObject) error {
	return nil
}
