package fin

import (
	"path"
	"path/filepath"
	"strings"

	"github.com/valyala/fasthttp"
)

type IRouter interface {
	Use(middleware ...HandlerFunc)

	Handle(relativePath string, method string, handlers ...HandlerFunc)
	ANY(relativePath string, handlers ...HandlerFunc)
	GET(relativePath string, handlers ...HandlerFunc)
	POST(relativePath string, handlers ...HandlerFunc)
	DELETE(relativePath string, handlers ...HandlerFunc)
	PATCH(relativePath string, handlers ...HandlerFunc)
	PUT(relativePath string, handlers ...HandlerFunc)
	OPTIONS(relativePath string, handlers ...HandlerFunc)
	HEAD(relativePath string, handlers ...HandlerFunc)

	Static(relativePath string, root string)
}

type Router struct {
	path string

	engine *Engine

	middlewares []HandlerFunc
}

// 检查Router是否实现了IRouter接口
var _ IRouter = &Router{}

// handle 注册一个新的路由函数
func (r *Router) handle(relativePath string, method string, h ...HandlerFunc) {
	// 计算路由的绝对路径
	absPath := r.path + relativePath
	// 组合路由的中间件到handlers
	handlers := make([]HandlerFunc, len(r.middlewares)+len(h))
	copy(handlers, r.middlewares)
	copy(handlers[len(r.middlewares):], h)
	// 注册路由
	r.engine.addRoute(absPath, method, handlers...)
}

func (r *Router) Use(middlewares ...HandlerFunc) {
	r.middlewares = append(r.middlewares, middlewares...)
}

// Handle 注册一个新的路由函数
func (r *Router) Handle(relativePath string, method string, h ...HandlerFunc) {
	r.handle(relativePath, method, h...)
}

func (r *Router) GET(relativePath string, h ...HandlerFunc) {
	r.Handle(relativePath, "GET", h...)
}

func (r *Router) POST(relativePath string, h ...HandlerFunc) {
	r.Handle(relativePath, "POST", h...)
}

func (r *Router) DELETE(relativePath string, h ...HandlerFunc) {
	r.Handle(relativePath, "DELETE", h...)
}

func (r *Router) PUT(relativePath string, h ...HandlerFunc) {
	r.Handle(relativePath, "PUT", h...)
}

func (r *Router) PATCH(relativePath string, h ...HandlerFunc) {
	r.Handle(relativePath, "PATCH", h...)
}

func (r *Router) HEAD(relativePath string, h ...HandlerFunc) {
	r.Handle(relativePath, "HEAD", h...)
}

func (r *Router) OPTIONS(relativePath string, h ...HandlerFunc) {
	r.Handle(relativePath, "OPTIONS", h...)
}

func (r *Router) ANY(relativePath string, h ...HandlerFunc) {
	r.Handle(relativePath, "GET", h...)
	r.Handle(relativePath, "POST", h...)
	r.Handle(relativePath, "DELETE", h...)
	r.Handle(relativePath, "PUT", h...)
	r.Handle(relativePath, "PATCH", h...)
	r.Handle(relativePath, "HEAD", h...)
	r.Handle(relativePath, "OPTIONS", h...)
	r.Handle(relativePath, "CONNECT", h...)
	r.Handle(relativePath, "TRACE", h...)
}

// Group新建一个路由分组
func (r *Router) Group(relativePath string, handlers ...HandlerFunc) *Router {
	// 计算路由的绝对路径
	absPath := r.path + relativePath
	// 复制当前路由的中间件到下一级
	middleware := make([]HandlerFunc, len(r.middlewares)+len(handlers))
	copy(middleware[0:len(r.middlewares)], r.middlewares)
	copy(middleware[len(r.middlewares):], handlers)
	router := &Router{
		path:        absPath,
		engine:      r.engine,
		middlewares: middleware,
	}

	return router
}

func (r *Router) StaticFile(relativePath string, filepath string) {
	if strings.Contains(relativePath, ":") || strings.Contains(relativePath, "*") {
		panic("URL parameters can not be used when serving a static file")
	}
	handler := func(c *Context) {
		c.SendFile(filepath)
	}
	r.GET(relativePath, handler)
	r.HEAD(relativePath, handler)
}

func (r *Router) Static(relativePath string, dir string) {
	if strings.Contains(relativePath, ":") || strings.Contains(relativePath, "*") {
		panic("URL parameters can not be used when serving a static folder")
	}
	fs := &fasthttp.FS{
		Root:            dir,
		AcceptByteRange: true,
		Compress:        false,
	}
	handler := r.createStaticHandler(dir, fs)
	urlPattern := path.Join(relativePath, "/*filepath")
	// Register GET handlers
	r.GET(urlPattern, handler)
	r.HEAD(urlPattern, handler)
}

func (r *Router) createStaticHandler(dir string, fs *fasthttp.FS) HandlerFunc {
	h := fs.NewRequestHandler()
	return func(c *Context) {
		file := c.Param("filepath")
		if len(file) == 0 || file[0] != '/' {
			// extend relative path to absolute path
			var err error
			if _, err = filepath.Abs(file); err != nil {
				c.String(fasthttp.StatusNotFound, "NOT FOUND")
				return
			}
		}
		before := string(c.Path())
		c.Request.SetRequestURI(file)
		h(c.RequestCtx)
		c.Request.SetRequestURI(before)
	}
}
