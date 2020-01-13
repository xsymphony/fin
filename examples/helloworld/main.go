package main

import (
	"fmt"
	"time"

	"github.com/xsymphony/fin"
)

func main() {
	hookFunc := func(ctx *fin.Context) {
		// 在执行下个handler函数之前打印请求信息
		fmt.Println("start serve request: ", string(ctx.Path()))
		// 调用Next()执行后面的handler函数
		ctx.Next()
		// 后面的handler函数执行完毕后打印请求信息
		fmt.Println("after serve request: ", string(ctx.Path()))
	}

	timedFunc := func(ctx *fin.Context) {
		fmt.Println("[timedFunc]start")
		start := time.Now()
		ctx.Next()
		fmt.Printf("[timedFunc]url: %s, used: %d\n", string(ctx.Path()), time.Now().Sub(start))
	}

	r := fin.New()
	r.Use(hookFunc, timedFunc)
	{
		api := r.Group("/api")
		{
			v1 := api.Group("/v1")
			v1.Use(func(ctx *fin.Context) {
				fmt.Println("I'm only used in v1")
				ctx.Next()
			})
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
