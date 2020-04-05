package fin

import (
	"log"
	"os"

	"github.com/valyala/fasthttp"
)

type HandlerFunc func(ctx *Context)

type Engine struct {
	Router

	server *fasthttp.Server

	methodTrees methodTrees

	HandleNotFound HandlerFunc
}

func New(handlers ...HandlerFunc) *Engine {
	engine := &Engine{
		Router: Router{
			path: "",
		},
		HandleNotFound: func(ctx *Context) {
			ctx.String(fasthttp.StatusNotFound, "404 NOT FOUND")
			return
		},
	}
	engine.server = &fasthttp.Server{
		Handler: engine.Dispatch,
		Name: "fin",
	}
	engine.Router.engine = engine
	engine.Router.middlewares = append(engine.Router.middlewares, handlers...)
	return engine
}

func (e *Engine) Apply(options ...Option) {
	for _, option := range options {
		option(e)
	}
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
	ctx := &Context{
		RequestCtx: fastCtx,
		index:      -1,
		Params:     Params{},
	}
	root := e.methodTrees.get(method)
	if root == nil {
		e.HandleNotFound(ctx)
		return
	}
	value := root.getValue(uri, ctx.Params)
	if value.handlers == nil {
		e.HandleNotFound(ctx)
		return
	}
	ctx.chain = value.handlers
	ctx.Params = value.params
	ctx.Next()
}

func (e *Engine) Run(addr string) error {
	return e.server.ListenAndServe(addr)
}

func (e *Engine) RunUnix(addr string, mode os.FileMode) error {
	return e.server.ListenAndServeUNIX(addr, mode)
}

func (e *Engine) RunTLS(addr, certFile, keyFile string) error {
	return e.server.ListenAndServeTLS(addr, certFile, keyFile)
}

func (e *Engine) Shutdown() error {
	return e.server.Shutdown()
}
