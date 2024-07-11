package goweb

import (
	"net/http"

	"github.com/techiehe/goweb/router"
)

type Framework struct {
	router *router.Router
}

func New() *Framework {
	return &Framework{
		router: router.New(),
	}
}

// GET
func (f *Framework) GET(path string, handleFunc http.HandlerFunc) {
}

func (f *Framework) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f.router.ServeHTTP(w, r)
}

// Run 启动
func (f *Framework) Run(addr string) error {
	return http.ListenAndServe(addr, f)
}
