package cachex

import (
	"context"
	"fmt"
)

type CacheInterface interface {
	Get(name string) (string, error)
	GetBytes(name string) ([]byte, error)
	Set(name, val string, exprBySecond ...int64) error
	SetBytes(name string, val []byte, exprBySecond ...int64) error
	Del(name string) error
	SetBig(name, val string, exprBySecond ...int64) error
	GetBig(name string) (string, error)
}

const (
	FastCache = "fastCache"
)

type Cache struct {
	c        CacheInterface
	maxBytes int
}

//maxBytes too less must greate 64m
func New(maxBytes int, cacheType ...string) (*Cache, error) {

	if maxBytes < 1<<26 {
		return nil, fmt.Errorf("maxBytes too less must greate 64m")
	}
	c := &Cache{c: getCacheInstance(maxBytes, cacheType...), maxBytes: maxBytes}
	return c, nil
}

func (c *Cache) Get(ctx context.Context, name string) (res string, s error) {
	if len(name) == 0 {
		s = fmt.Errorf("name  is empty")
		return
	}

	return c.c.Get(name)
}

func (c *Cache) Set(ctx context.Context, name, val string, exprBySecond ...int64) (err error) {
	if len(name) == 0 || len(val) == 0 {
		err = fmt.Errorf("name or val is empty")
		return
	}

	if len(val) >= ((1 << 16) - timestampSizeInBytes) {
		err = fmt.Errorf("val too large must use SetBig")
		return
	}

	return c.c.Set(name, val, exprBySecond...)
}

func (c *Cache) GetBytes(ctx context.Context, name string) (res []byte, s error) {
	if len(name) == 0 {
		s = fmt.Errorf("name  is empty")
		return
	}

	return c.c.GetBytes(name)
}

func (c *Cache) SetBytes(ctx context.Context, name string, val []byte, exprBySecond ...int64) (err error) {
	if len(name) == 0 || len(val) == 0 {
		err = fmt.Errorf("name or val is empty")
		return
	}

	if len(val) >= ((1 << 16) - timestampSizeInBytes) {
		err = fmt.Errorf("val too large must use SetBig")
		return
	}

	return c.c.SetBytes(name, val, exprBySecond...)
}

func (c *Cache) SetBig(ctx context.Context, name, val string, exprBySecond ...int64) (err error) {
	if len(name) == 0 || len(val) == 0 {
		err = fmt.Errorf("name or val is empty")
		return
	}

	return c.c.SetBig(name, val, exprBySecond...)
}

func (c *Cache) GetBig(ctx context.Context, name string) (res string, s error) {
	if len(name) == 0 {
		s = fmt.Errorf("name  is empty")
		return
	}

	return c.c.GetBig(name)
}

func (c *Cache) Del(ctx context.Context, name string) (err error) {
	if len(name) == 0 {
		return
	}
	return c.c.Del(name)
}

func getCacheInstance(maxBytes int, cacheType ...string) (res CacheInterface) {
	if len(cacheType) == 0 {
		return newFastCache(maxBytes)
	}

	//可使用多种缓存机制
	switch cacheType[0] {
	case FastCache:
		res = newFastCache(maxBytes)

	}

	return
}
