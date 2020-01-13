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

	trees trees
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

func (e *Engine) addRoute(path string, method string, h ...HandlerFunc) {
	// 获取此方法下的所有路由函数map，不存在则新建
	handlers := e.trees.get(method)
	if handlers == nil {
		handlers = make(node)
		e.trees = append(e.trees, tree{method: method, node: handlers})
	}
	handlers.addRoute(path, h...)
}

func (e *Engine) dispatch(fastCtx *fasthttp.RequestCtx) {
	uri := string(fastCtx.Path())
	method := string(fastCtx.Method())
	tree := e.trees.get(method)
	if tree == nil {
		fastCtx.SetStatusCode(404)
		fastCtx.WriteString("404 NOT FOUND")
		return
	}
	chain, ok := tree[uri]
	if !ok {
		fastCtx.SetStatusCode(404)
		fastCtx.WriteString("404 NOT FOUND")
		return
	}
	ctx := &Context{
		RequestCtx: fastCtx,
		chain:      chain,
		index:      -1,
	}
	ctx.Next()
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
