package main

import (
    "github.com/valyala/fasthttp"
    "github.com/xsymphony/fin"
)

func main() {
    r := fin.Default("/api")
    r.Use(fin.AllowMethodMiddleware("GET"))
    {
        r.ANY("/v1/hello", fin.NewAdapter(func(ctx *fasthttp.RequestCtx) {
            ctx.WriteString("hello world")
        }))
        r.ANY("/v2/hello", fin.NewAdapter(func(ctx *fasthttp.RequestCtx) {
            ctx.WriteString("你好")
        }))
        r.GET("/panic", fin.NewAdapter(func(ctx *fasthttp.RequestCtx) {
            panic("just for panic")
        }))
    }
    fasthttp.ListenAndServe(":8080", r.HandleRequest)
}
