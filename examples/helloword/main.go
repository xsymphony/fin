package main

import (
	"fmt"

	"github.com/xsymphony/fin"
)

func main() {
	beforeFunc := func(ctx *fin.Context) {
		fmt.Println("start serve request: ", string(ctx.Path()))
	}
	afterFunc := func(ctx *fin.Context) {
		fmt.Println("after serve request: ", string(ctx.Path()))
	}

	r := fin.New()
	{
		api := r.Group("/api")
		{
			v1 := api.Group("/v1")
			{
				v1.GET("/hello", beforeFunc, func(ctx *fin.Context) {
					ctx.WriteString("hello world")
				}, afterFunc)
			}
			v2 := api.Group("/v2")
			{
				v2.GET("/hello", beforeFunc, func(ctx *fin.Context) {
					ctx.WriteString("你好")
				}, afterFunc)
			}
		}
	}
	r.Run(":8080")
}
