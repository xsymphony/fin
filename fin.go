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

	methodTrees trees
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
	// 获取此方法下的所有路由函数树, 不存在则新建
	tree := e.methodTrees.get(method)
	if tree == nil {
		tree = newTree(method)
		e.methodTrees = append(e.methodTrees, tree)
	}
	tree.addRoute(path, h...)
	log.Printf("[fin-debug]Register Method: %s | URL: %s", method, path)
}

func (e *Engine) Dispatch(fastCtx *fasthttp.RequestCtx) {
	uri := string(fastCtx.Path())
	method := string(fastCtx.Method())
	tree := e.methodTrees.get(method)
	if tree == nil {
		fastCtx.SetStatusCode(404)
		fastCtx.WriteString("404 NOT FOUND")
		return
	}
	node := tree.search(uri)
	if node == nil {
		fastCtx.SetStatusCode(404)
		fastCtx.WriteString("404 NOT FOUND")
		return
	}
	ctx := &Context{
		RequestCtx: fastCtx,
		chain:      node.chain,
		index:      -1,
	}
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
