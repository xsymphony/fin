package fin

import "github.com/valyala/fasthttp"

const abortIndex = 2<<6 - 1

type Context struct {
	*fasthttp.RequestCtx

	index int8
	chain []HandlerFunc

	Params   Params
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
