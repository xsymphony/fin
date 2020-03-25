package fin

import (
	"log"

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

	methodTrees methodTrees
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
	// 获取此方法下的根节点, 不存在则新建
	root := e.methodTrees.get(method)
	if root == nil {
		root = new(node)
		e.methodTrees = append(e.methodTrees, methodTree{method: method, root: root})
	}
	root.addRoute(path, h)
	log.Printf("[fin-debug]Register Method: %s | URL: %s", method, path)
}

func (e *Engine) Dispatch(fastCtx *fasthttp.RequestCtx) {
	uri := string(fastCtx.Path())
	method := string(fastCtx.Method())
	root := e.methodTrees.get(method)
	if root == nil {
		fastCtx.SetStatusCode(404)
		fastCtx.WriteString("404 NOT FOUND")
		return
	}
	ctx := &Context{
		RequestCtx: fastCtx,
		index:      -1,
		Params:     Params{},
	}
	value := root.getValue(uri, ctx.Params)
	if value.handlers == nil {
		fastCtx.SetStatusCode(404)
		fastCtx.WriteString("404 NOT FOUND")
		return
	}
	ctx.chain = value.handlers
	ctx.Params = value.params
	ctx.Next()
}

func (e *Engine) Run(addr string) error {
	server := &fasthttp.Server{
		Handler: e.Dispatch,
	}
	e.server = server
	return server.ListenAndServe(addr)
}

func (e *Engine) Shutdown() error {
	return e.server.Shutdown()
}
