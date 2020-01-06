package fin

import "fmt"

type node map[string]HandlerFunc

func (n node) addRoute(path string, h HandlerFunc) {
	if _, ok := n[path]; ok {
		panic(fmt.Sprintf("duplicate uri %s", path))
	}
	n[path] = h
}

type tree struct {
	method string
	node   node
}

type trees []tree

func (t trees) get(method string) node {
	for _, tree := range t {
		if tree.method == method {
			return tree.node
		}
	}

	return nil
}
