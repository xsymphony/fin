package fin

import (
	"github.com/valyala/fasthttp"
)

type HandlerFunc fasthttp.RequestHandler

type Engine struct {
	Router

	handlers trees
}

func New() *Engine {
	engine := &Engine{
		Router: Router{
			path: "",
		},
	}
	engine.Router.engine = engine
	return engine
}

func (e *Engine) addRoute(path string, method string, h HandlerFunc) {
	// 获取此方法下的所有路由函数map，不存在则新建
	handlers := e.handlers.get(method)
	if handlers == nil {
		handlers = make(node)
		e.handlers = append(e.handlers, tree{method: method, node: handlers})
	}
	handlers.addRoute(path, h)
}

func (e *Engine) dispatch(ctx *fasthttp.RequestCtx) {
	uri := string(ctx.Path())
	method := string(ctx.Method())
	handlers := e.handlers.get(method)
	if handlers == nil {
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
