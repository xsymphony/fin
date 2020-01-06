package fin

type IRouter interface {
	Handle(string, string, ...HandlerFunc)
	ANY(string, ...HandlerFunc)
	GET(string, ...HandlerFunc)
	POST(string, ...HandlerFunc)
	DELETE(string, ...HandlerFunc)
	PATCH(string, ...HandlerFunc)
	PUT(string, ...HandlerFunc)
	OPTIONS(string, ...HandlerFunc)
	HEAD(string, ...HandlerFunc)
}

type Router struct {
	engine *Engine
	path   string
}

// 检查Router是否实现了IRouter接口
var _ IRouter = &Router{}

// handle 注册一个新的路由函数
func (r *Router) handle(relativePath string, method string, h ...HandlerFunc) {
	// 计算路由的绝对路径
	path := r.path + relativePath
	// 注册路由
	r.engine.addRoute(path, method, h...)
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
func (r *Router) Group(relativePath string) *Router {
	// 计算路由的绝对路径
	path := r.path + relativePath
	router := &Router{
		path:   path,
		engine: r.engine,
	}

	return router
}
