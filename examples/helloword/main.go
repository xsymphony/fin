package main

import (
	"github.com/valyala/fasthttp"
	"github.com/xsymphony/fin"
)

func main() {
	r := fin.New()
	{
		api := r.Group("/api")
		{
			v1 := api.Group("/v1")
			{
				v1.Handle("/hello", "GET", func(ctx *fasthttp.RequestCtx) {
					ctx.WriteString("hello world")
				})
			}
			v2 := api.Group("/v2")
			{
				v2.Handle("/hello", "GET", func(ctx *fasthttp.RequestCtx) {
					ctx.WriteString("你好")
				})
			}
		}
	}
	r.Run(":8080")
}
