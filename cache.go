package main

import (
	"fmt"
	"strings"
	"time"
)

type treeNode struct {
	Key      string      `json:"key" yaml:"key" bson:"key"`
	Value    interface{} `json:"value" bson:"value" yaml:"value"`
	Index    int         `json:"index" bson:"index" yaml:"index"`
	ExTime   int64       `json:"ex_time" bson:"ex_time" yaml:"ex_time"`
	Parent   *treeNode   `json:"parent" bson:"parent" yaml:"parent"`
	Children []*treeNode `json:"children" bson:"children" yaml:"children"`
}

func newRootNode() *treeNode {
	root := new(treeNode)
	root.Value = ""
	root.Index = -1
	root.Parent = nil
	root.Children = []*treeNode{}
	return root
}

func newNode(key string, value interface{}, index int, ex int64) *treeNode {
	node := new(treeNode)
	node.Key = key
	node.Value = value
	node.Index = index
	node.ExTime = time.Now().Unix() + ex
	node.Children = []*treeNode{}
	return node
}

func (n *treeNode) search(index int, key string) (node *treeNode) {
	if n.Key == key && index == n.Index {
		return n
	}
	for _, v := range n.Children {
		if v.Key == key {
			return v
		} else {
			node = v.search(index, key)
			if node != nil && node.Value == key {
				return node
			}
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
	node.Index = n.Index + 1
	if n.Children == nil {
		n.Children = []*treeNode{}
	}
	node.Parent = n
	n.Children = append(n.Children, node)
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

func (n *treeNode) getChildKeys(key string) (node []string) {
	for _, v := range n.Children {
		if v.Value != nil {
			node = append(node, key+v.Key)
		}
		node = append(node, v.getChildKeys(key+v.Key)...)
	}
	return node
}

func (n *treeNode) getParent() {

}

type Cache struct {
	root *treeNode
}

func NewCache() *Cache {
	c := new(Cache)
	c.root = newRootNode()
	return c
}

func (c *Cache) Set(key string, value interface{}, ex time.Duration) error {
	var node *treeNode
	node = c.root
	for i := 0; i < len(key); i++ {
		tmp := node.search(i, string(key[i]))
		if tmp == nil {
			tmp = newNode(string(key[i]), nil, i, int64(ex.Seconds()))
			node.addChild(tmp)
		}
		if i == len(key)-1 {
			tmp.Value = value
		}
		node = tmp
	}
	return nil
}

func (c *Cache) Get(key string) (interface{}, error) {
	var node *treeNode
	node = c.root
	for i := 0; i < len(key); i++ {
		tmp := node.search(i, string(key[i]))
		if tmp == nil {
			return nil, fmt.Errorf("not key")
		}
		node = tmp
	}
	if node.Value == nil {
		return nil, fmt.Errorf("not value")
	}
	return node.Value, nil
}

func (c *Cache) Scan(pattern string) ([]string, error) {
	var node *treeNode
	node = c.root
	if pattern[len(pattern)-1] != '*' {
		return nil, fmt.Errorf("not *")
	}
	for i := 0; i < len(pattern)-1; i++ {
		tmp := node.search(i, string(pattern[i]))
		if tmp == nil {
			return nil, fmt.Errorf("not key")
		}
		node = tmp
	}
	nodes := node.getChildKeys(pattern[0 : len(pattern)-1])
	return nodes, nil
}

func (c *Cache) ExPire(key string, ex time.Duration) error {
	var node *treeNode
	node = c.root
	for i := 0; i < len(key); i++ {
		tmp := node.search(i, string(key[i]))
		if tmp == nil {
			return fmt.Errorf("not key")
		}
		node = tmp
	}
	if node.Value == nil {
		return fmt.Errorf("not value")
	}
	node.ExTime = time.Now().Unix() + int64(ex.Seconds())
	return nil
}
