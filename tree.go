package fin

import "fmt"

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

// longestSharePrefix查找两个字符串的最长公共前缀
func longestSharePrefix(a, b string) int {
	i := 0
	max := min(len(a), len(b))
	for i < max && a[i] == b[i] {
		i++
	}

	return i
}

type node struct {
	chain []HandlerFunc

	path     string
	children []*node
	indices  string // 每个子节点的path首字母
	isEnd    bool
}

func (n *node) search(s string) *node {
	if s == n.path {
		return n
	}

	// 先匹配自身的path，for compressed
	i := 0
	for i < min(len(n.path), len(s)) {
		if n.path[i] != s[i] {
			break
		}
		i++
	}

	// 接着s的剩余部分匹配indices子节点
	s = s[i:]
	if len(s) == 0 {
		return nil
	}

	head := s[0]
	for i := 0; i < len(n.indices); i++ {
		if n.indices[i] == head {
			return n.children[i].search(s)
		}
	}

	return nil
}

func (n *node) insert(s string, child *node) {
	// 获取当前节点path与待插入的s的最长公共前缀长度
	sharePrefix := longestSharePrefix(n.path, s)
	// 当前节点的path长度大于共有的前缀长度，则需要分裂当前节点的path为两个
	if sharePrefix < len(n.path) {
		// 当前节点的存储的chain等信息转移到分裂后的子节点
		splitNode := &node{
			chain:    n.chain,
			path:     n.path[sharePrefix:],
			indices:  n.indices,
			children: n.children,
			isEnd:    n.isEnd,
		}
		// 分裂当前节点
		n.path = n.path[:sharePrefix]
		n.children = []*node{splitNode}
		n.indices = string([]byte{splitNode.path[0]})
		n.isEnd = false
		n.chain = nil
	}
	// 待插入节点的s长度大于共有的前缀长度，则需要分割带插入字符串s
	if sharePrefix < len(s) {
		s = s[sharePrefix:]
		child.path = s
		head := s[0]
		// 若当前节点存在以s[0]的indices索引，则转移到对应的indices子节点，继续插入s
		for i := 0; i < len(n.indices); i++ {
			if n.indices[i] == head {
				n.children[i].insert(s, child)
				return
			}
		}
		// 若不存在对应的indices索引，则直接插入到当前节点的子节点
		n.children = append(n.children, child)
		n.indices += string([]byte{head})
		return
	}
	// 待插入的节点刚好在当前节点上，直接补充当前节点的信息
	if sharePrefix == len(s) {
		if n.isEnd {
			fmt.Printf("current:%+v, child:%+v", n, child)
			panic("duplicate uri for " + s + ", current path is " + n.path)
		}
		n.isEnd = true
		n.chain = child.chain
	}
}

func (n *node) addRoute(path string, h ...HandlerFunc) {
	n.insert(path, &node{
		chain: h,
		path:  path,
		isEnd: true,
	})
}

func (n *node) getValue(path string) (handlers []HandlerFunc) {
	matched := n.search(path)
	if matched == nil || !matched.isEnd {
		return
	}

	handlers = matched.chain
	return
}

type tree struct {
	method string
	root   *node
}

type trees []*tree

func (t trees) get(method string) *node {
	for _, tree := range t {
		if tree.method == method {
			return tree.root
		}
	}

	return nil
}
