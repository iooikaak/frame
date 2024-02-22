# 项目简介

HTTP 加签名

1. GET请求 参数用所有GET参数，以下划线开头的除外，key1=value1&key2=value2类型排序，然后拼接key，取md5值，GET请求复杂
   参数改用POST JSON BODY，此处不支持

2. POST请求 统一使用POST JSON BODY，对JSON 字符串拼接key，取md5值

```
