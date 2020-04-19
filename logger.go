package fin

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/valyala/fasthttp"
)

var logger *log.Logger

func SimpleLogger() HandlerFunc {
	return func(c *Context) {
		if logger == nil {
			logger = log.New(os.Stdout, "[fin]", 0)
		}
		start := time.Now()
		c.Next()
		end := time.Now()
		logger.Printf("%v | %3d | %8v | %15s |%-7s %#v\n",
			end.Format("2006/01/02 - 15:04:05"),
			c.Response.StatusCode(),
			end.Sub(start),
			c.RemoteIP(),
			string(c.Method()),
			string(c.Path()),
		)
	}
}

func Logger() HandlerFunc {
	return func(c *Context) {
		start := time.Now()
		c.Next()
		end := time.Now()
		kvs := []interface{}{
			"method", string(c.Method()),
			"path", string(c.Path()),
			"time", end.Format("2006/01/02 15:04:05"),
			"status", c.Response.StatusCode(),
			"cost", end.Sub(start),
			"remote_ip", c.RemoteIP(),
			"request_body", string(c.PostBody()),
		}
		contentType := string(c.Response.Header.Peek(fasthttp.HeaderContentType))
		if strings.Contains(contentType, "application/json") {
			kvs = append(kvs, "response", string(c.Response.Body()))
		}
		c.engine.logger.Infow("http request with", kvs...)
	}
}
