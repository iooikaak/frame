#### paladin

##### 项目简介

paladin 是一个config SDK客户端，包括了file、mock几个抽象功能，方便使用本地文件或者sven\apollo配置中心，并且集成了对象自动reload功能。  

local files:
```
demo -conf=/data/conf/app/msm-servie.toml
// or dir
demo -conf=/data/conf/app/
```

*注：使用远程配置中心的用户在执行应用，如这里的`demo`时务必**不要**带上`-conf`参数，具体见下文远程配置中心的例子*

example:
```
见example_test文件

```

##### 编译环境

- **请只用 Golang v1.12.x 以上版本编译执行**

##### 依赖包
