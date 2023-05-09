package web

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"reflect"
	"testing"
)

func TestRouter_addRoute(t *testing.T) {
	// 构建注册路由
	testRoutes := []struct {
		method string
		path   string
	}{
		{
			method: http.MethodGet,
			path:   "/",
		},
		{
			method: http.MethodGet,
			path:   "/user",
		},
		{
			method: http.MethodGet,
			path:   "/user/home",
		},
		{
			method: http.MethodGet,
			path:   "/order/detail",
		},
		{
			method: http.MethodPost,
			path:   "/order/create",
		},
		{
			method: http.MethodPost,
			path:   "/login",
		},
	}

	mockHandler := func(ctx *Context) {}

	r := newRouter()

	for _, route := range testRoutes {
		r.addRoute(route.method, route.path, mockHandler)
	}

	wantRouter := &router{
		trees: map[string]*node{
			http.MethodGet: &node{
				path:    "/",
				handler: mockHandler,
				children: map[string]*node{
					"user": &node{
						path:    "user",
						handler: mockHandler,
						children: map[string]*node{
							"home": &node{
								path:    "home",
								handler: mockHandler,
							},
						},
					},
					"order": &node{
						path: "order",
						children: map[string]*node{
							"detail": &node{
								path:    "detail",
								handler: mockHandler,
							},
						},
					},
				},
			},
			http.MethodPost: &node{
				path: "/",
				children: map[string]*node{
					"order": &node{
						path: "order",
						children: map[string]*node{
							"create": &node{
								path:    "create",
								handler: mockHandler,
							},
						},
					},
					"login": &node{
						path:    "login",
						handler: mockHandler,
					},
				},
			},
		},
	}

	msg, ok := wantRouter.equal(r)
	assert.True(t, ok, msg)

	assert.PanicsWithValue(t, "web: 路由是空字符串", func() {
		r.addRoute(http.MethodGet, "", mockHandler)
	})
	assert.PanicsWithValue(t, "web: 路由必须以 / 开头", func() {
		r.addRoute(http.MethodGet, "a/b/c", mockHandler)
	})
	assert.PanicsWithValue(t, "web: 路由不能以 / 结尾", func() {
		r.addRoute(http.MethodGet, "/a/b/c/", mockHandler)
	})
	assert.PanicsWithValue(t, "web: 不能有连续的 /", func() {
		r.addRoute(http.MethodGet, "/a//b/c", mockHandler)
	})

	assert.PanicsWithValue(t, "web: 路由冲突[/]", func() {
		r.addRoute(http.MethodGet, "/", mockHandler)
	})
	assert.PanicsWithValue(t, "web: 路由重复注册", func() {
		r.addRoute(http.MethodGet, "/user/home", mockHandler)
	})
}

func TestRouter_findRoute(t *testing.T) {
	// 构建注册路由
	testRoutes := []struct {
		method string
		path   string
	}{
		{
			method: http.MethodGet,
			path:   "/",
		},
		{
			method: http.MethodGet,
			path:   "/user",
		},
		//{
		//	method: http.MethodPost,
		//	path:   "/order/create",
		//},
	}

	r := newRouter()
	mockHandler := func(ctx *Context) {}
	for _, route := range testRoutes {
		r.addRoute(route.method, route.path, mockHandler)
	}

	// 构建查找路由
	testCases := []struct {
		name     string
		method   string
		path     string
		found    bool
		wantNode *node
	}{
		{
			name:   "method not found",
			method: http.MethodHead,
		},
		{
			name:   "path not found",
			method: http.MethodGet,
			path:   "/a/b/c",
		},
		{
			name:   "root",
			method: http.MethodGet,
			path:   "/",
			found:  true,
			wantNode: &node{
				path:    "/",
				handler: mockHandler,
				children: map[string]*node{
					"user": &node{
						path:    "user",
						handler: mockHandler,
					},
				},
			},
		},
		{
			name:   "user",
			method: http.MethodGet,
			path:   "user",
			found:  true,
			wantNode: &node{
				path:    "user",
				handler: mockHandler,
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			node, found := r.findRoute(testCase.method, testCase.path)
			assert.Equal(t, testCase.found, found)
			if !found {
				return
			}
			fmt.Println(testCase.path, node.path)
			assert.Equal(t, testCase.path, node.path)
			//assert.Equal(t, testCase.wantNode.children, node.children)
			msg, ok := testCase.wantNode.equal(node)
			assert.True(t, ok, msg)
		})

	}
}

func (r *router) equal(y router) (string, bool) {
	for k, v := range r.trees {
		dst, ok := y.trees[k]
		if !ok {
			return fmt.Sprintf("找不到对应 http method"), false
		}
		msg, equal := v.equal(dst)
		if !equal {
			return msg, false
		}
	}
	return "", true
}

func (n *node) equal(y *node) (string, bool) {
	if n.path != y.path {
		return fmt.Sprintf("路径不匹配"), false
	}
	if len(n.children) != len(y.children) {
		return fmt.Sprintf("子节点数量不相等"), false
	}

	nHandler := reflect.ValueOf(n.handler)
	yHandler := reflect.ValueOf(y.handler)
	if nHandler != yHandler {
		return fmt.Sprintf("handler 不匹配"), false
	}

	for path, child := range n.children {
		dst, ok := y.children[path]
		if !ok {
			return fmt.Sprintf("子节点路由不匹配"), false
		}
		msg, equal := child.equal(dst)
		if !equal {
			return msg, false
		}
	}
	return "", true
}
