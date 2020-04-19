package fin

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/valyala/fasthttp"
)

func Recovery() HandlerFunc {
	return func(c *Context) {
		defer func() {
			if e := recover(); e != nil {
				buf := make([]byte, 64<<10)
				buf = buf[:runtime.Stack(buf, false)]
				msg := fmt.Sprintf("[%s] when server %s panic: %s\n ==> %s\n",
					time.Now().Format("2006-01-02 15:04:05"), c.Path(), e, buf)
				_, _ = fmt.Fprintln(os.Stderr, msg)
				c.SetStatusCode(fasthttp.StatusInternalServerError)
				c.Abort()
				return
			}
		}()
		c.Next()
	}
}
