package cachex

import (
	"encoding/binary"
	"reflect"
	"time"
	"unsafe"

	"github.com/VictoriaMetrics/fastcache"
)

const timestampSizeInBytes = 8

type fastCache struct {
	fCache *fastcache.Cache
}

func newFastCache(maxBytes int) CacheInterface {
	return &fastCache{fCache: fastcache.New(maxBytes)}
}

func (c *fastCache) Get(name string) (res string, err error) {
	var re []byte
	cName := c.stringToBytes(name)
	re = c.fCache.Get(nil, cName)
	if len(re) > timestampSizeInBytes {
		expire := int64(binary.LittleEndian.Uint64(re))
		if expire > 0 && c.isExpire(expire) {
			c.fCache.Del(cName)
		} else {
			res = c.bytesToString(re[timestampSizeInBytes:])
		}
	}

	return
}

func (c *fastCache) GetBytes(name string) (res []byte, err error) {
	var re []byte
	cName := c.stringToBytes(name)
	re = c.fCache.Get(nil, cName)
	if len(re) > timestampSizeInBytes {
		expire := int64(binary.LittleEndian.Uint64(re))
		if expire > 0 && c.isExpire(expire) {
			c.fCache.Del(cName)
		} else {
			res = re[timestampSizeInBytes:]
		}
	}

	return
}

func (c *fastCache) GetBig(name string) (res string, err error) {
	var re []byte
	cName := c.stringToBytes(name)
	re = c.fCache.GetBig(nil, cName)
	if len(re) > timestampSizeInBytes {
		expire := int64(binary.LittleEndian.Uint64(re))
		if expire > 0 && c.isExpire(expire) {
			c.fCache.Del(cName)
		} else {
			res = c.bytesToString(re[timestampSizeInBytes:])
		}
	}

	return
}
func (c *fastCache) Set(name, val string, expr ...int64) (err error) {
	var ex uint64
	if len(expr) > 0 {
		ex = uint64(time.Now().Add(time.Second * time.Duration(expr[0])).Unix())
	}
	c.fCache.Set(c.stringToBytes(name), append(c.setExpire(ex), c.stringToBytes(val)...))
	return
}

func (c *fastCache) SetBytes(name string, val []byte, expr ...int64) (err error) {
	var ex uint64
	if len(expr) > 0 {
		ex = uint64(time.Now().Add(time.Second * time.Duration(expr[0])).Unix())
	}
	c.fCache.Set(c.stringToBytes(name), append(c.setExpire(ex), val...))
	return
}

func (c *fastCache) SetBig(name, val string, expr ...int64) (err error) {
	var ex uint64
	if len(expr) > 0 {
		ex = uint64(time.Now().Add(time.Second * time.Duration(expr[0])).Unix())
	}
	c.fCache.SetBig(c.stringToBytes(name), append(c.setExpire(ex), c.stringToBytes(val)...))
	return
}

func (c *fastCache) Del(name string) (err error) {
	c.fCache.Del(c.stringToBytes(name))
	return
}

func (c *fastCache) setExpire(expire uint64) (expireByte []byte) {
	expireByte = make([]byte, timestampSizeInBytes)
	binary.LittleEndian.PutUint64(expireByte, expire)
	return
}

func (c *fastCache) isExpire(expire int64) bool {
	return !time.Unix(expire, 0).After(time.Now())
}

func (c *fastCache) stringToBytes(s string) []byte {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := reflect.SliceHeader{
		Data: sh.Data,
		Len:  sh.Len,
		Cap:  sh.Len,
	}
	return *(*[]byte)(unsafe.Pointer(&bh)) //nolint
}

func (c *fastCache) bytesToString(bytes []byte) string {
	return *(*string)(unsafe.Pointer(&bytes))
}
