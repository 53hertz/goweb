package web

import (
	"strings"
)

type router struct {
	trees map[string]*node
}

type node struct {
	path     string
	children map[string]*node
	handler  HandleFunc
}

func newRouter() router {
	return router{
		trees: map[string]*node{},
	}
}

// get /user/home  "" user home
func (r *router) addRoute(method string, path string, handleFunc HandleFunc) {
	//if r.trees == nil {
	//	r.trees = map[string]*node{}
	//}
	if path == "" {
		panic("web: 路由是空字符串")
	}
	if path[0] != '/' {
		panic("web: 路由必须以 / 开头")
	}
	if path != "/" && path[len(path)-1] == '/' {
		panic("web: 路由不能以 / 结尾")
	}
	root, ok := r.trees[method]
	if !ok {
		root = &node{
			path: "/",
		}
		r.trees[method] = root
	}
	if path == "/" {
		if root.handler != nil {
			panic("web: 路由冲突[/]")
		}
		root.handler = handleFunc
		return
	}

	segs := strings.Split(path[1:], "/")
	for _, seg := range segs {
		if seg == "" {
			panic("web: 不能有连续的 /")
		}
		child := root.childOrCreate(seg)
		root = child
	}
	if root.handler != nil {
		panic("web: 路由重复注册")
	}
	root.handler = handleFunc
}

func (n *node) childOrCreate(seg string) *node {
	if n.children == nil {
		n.children = map[string]*node{}
	}

	child, ok := n.children[seg]
	if !ok {
		child = &node{
			path: seg,
		}
		n.children[seg] = child
	}
	return child
}

func (r *router) findRoute(method string, path string) (*node, bool) {
	root, ok := r.trees[method]
	if !ok {
		return nil, false
	}
	if path == "/" {
		return root, true
	}

	segs := strings.Split(strings.Trim(path, "/"), "/")
	for _, seg := range segs {
		root, ok = root.childOf(seg)
		if !ok {
			return nil, false
		}
	}
	return root, true
}

func (n *node) childOf(path string) (*node, bool) {
	if n.children == nil {
		return nil, false
	}

	child, ok := n.children[path]
	return child, ok
}
