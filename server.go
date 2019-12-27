package main

import "github.com/valyala/fasthttp"


// dispatch 是路由分发的原型
// 根据请求路由执行不同的处理逻辑
func dispatch(ctx *fasthttp.RequestCtx) {
    switch string(ctx.Path()) {
    case "/api/v1/hello":
        ctx.WriteString("hello world")
    case "/api/v2/hello":
        ctx.WriteString("你好")
    default:
        ctx.SetStatusCode(404)
        ctx.WriteString("404 NOT FOUND")
    }
}

func main() {
    fasthttp.ListenAndServe(":8080", dispatch)
}
