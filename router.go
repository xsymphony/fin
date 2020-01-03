package fin

import (
	"fmt"

	"github.com/valyala/fasthttp"
)

type Router struct {
	engine   *Engine
	prefix   string                             // prefix是整个路由的公共前缀
	handlers map[string]fasthttp.RequestHandler // 存放所有handle函数与url的映射关系
}

// AddRouter 注册一个新的路由函数
func (r *Router) AddRouter(uri string, h fasthttp.RequestHandler) {
	uri = r.prefix + uri
	if _, ok := r.handlers[uri]; ok {
		panic(fmt.Sprintf("duplicate uri %s", uri))
	}
	r.handlers[uri] = h
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
