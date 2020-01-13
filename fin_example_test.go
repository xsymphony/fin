package fin_test

import (
	"io"
	"net/http"
	"testing"

	"github.com/xsymphony/fin"
)

func TestNewRouter(t *testing.T) {
	r := fin.New()
	r.ANY("/hello", func(ctx *fin.Context) {
		ctx.WriteString("hello")
	})
	go func() {
		r.Run(":8080")
	}()

	resp, err := http.Get("http://127.0.0.1:8080/hello")
	if err != nil {
		t.Fatalf("fetch fin server fail with %s", err)
	}
	defer resp.Body.Close()
	payload := make([]byte, 5)
	if _, err := resp.Body.Read(payload); err != io.EOF && err != nil {
		t.Fatalf("read body fail with %s", err)
	}
	if string(payload) != "hello" {
		t.Fatalf("server response body is not excepted %s", string(payload))
	}
}
