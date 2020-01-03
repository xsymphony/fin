package fin

import "github.com/valyala/fasthttp"

type Engine struct {
	Router
}

func New() *Engine {
	engine := &Engine{
		Router: Router{
			path:     "",
			handlers: make(map[string]fasthttp.RequestHandler),
		},
	}
	engine.Router.engine = engine
	return engine
}

func (e *Engine) Run(addr string) error {
	return fasthttp.ListenAndServe(addr, e.HandleRequest)
}
