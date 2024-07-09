package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Context 是一个上下文类型
type Context struct {
	Writer  http.ResponseWriter
	Request *http.Request
	Params  map[string]string
}

// Param 函数通过传入键 key 获取 Context 结构体中 Params 映射里对应的值
func (c *Context) Param(key string) string {
	return c.Params[key]
}

// String 函数用于向客户端返回纯文本格式的字符串响应。
// 参数 code 表示 HTTP 状态码，format 表示格式化字符串，values 表示格式化字符串所需的参数。
func (c *Context) String(code int, format string, values ...interface{}) {
	c.Writer.WriteHeader(code)
	c.Writer.Header().Set("Content-Type", "text/plain")
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

// JSON 函数用于将一个对象编码为 JSON 格式，并将其作为 HTTP 响应返回。
func (c *Context) JSON(code int, obj interface{}) {
	c.Writer.WriteHeader(code)
	c.Writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(c.Writer).Encode(obj)
}
