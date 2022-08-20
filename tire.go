package triecache

import (
	"strings"
	"sync"
	"time"
)

type treeNode struct {
	Key      string      `json:"key" yaml:"key" bson:"key"`
	Value    interface{} `json:"value" bson:"value" yaml:"value"`
	Index    int         `json:"index" bson:"index" yaml:"index"`
	ExTime   int64       `json:"ex_time" bson:"ex_time" yaml:"ex_time"`
	Parent   *treeNode   `json:"parent" bson:"parent" yaml:"parent"`
	Children []*treeNode `json:"children" bson:"children" yaml:"children"`
	mu       *sync.Mutex
}

func newRootNode() *treeNode {
	root := new(treeNode)
	root.Value = ""
	root.Index = -1
	root.Parent = nil
	root.ExTime = -1
	root.Children = []*treeNode{}
	root.mu = &sync.Mutex{}

	return root
}

func newNode(key string, value interface{}, index int) *treeNode {
	node := new(treeNode)
	node.Key = key
	node.Value = value
	node.Index = index
	node.Children = []*treeNode{}
	node.mu = &sync.Mutex{}
	return node
}

func (n *treeNode) search(index int, key string) (node *treeNode) {
	if n.Key == key && index == n.Index {
		return n
	}
	for _, v := range n.Children {
		if v.search(v.Index, key) != nil {
			return v
		}
	}
	return node
}

func (n *treeNode) fuzzySearch(key string, node *[]treeNode) {
	if strings.Contains(n.Key, key) {
		if node != nil {
			*node = append(*node, *n)
		}
	}
	for _, v := range n.Children {
		v.fuzzySearch(key, node)
	}
}

func (n *treeNode) addChild(node *treeNode) {
	n.mu.Lock()
	defer n.mu.Unlock()
	node.Index = n.Index + 1
	if n.Children == nil {
		n.Children = []*treeNode{}
	}
	node.Parent = n
	have := false
	for _, child := range n.Children {
		if child.Key == node.Key {
			have = true
		}
	}
	if have == false {
		n.Children = append(n.Children, node)
	}
}

func (n *treeNode) delete(index int, key string) {
	if n.Key == key {
		n = new(treeNode)
	} else {
		parent := n.search(index, n.search(index, key).Key)
		var tmp []*treeNode
		for _, v := range parent.Children {
			if v.Value == key {
				continue
			}
			tmp = append(tmp, v)
		}
		parent.Children = tmp
	}
}

func (n *treeNode) nodeAdd(index int, key string, child *treeNode) {
	n.search(index, key).addChild(child)
}

func (n *treeNode) getChildKeys(key string, node *[]string) {
	//if len(key)-1 != n.Index {
	//	return
	//}
	if n.ExTime < time.Now().Unix() {
		n.del()
	}
	for _, v := range n.Children {
		if v == nil {
			continue
		}
		if v.ExTime < time.Now().Unix() {
			v.del()
		}
		if v.Value != nil {
			*node = append(*node, key+v.Key)
			return
		} else {
			v.getChildKeys(key+v.Key, node)
		}
	}
}

func (n *treeNode) checkExpire() {
	for _, v := range n.Children {
		if v.ExTime < time.Now().Unix() && v.Value != nil {
			v.del()
		}
		if len(v.Children) > 0 {
			v.checkExpire()
		}
	}
}

func (n *treeNode) set(value interface{}, ex time.Duration) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.Value = value
	n.ExTime = time.Now().Unix() + int64(ex.Seconds())
}

func (n *treeNode) del() {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.Value = nil
	n.ExTime = 0
}
