package fin

import (
    "fmt"

    "github.com/valyala/fasthttp"
)

type Handler interface {
    Serve(ctx *fasthttp.RequestCtx)
}

type Router struct {
    prefix     string               // prefix是整个路由的公共前缀
    handlers   map[string]Handler   // 存放所有handle函数与url的映射关系
    middleware []Middleware         // 存放所有中间件
}

func New(prefix string) *Router {
    return &Router{
        prefix: prefix,
        handlers: make(map[string]Handler),
    }
}

func Default(prefix string) *Router {
    r := New(prefix)
    r.Use(Recovery())
    return r
}

func (r *Router) Use(m ...Middleware) {
    r.middleware = append(r.middleware, m...)
}

// AddRouter 注册一个新的路由函数
func (r *Router) AddRouter(uri string, h Handler) {
    uri = r.prefix + uri
    if _, ok := r.handlers[uri]; ok {
        panic(fmt.Sprintf("duplicate uri %s", uri))
    }
    // 添加中间件
    for _, m := range r.middleware {
        h = m(h)
    }
    r.handlers[uri] = h
}

func (r *Router) ANY(uri string, h Handler) {
    r.AddRouter(uri, h)
}

func (r *Router) GET(uri string, h Handler) {
    h = AllowMethodMiddleware("GET")(h)
    r.AddRouter(uri, h)
}

func (r *Router) POST(uri string, h Handler) {
    h = AllowMethodMiddleware("POST")(h)
    r.AddRouter(uri, h)
}

func (r *Router) GETPOST(uri string, h Handler) {
    h = AllowMethodMiddleware("GET", "POST")(h)
    r.AddRouter(uri, h)
}

// HandleRequest 作为fasthttp总的入口函数，进行路由分发
func (r *Router) HandleRequest(ctx *fasthttp.RequestCtx) {
    uri := string(ctx.Path())
    h, ok := r.handlers[uri]
    if !ok {
        ctx.SetStatusCode(404)
        ctx.WriteString("404 NOT FOUND")
    }
    h.Serve(ctx)
}
