package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/xsymphony/fin"
)

var (
	logger = func(ctx *fin.Context) {
		// 在执行下个handler函数之前打印请求信息
		fmt.Println("start serve request: ", string(ctx.Path()))
		// 调用Next()执行后面的handler函数
		ctx.Next()
		// 后面的handler函数执行完毕后打印请求信息
		fmt.Println("after serve request: ", string(ctx.Path()))
	}

	timer = func(ctx *fin.Context) {
		fmt.Println("[timedFunc]start")
		start := time.Now()
		ctx.Next()
		fmt.Printf("[timedFunc]url: %s, used: %d\n", string(ctx.Path()), time.Now().Sub(start))
	}
)

func main() {
	r := fin.New()
	r.Use(logger, timer)
	r.Apply(fin.HandleNotFound(func(c *fin.Context) {
		c.HTML(http.StatusNotFound, "404.html", map[string]interface{}{
			"path": string(c.Path()),
		})
	}))
	r.LoadHTMLGlob("templates/*")
	{
		r.Static("/assets", "./static")
		r.GET("/index", func(c *fin.Context) {
			c.HTML(http.StatusOK, "index.html", nil)
		})
		v1 := r.Group("/api")
		{
			v1.GET("/echo/:name", func(c *fin.Context) {
				name := c.Param("name")
				age, _ := c.Query("age")
				c.JSON(http.StatusOK, map[string]interface{}{
					"name": name,
					"age":  age,
				})
			})
		}
	}
	r.Run(":8080")
}
