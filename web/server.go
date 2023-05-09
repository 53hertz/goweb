package web

import (
	"net"
	"net/http"
)

type HandleFunc func(ctx *Context)

// 确保 HttpServer 一定实现了 Server
var _ Server = &HttpServer{}

type Server interface {
	http.Handler
	Start(addr string) error

	// AddRoute 路由注册
	addRoute(method string, path string, handleFunc HandleFunc)
}

type HttpServer struct {
	router
}

func NewHttpServer() *HttpServer {
	return &HttpServer{
		newRouter(),
	}
}

func (h *HttpServer) Get(path string, handleFunc HandleFunc) {
	h.addRoute(http.MethodGet, path, handleFunc)
}

// ServeHTTP 处理请求入口
func (h *HttpServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := &Context{
		Req:  r,
		Resp: w,
	}

	h.Serve(ctx)
}

func (h *HttpServer) Serve(ctx *Context) {
	node, ok := h.findRoute(ctx.Req.Method, ctx.Req.URL.Path)
	if !ok || node.handler == nil {
		ctx.Resp.WriteHeader(404)
		ctx.Resp.Write([]byte("not found"))
		return
	}
	node.handler(ctx)
}

func (h *HttpServer) Start(addr string) error {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	return http.Serve(l, h)
}
