package fin

import (
	"github.com/valyala/fasthttp"
)

type Router struct {
	engine *Engine
	path   string
}

// AddRouter 注册一个新的路由函数
func (r *Router) AddRouter(relativePath string, h fasthttp.RequestHandler) {
	// 计算路由的绝对路径
	path := r.path + relativePath
	// 注册路由
	r.engine.addRouter(path, h)
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
