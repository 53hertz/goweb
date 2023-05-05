//go:build e2e

// 集成测试
package web

import (
	"fmt"
	"testing"
)

func TestServer(t *testing.T) {
	h := &HttpServer{}

	h.Get("/demo", func(ctx Context) {
		fmt.Println("xxx")
	})

	h.Start(":8081")
}
