package fin

import "github.com/valyala/fasthttp"

type Context struct {
	*fasthttp.RequestCtx
	Request  *fasthttp.Request
	Response *fasthttp.Response
}
