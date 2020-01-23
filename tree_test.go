package fin

import (
	"fmt"
	"testing"
)

var fakeHandlerValue string

func fakeHandler(route string) func(ctx *Context) {
	return func(ctx *Context) {
		fakeHandlerValue = route
	}
}

type testRequests []struct {
	path       string
	nilHandler bool
	route      string
}

func checkRequests(t *testing.T, root *node, requests testRequests) {
	for _, request := range requests {
		handlers := root.getValue(request.path)
		fmt.Printf("path:%s, handlers:%+v\n", request.path, handlers)
		if handlers == nil {
			if !request.nilHandler {
				t.Errorf("handle mismatch for route '%s': Expected non-nil handle", request.path)
			}
		} else if request.nilHandler {
			t.Errorf("handle mismatch for route '%s': Expected nil handle", request.path)
		} else {
			handlers[0](nil)
			if fakeHandlerValue != request.route {
				t.Errorf("handle mismatch for route '%s': Wrong handle (%s != %s)", request.path, fakeHandlerValue, request.route)
			}
		}
	}
}

func TestTree(t *testing.T) {
	root := new(node)

	routes := [...]string{
		"/hi",
		"/contact",
		"/co",
		"/c",
		"/a",
		"/ab",
		"/doc/",
		"/doc/go_faq.html",
		"/doc/go1.html",
		"/α",
		"/β",
	}
	for _, route := range routes {
		root.addRoute(route, fakeHandler(route))
	}

	checkRequests(t, root, testRequests{
		{"/a", false, "/a"},
		{"/", true, ""},
		{"/hi", false, "/hi"},
		{"/contact", false, "/contact"},
		{"/co", false, "/co"},
		{"/con", true, ""},  // key mismatch
		{"/cona", true, ""}, // key mismatch
		{"/no", true, ""},   // no matching child
		{"/ab", false, "/ab"},
		{"/α", false, "/α"},
		{"/β", false, "/β"},
	})
}
