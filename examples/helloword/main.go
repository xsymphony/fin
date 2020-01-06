package main

import (
	"github.com/xsymphony/fin"
)

func main() {
	r := fin.New()
	{
		api := r.Group("/api")
		{
			v1 := api.Group("/v1")
			{
				v1.GET("/hello", func(ctx *fin.Context) {
					ctx.WriteString("hello world")
				})
			}
			v2 := api.Group("/v2")
			{
				v2.GET("/hello", func(ctx *fin.Context) {
					ctx.WriteString("你好")
				})
			}
		}
	}
	r.Run(":8080")
}
