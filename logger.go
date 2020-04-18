package fin

import (
	"log"
	"os"
	"time"
)

var logger *log.Logger

func Logger() HandlerFunc {
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
