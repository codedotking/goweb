package main

import (
	"fmt"
	"net/http"
)

func main() {
	r := NewRouter()

	// 设置路由
	r.GET("/", func(c *Context) {
		c.String(http.StatusOK, "Hello, World!")
	})

	r.GET("/hello/:name", func(c *Context) {
		name := c.Param("name")
		c.String(http.StatusOK, fmt.Sprintf("Hello, %s!", name))
	})

	// 启动服务器
	r.Run(":8080")
}
