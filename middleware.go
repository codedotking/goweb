package main

import "net/http"

type MiddlewareFunc func(http.Handler) http.Handler

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 记录日志
		next.ServeHTTP(w, r)
	})
}
