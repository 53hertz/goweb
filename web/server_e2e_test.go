//go:build e2e

// 集成测试
package web

import (
	"net/http"
	"testing"
)

func TestServer(t *testing.T) {
	h := NewHttpServer()
	h.addRoute(http.MethodGet, "/user", func(ctx *Context) {
		ctx.Resp.Write([]byte("hello gopher"))
	})

	h.Start(":8081")

}
