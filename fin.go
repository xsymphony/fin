package fin

import (
	"fmt"

	"github.com/valyala/fasthttp"
)

type Engine struct {
	Router

	handlers map[string]fasthttp.RequestHandler
}

func New() *Engine {
	engine := &Engine{
		Router: Router{
			path: "",
		},
		handlers: make(map[string]fasthttp.RequestHandler),
	}
	engine.Router.engine = engine
	return engine
}

func (e *Engine) addRouter(path string, h fasthttp.RequestHandler) {
	// 从根router中查找路由是否存在
	if _, ok := e.handlers[path]; ok {
		panic(fmt.Sprintf("duplicate uri %s", path))
	}
	// 注册路由函数到根router中
	e.handlers[path] = h
}

func (e *Engine) dispatch(ctx *fasthttp.RequestCtx) {
	uri := string(ctx.Path())
	h, ok := e.handlers[uri]
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
