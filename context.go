package fin

import (
	"encoding/json"
	"fmt"

	"github.com/valyala/fasthttp"
)

const abortIndex = 2<<6 - 1

type Context struct {
	*fasthttp.RequestCtx

	index int8
	chain []HandlerFunc

	Params Params
}

func (c *Context) reset() {
	c.index = -1
	c.chain = nil
	c.Params = c.Params[0:0]
}

func (c *Context) Next() {
	c.index++
	for s := int8(len(c.chain)); c.index < s; c.index++ {
		c.chain[c.index](c)
	}
}

func (c *Context) Abort() {
	c.index = abortIndex
}

func (c *Context) Param(key string) string {
	return c.Params.ByName(key)
}

func (c *Context) Query(key string) (string, bool) {
	v := c.QueryArgs().Peek(key)
	if v == nil {
		return "", false
	}

	return string(v), true
}

func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetStatusCode(code)
	if len(values) > 0 {
		format = fmt.Sprintf(format, values...)
	}
	if _, err := c.WriteString(format); err != nil {
		panic(err)
	}
}

func (c *Context) JSON(code int, obj interface{}) {
	c.Response.Header.Set("Content-type", "application/json")
	d, err := json.Marshal(obj)
	if err != nil {
		return
	}
	if _, err := c.Write(d); err != nil {
		panic(err)
	}
}

func (c *Context) JSONAbort(code int, obj interface{}) {
	c.JSON(code, obj)
	c.Abort()
}
