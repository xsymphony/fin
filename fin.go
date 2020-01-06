package fin

import (
	"github.com/valyala/fasthttp"
)

type HandlerFunc func(ctx *Context)

type IEngine interface {
	IRouter

	Run(string) error
	Shutdown() error
}

type Engine struct {
	Router

	server *fasthttp.Server

	handlers trees
}

var _ IEngine = &Engine{}

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

func (e *Engine) dispatch(fastCtx *fasthttp.RequestCtx) {
	ctx := &Context{
		RequestCtx: fastCtx,
		Response:   &fastCtx.Response,
		Request:    &fastCtx.Request,
	}
	uri := string(ctx.Path())
	method := string(ctx.Method())
	handlers := e.handlers.get(method)
	if handlers == nil {
		fastCtx.SetStatusCode(404)
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
	server := &fasthttp.Server{
		Handler: e.dispatch,
	}
	e.server = server
	return server.ListenAndServe(addr)
}

func (e *Engine) Shutdown() error {
	return e.server.Shutdown()
}
