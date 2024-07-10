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
		router: router.NewRouter(),
	}
}

func (f *Framework) GET(path string, handler http.HandlerFunc) {
	f.router.AddRoute("GET", path, handler)
}

func (f *Framework) POST(path string, handler http.HandlerFunc) {
	f.router.AddRoute("POST", path, handler)
}

func (f *Framework) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f.router.ServeHTTP(w, r)
}
