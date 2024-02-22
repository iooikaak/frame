package cachex

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/VictoriaMetrics/fastcache"
)

func BenchmarkFastCacheSet(b *testing.B) {
	cache := fastcache.New(32 * 1024 * 1024)
	for i := 0; i < b.N; i++ {
		cache.Set([]byte(key(i)), value())
	}
}

func BenchmarkFastCacheGet(b *testing.B) {
	b.StopTimer()
	cache := fastcache.New(64 * 1024 * 1024)
	for i := 0; i < b.N; i++ {
		cache.Set([]byte(key(i)), value())
	}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		cache.Get(nil, []byte(key(i)))
	}
}

func BenchmarkFastCacheSetParallel(b *testing.B) {
	cache := fastcache.New(64 * 1024 * 1024)
	rand.Seed(time.Now().Unix())

	b.RunParallel(func(pb *testing.PB) {
		id := rand.Intn(1000)
		counter := 0
		for pb.Next() {
			cache.Set([]byte(parallelKey(id, counter)), value())
			counter = counter + 1
		}
	})
}

func BenchmarkFastCacheGetParallel(b *testing.B) {
	b.StopTimer()
	cache := fastcache.New(64 * 1024 * 1024)
	for i := 0; i < b.N; i++ {
		cache.Set([]byte(key(i)), value())
	}

	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {
		counter := 0
		for pb.Next() {
			cache.Get(nil, []byte(key(counter)))
			counter = counter + 1
		}
	})
}

func key(i int) string {
	return fmt.Sprintf("key-%010d", i)
}

func value() []byte {
	return make([]byte, 100)
}

func parallelKey(threadID int, counter int) string {
	return fmt.Sprintf("key-%04d-%06d", threadID, counter)
}

func TestFastCache(t *testing.T) {

	fastcache, err := New(1 << 26)
	if err != nil {
		t.Error(err)
		return
	}

	t.Run("fastcache-set", func(t *testing.T) {
		err := fastcache.Set(context.Background(), "fastcache-set", "fastcache-set", 10)
		if err != nil {
			t.Error(err)
			return
		}
	})

	t.Run("fastcache-get", func(t *testing.T) {
		err := fastcache.Set(context.Background(), "fastcache-get", "fastcache-get", 10)
		if err != nil {
			t.Error(err)
			return
		}

		res, err := fastcache.Get(context.Background(), "fastcache-get")
		if err != nil {
			t.Error(err)
			return
		}
		t.Log("result:", res)
	})

	t.Run("fastcache-del", func(t *testing.T) {
		err := fastcache.Set(context.Background(), "fastcache-del", "fastcache-del", 10)
		if err != nil {
			t.Error(err)
			return
		}
		err = fastcache.Del(context.Background(), "fastcache-del")
		if err != nil {
			t.Error(err)
			return
		}
	})

}
