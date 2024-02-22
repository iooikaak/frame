# Gin Default Server

This is API experiment for Gin.

```go
package main

import (
	"github.com/iooikaak/frame/gin"
	"github.com/iooikaak/frame/gin/ginS"
)

func main() {
	ginS.GET("/", func(c *gin.Context) { c.String(200, "Hello World") })
	ginS.Run()
}
```
