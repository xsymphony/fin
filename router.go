package fin

import (
	"fmt"

	"github.com/valyala/fasthttp"
)

type Router struct {
	engine   *Engine
	path     string                             // prefix是整个路由的公共前缀
	handlers map[string]fasthttp.RequestHandler // 存放所有handle函数与url的映射关系
}

// AddRouter 注册一个新的路由函数
func (r *Router) AddRouter(relativePath string, h fasthttp.RequestHandler) {
	// 计算路由的绝对路径
	path := r.path + relativePath
	// 从根router中查找路由是否存在
	if _, ok := r.engine.handlers[path]; ok {
		panic(fmt.Sprintf("duplicate uri %s", path))
	}
	// 注册路由函数到根router中
	r.engine.handlers[path] = h
}

// Group新建一个路由分组
func (r *Router) Group(relativePath string) *Router {
	// 计算路由的绝对路径
	path := r.path + relativePath
	router := &Router{
		path:     path,
		handlers: make(map[string]fasthttp.RequestHandler),
		engine:   r.engine,
	}

	return router
}

// HandleRequest 作为fasthttp总的入口函数，进行路由分发
func (r *Router) HandleRequest(ctx *fasthttp.RequestCtx) {
	uri := string(ctx.Path())
	h, ok := r.handlers[uri]
	if !ok {
		ctx.SetStatusCode(404)
		ctx.WriteString("404 NOT FOUND")
		return
	}
	h(ctx)
}
