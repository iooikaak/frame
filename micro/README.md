# 快速开始

* micro使用例子

```go
//step 1：new micro实例
//TODO 记得配置CONSUL_IP环境变量

m, err := New("go.micro.srv.xxx")
if err != nil {
panic("parse micro err：" + err.Error())
}

//step 2：初始化micro参数
err = m.Initialize()
if err != nil {
panic("micro init failed err：" + err.Error())
}

//grcp 需要先生成pb文件
if err := pb.RegisterSupplierHandler(m.Server(), new(handler.SupplierHanders)); err != nil {
xlog.Error(err)
}
//step 2：运行
if err := m.Run(); err != nil {
panic(err)
}

```
