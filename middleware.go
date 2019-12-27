package fin

import (
    "fmt"
    "os"
    "runtime"
    "time"

    "github.com/valyala/fasthttp"
)

type Middleware func(pre Handler) (wrapped Handler)

type AllowMethodHandler struct {
    methods []string
    next    Handler
}

func (amh *AllowMethodHandler) isAllowMethod(method string) bool {
    for _, m := range amh.methods {
        if m == method {
            return true
        }
    }

    return false
}

func (amh *AllowMethodHandler) Serve(ctx *fasthttp.RequestCtx) {
    if !amh.isAllowMethod(string(ctx.Method())) {
        ctx.SetStatusCode(405)
        ctx.WriteString("405 METHOD NOT ALLOWED")
        return
    }
    amh.next.Serve(ctx)
}

// AllowMethodMiddleware 借助自定义结构体实现限制请求方法的中间件
func AllowMethodMiddleware(method ...string) Middleware {
    return func(pre Handler) (wrapped Handler) {
        h := &AllowMethodHandler{
            methods: method,
            next: pre,
        }

        return h
    }
}

// AllowMethodMiddleware 借助adapter实现recover中间件
func Recovery() Middleware {
    return func(pre Handler) (wrapped Handler) {
        h := NewAdapter(func(ctx *fasthttp.RequestCtx) {
            defer func() {
                if e := recover(); e != nil {
                    buf := make([]byte, 64<<10)
                    buf = buf[:runtime.Stack(buf, false)]
                    msg := fmt.Sprintf("[%s] When server request %s panic: %s\n ==> %s\n",
                        time.Now().Format("2006-01-02 15:04:05"), string(ctx.Path()), e, buf)
                    _, _ = fmt.Fprintln(os.Stderr, msg)

                    ctx.WriteString("500 Internal Serve Error")
                }

            }()

            pre.Serve(ctx)
        })

        return h
    }
}