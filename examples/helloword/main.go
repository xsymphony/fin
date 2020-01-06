package main

import (
    "github.com/valyala/fasthttp"
    "github.com/xsymphony/fin"
)

func main() {
    r := fin.NewRouter("/api")
    {
        r.AddRouter("/v1/hello", func(ctx *fasthttp.RequestCtx) {
            ctx.WriteString("hello world")
        })
        r.AddRouter("/v2/hello",func(ctx *fasthttp.RequestCtx) {
            ctx.WriteString("你好")
        })
    }
    fasthttp.ListenAndServe(":8080", r.HandleRequest)
}
