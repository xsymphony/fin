package fin

import (
	"html/template"
	"log"
	"os"
	"sync"

	"github.com/valyala/fasthttp"
)

type HandlerFunc func(ctx *Context)

type Engine struct {
	Router

	HandleNotFound HandlerFunc

	server        *fasthttp.Server
	methodTrees   methodTrees
	ctxPool       sync.Pool
	htmlTemplates *template.Template
	funcMap       template.FuncMap

	logger *log.Logger
}

func New() *Engine {
	engine := &Engine{
		Router: Router{
			path: "",
		},
		HandleNotFound: func(ctx *Context) {
			ctx.String(fasthttp.StatusNotFound, "404 NOT FOUND")
			return
		},
		funcMap: template.FuncMap{},
		logger:  log.New(os.Stderr, "[fin]", log.LstdFlags|log.LstdFlags),
	}
	engine.ctxPool.New = func() interface{} {
		return &Context{}
	}
	engine.server = &fasthttp.Server{
		Handler: engine.dispatch,
		Name:    "fin",
	}
	engine.Router.engine = engine
	return engine
}

func Default() *Engine {
	engine := New()
	engine.Use(Recovery(), Logger())
	return engine
}

func (e *Engine) Apply(options ...Option) {
	for _, option := range options {
		option(e)
	}
}

func (e *Engine) SetFuncMap(funcMap template.FuncMap) {
	e.funcMap = funcMap
}

func (e *Engine) LoadHTMLGlob(pattern string) {
	e.logger.Printf("set html template pattern: %s", pattern)
	e.htmlTemplates = template.Must(template.New("").Funcs(e.funcMap).ParseGlob(pattern))
}

func (e *Engine) addRoute(path string, method string, h ...HandlerFunc) {
	// 获取此方法下的根节点, 不存在则新建
	root := e.methodTrees.get(method)
	if root == nil {
		root = new(node)
		e.methodTrees = append(e.methodTrees, methodTree{method: method, root: root})
	}
	root.addRoute(path, h)
	e.logger.Printf("register Method: %-6s URL: %s", method, path)
}

func (e *Engine) dispatch(fastCtx *fasthttp.RequestCtx) {
	ctx := e.ctxPool.Get().(*Context)
	ctx.engine = e
	ctx.RequestCtx = fastCtx
	ctx.reset()

	e.handleHTTPRequest(ctx)

	e.ctxPool.Put(ctx)
}

func (e *Engine) handleHTTPRequest(ctx *Context) {
	uri := string(ctx.RequestCtx.Path())
	method := string(ctx.RequestCtx.Method())
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
	e.logger.Printf("run %s", addr)
	return e.server.ListenAndServe(addr)
}

func (e *Engine) RunUnix(addr string, mode os.FileMode) error {
	e.logger.Printf("runUnix %s", addr)
	return e.server.ListenAndServeUNIX(addr, mode)
}

func (e *Engine) RunTLS(addr, certFile, keyFile string) error {
	e.logger.Printf("runTLS %s", addr)
	return e.server.ListenAndServeTLS(addr, certFile, keyFile)
}

func (e *Engine) Shutdown() error {
	e.logger.Println("receive to shutdown")
	return e.server.Shutdown()
}
