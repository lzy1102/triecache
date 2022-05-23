package treecache

import (
	"strings"
)

type treeNode struct {
	Value    string      `json:"value" bson:"value" yaml:"value"`
	Index    int         `json:"index" bson:"index" yaml:"index"`
	Children []*treeNode `json:"children" bson:"children" yaml:"children"`
}

func newRootNode() *treeNode {
	root := new(treeNode)
	root.Value = ""
	root.Index = -1
	root.Children = []*treeNode{}
	return root
}

func (n *treeNode) search(index int, value string) (node *treeNode) {
	if n.Value == value && index == n.Index {
		return n
	}
	for _, v := range n.Children {
		if v.Value == value {
			return v
		} else {
			node = v.search(index, value)
			if node != nil && node.Value == value {
				return node
			}
		}
	}
	return node
}

func (n *treeNode) fuzzySearch(label string, node *[]treeNode) {
	if strings.Contains(n.Label, label) {
		if node != nil {
			*node = append(*node, *n)
		}
	}
	for _, v := range n.Children {
		v.fuzzySearch(label, node)
	}
}

func (n *treeNode) addChild(node *treeNode) {
	node.Index = n.Index + 1
	if n.Children == nil {
		n.Children = []*treeNode{}
	}
	n.Children = append(n.Children, node)
}

func (n *treeNode) delete(value string) {
	if n.Value == value {
		n = new(treeNode)
	} else {
		parent := n.search(n.search(value).Value)
		var tmp []*treeNode
		for _, v := range parent.Children {
			if v.Value == value {
				continue
			}
			tmp = append(tmp, v)
		}
		parent.Children = tmp
	}
}

func (n *treeNode) nodeAdd(uuid string, child *treeNode) {
	n.search(uuid).addChild(child)
}

func (n *treeNode) getChild(uuid string) {
	for _, v := range n.search(uuid).Children {
		if len(v.Children) > 0 && v.Children != nil {
			v.Children = nil
		}
	}
}
