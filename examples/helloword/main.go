package main

import (
	"github.com/valyala/fasthttp"
	"github.com/xsymphony/fin"
)

func main() {
	r := fin.New()
	{
		r.AddRouter("/api/v1/hello", func(ctx *fasthttp.RequestCtx) {
			ctx.WriteString("hello world")
		})
		r.AddRouter("/api/v2/hello", func(ctx *fasthttp.RequestCtx) {
			ctx.WriteString("你好")
		})
	}
	r.Run(":8080")
}
