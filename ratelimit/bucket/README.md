# 项目简介

本限流是基于redis基础之上的令牌限流控制，可在多个节点之间进行限流，支持参数可配置

#快速开始

## 先实例化redis

    err := sdk.Client.Init(ServiceName, false)
 	if err != nil {
 		panic(err)
 	}

 	confStr, err := sdk.Client.GetRedisConfig(ServiceName) // redisInstanceName指的是redis配置的name值
 	if err != nil {
 		panic("无法从 goms 获取redis配置信息：" + err.Error())
 	}

 	err = redis.GetRedisForJsonStr(confStr)
 	if err != nil {
 		panic("redis配置错误：" + err.Error())
 	}

 	RedisInstance, err = cache.InitRedisInstance(ServiceName, ServiceName)
 	if err != nil {
 		panic(err.Error())
 	}
 	//检测实例链接是否可用
 	_, err = RedisInstance.Str().Ping(true)
 	if err != nil {
 		panic(err.Error())
 	}

## 将redis连接添加，并传入每次令牌桶的每次申请的数量和令牌的上限

    // applyNum 每次从redis中申请的令牌数量
    // Capacity 服务每秒的令牌上限 实际最大上限为applyNum+Capacity
    // redisPrefix redis key名称
    // action 服务名称
    func NewBucket(applyNum, Capacity int64,redisPrefix,action string) *ratelimit.Bucket {
   	prefixKey := redisPrefix + action
   	bucketOptions := ratelimit.NewDefaultOptions(
   		ratelimit.SetRedisOptionsPrefix(prefixKey),
   		ratelimit.SetRedisOptionsCapacity(Capacity),
   		ratelimit.SetRedisInstance(RedisInstance),
   	)

   	bucket := ratelimit.NewBucket(
   		ratelimit.SetApplyNum(applyNum),
   		ratelimit.SetAvailableTokens(applyNum),
   		ratelimit.SetOptions(bucketOptions),
   		ratelimit.SetQuantum(applyNum),
   		ratelimit.SetCapacity(applyNum),
   	)

   	return bucket

}

## 在获取到Bucket之后，每次有请求进来之后，减去相应的自己规定的最大数量

    var count = 1
    bucker := NewBucket(5, 40, "redidPrefix", "test")
    bucker.Take(count)
