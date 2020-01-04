package fin

import (
	"fmt"

	"github.com/valyala/fasthttp"
)

type Engine struct {
	Router

	handlers map[string]map[string]fasthttp.RequestHandler
}

func New() *Engine {
	engine := &Engine{
		Router: Router{
			path: "",
		},
		handlers: make(map[string]map[string]fasthttp.RequestHandler),
	}
	engine.Router.engine = engine
	return engine
}

func (e *Engine) addRoute(path string, method string, h fasthttp.RequestHandler) {
	// 获取此方法下的所有路由函数map，不存在则新建
	handlers, ok := e.handlers[method]
	if !ok {
		handlers = make(map[string]fasthttp.RequestHandler)
		e.handlers[method] = handlers
	}
	// 从根router中查找路由是否存在
	if _, ok := handlers[path]; ok {
		panic(fmt.Sprintf("duplicate uri %s", path))
	}
	// 注册路由函数到根router中
	handlers[path] = h
}

func (e *Engine) dispatch(ctx *fasthttp.RequestCtx) {
	uri := string(ctx.Path())
	method := string(ctx.Method())
	handlers, ok := e.handlers[method]
	if !ok {
		ctx.SetStatusCode(404)
		ctx.WriteString("404 NOT FOUND")
		return
	}
	h, ok := handlers[uri]
	if !ok {
		ctx.SetStatusCode(404)
		ctx.WriteString("404 NOT FOUND")
		return
	}
	h(ctx)
}

func (e *Engine) Run(addr string) error {
	return fasthttp.ListenAndServe(addr, e.dispatch)
}
