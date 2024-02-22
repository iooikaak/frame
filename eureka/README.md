#### eureka

##### 项目简介

eureka, java 版注册中心，此包用来对接eureka golang版client

```go
// create eureka client
client := NewClient(&Config{
    DefaultZone:           "http://10.1.9.5:8761",
    App:                   "go-example",
    Port:                  10000,
    RenewalIntervalInSecs: 10,
    DurationInSecs:        30,
    DataCenterName:        "discovery",
    Metadata: map[string]interface{}{
        "VERSION":              "0.1.0",
        "NODE_GROUP_ID":        0,
        "PRODUCT_CODE":         "DEFAULT",
        "PRODUCT_VERSION_CODE": "DEFAULT",
        "PRODUCT_ENV_CODE":     "DEFAULT",
        "SERVICE_VERSION_CODE": "DEFAULT",
    },
})
// start client, register、heartbeat、refresh
client.Start()

// http server
http.HandleFunc("/v1/services", func(writer http.ResponseWriter, request *http.Request) {
    // full applications from eureka server
    apps := client.Applications
    b, _ := json.Marshal(apps)
    _, _ = writer.Write(b)
})

// start http server
if err := http.ListenAndServe(":10000", nil); err != nil {
    fmt.Println(err)
}
```

##### eureka 官方文档

https://github.com/Netflix/eureka/wiki/Eureka-REST-operations
