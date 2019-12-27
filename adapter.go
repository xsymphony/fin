package fin

import "github.com/valyala/fasthttp"

// HandlerAdapter是通用的handler函数适配器
// 使外层调用不再用申明结构体
type HandlerAdapter struct {
    f func(ctx *fasthttp.RequestCtx)
}

func (adapter *HandlerAdapter) Serve(ctx *fasthttp.RequestCtx) {
    adapter.f(ctx)
}

func NewAdapter(f func(ctx *fasthttp.RequestCtx)) *HandlerAdapter {
    return &HandlerAdapter{
        f: f,
    }
}
