# 项目简介
1.基于fastcache为实现层的本地缓存


# 快速开始
* 使用例子

```go
        //全局实例化
        fastcache, err := New(1 << 26)
    	if err != nil {
    		t.Error(err)
    		return
    	}

        //set操作
    	err := fastcache.Set(context.Background(), "fastcache-set", "fastcache-set", 10)
        if err != nil {
            t.Error(err)
            return
        }

        //get操作
        res, err := fastcache.Get(context.Background(), "fastcache-get")
        if err != nil {
            t.Error(err)
            return
        }

    	err := fastcache.Set(context.Background(), "fastcache-del", "fastcache-del", 10)
        if err != nil {
            t.Error(err)
            return
        }

         //del操作
        err = fastcache.Del(context.Background(), "fastcache-del")
        if err != nil {
            t.Error(err)
            return
        }

```
