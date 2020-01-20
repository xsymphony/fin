package fin

import (
	"fmt"
)

type node struct {
	chain []HandlerFunc

	path     string
	children map[string]*node
	isEnd    bool
}

func (n *node) getChild(s string) *node {
	if n.children == nil {
		return nil
	}

	if child, exists := n.children[s]; exists {
		return child
	}

	return nil
}

func (n *node) hasChild(s string) bool {
	return n.getChild(s) != nil
}

func (n *node) addChild(s string, child *node) {
	if n.children == nil {
		n.children = make(map[string]*node)
	}

	if _, exists := n.children[s]; exists {
		return
	}

	n.children[s] = child
}

type tree struct {
	method string
	root   *node
}

func newTree(method string) *tree {
	return &tree{
		method: method,
		root:   new(node),
	}
}

func (t tree) addRoute(path string, h ...HandlerFunc) {
	root := t.root
	for _, s := range path {
		if !root.hasChild(string(s)) {
			root.addChild(string(s), &node{
				path: string(s),
			})
		}
		root = root.getChild(string(s))
	}

	// 节点的isEnd=true，说明已经注册了重复的路由
	if root.isEnd {
		panic(fmt.Sprintf("duplicate uri: %s", path))
	}

	root.isEnd = true
	root.chain = h
}

func (t tree) search(path string) *node {
	root := t.root
	for _, s := range path {
		ss := string(s)
		next := root.getChild(ss)
		if next == nil {
			return nil
		}
		root = next
	}

	// 没有匹配到完整的路径
	if !root.isEnd {
		return nil
	}

	return root
}

type trees []*tree

func (t trees) get(method string) *tree {
	for _, tree := range t {
		if tree.method == method {
			return tree
		}
	}

	return nil
}
