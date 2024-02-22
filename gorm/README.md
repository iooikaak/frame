# 项目简介
1.基于gorm库，做了封装，增加了tracing

# 详细文档地址
https://jasperxu.github.io/gorm-zh


# 快速开始
```go
    db := gorm.New(&config.Config{
        Type:     "mysql",
        Server:   "127.0.0.1",
        Port:     3306,
        Database: "db",
        User:     "db",
        Password: "123456",
    })

    //应用中使用，sess的具体方法查看 https://jasperxu.github.io/gorm-zh
    sess = db.Context(context.Background())
```
