package goweb_test

import (
	"testing"

	"github.com/techiehe/goweb"
)

// 测试 goweb
func TestGoweb(t *testing.T) {
	gw := goweb.New()
	gw.Run(":8080")
}
