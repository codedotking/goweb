package main

import (
	"net/http"
)

// HandlerFunc 是一种函数类型，它接收一个 Context 类型的指针作为参数并且没有返回值
type HandlerFunc func(*Context)

// 定义一个名为 Router 的结构体
type Router struct {
	handlers map[string]HandlerFunc
}

// NewRouter 函数用于创建一个新的路由实例
func NewRouter() *Router {
	return &Router{handlers: make(map[string]HandlerFunc)}
}

// addRoute 函数用于向 Router 添加一条新的路由
func (r *Router) addRoute(method string, pattern string, handler HandlerFunc) {
	key := method + "-" + pattern
	r.handlers[key] = handler
}

// GET 是一个处理 GET 请求的路由处理函数
func (r *Router) GET(pattern string, handler HandlerFunc) {
	r.addRoute("GET", pattern, handler)
}

// POST 是一个 HTTP 方法，通常用于向服务器提交要被处理的数据。这里的 POST 方法用于绑定一个路由模式和处理函数到 HTTP 服务器上。
func (r *Router) POST(pattern string, handler HandlerFunc) {
	r.addRoute("POST", pattern, handler)
}

// Run 在指定的地址上启动路由器服务
func (r *Router) Run(addr string) error {
	return http.ListenAndServe(addr, r)
}

// ServeHTTP 处理传入的 HTTP 请求并调度到适当的处理程序
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	key := req.Method + "-" + req.URL.Path
	if handler, ok := r.handlers[key]; ok {
		c := &Context{
			Writer:  w,
			Request: req,
			Params:  make(map[string]string),
		}
		handler(c)
	} else {
		http.NotFound(w, req)
	}
}
