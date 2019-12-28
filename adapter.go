package fin

import "github.com/valyala/fasthttp"

// HandlerFunc是通用的handler函数适配器
// 使外层调用不再用申明结构体
type HandlerFunc func(ctx *fasthttp.RequestCtx)

func (f HandlerFunc) Serve(ctx *fasthttp.RequestCtx) {
	f(ctx)
}
