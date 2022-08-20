package triecache

import (
	"encoding/json"
	"strings"
	"sync"
	"time"
)

type TreeNode struct {
	Key      string      `json:"key"  bson:"key" `
	Value    interface{} `json:"value" bson:"value" `
	Index    int         `json:"index" bson:"index" `
	ExTime   int64       `json:"ex_time" bson:"ex_time" `
	Parent   *TreeNode   `json:"parent" bson:"parent" `
	Children []*TreeNode `json:"children" bson:"children" `
	Mu       *sync.Mutex
}

func newRootNode() *TreeNode {
	root := new(TreeNode)
	root.Value = ""
	root.Index = -1
	root.Parent = nil
	root.ExTime = -1
	root.Children = []*TreeNode{}
	root.Mu = &sync.Mutex{}

	return root
}

func newNode(key string, value interface{}, index int) *TreeNode {
	node := new(TreeNode)
	node.Key = key
	node.Value = value
	node.Index = index
	node.Children = []*TreeNode{}
	node.Mu = &sync.Mutex{}
	return node
}

func (n *TreeNode) marshal() []byte {
	marshal, err := json.Marshal(n)
	if err != nil {
		panic(err)
	}
	return marshal
}

func (n *TreeNode) search(index int, key string) (node *TreeNode) {
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

func (n *TreeNode) fuzzySearch(key string, node *[]TreeNode) {
	if strings.Contains(n.Key, key) {
		if node != nil {
			*node = append(*node, *n)
		}
	}
	for _, v := range n.Children {
		v.fuzzySearch(key, node)
	}
}

func (n *TreeNode) addChild(node *TreeNode) {
	n.Mu.Lock()
	defer n.Mu.Unlock()
	node.Index = n.Index + 1
	if n.Children == nil {
		n.Children = []*TreeNode{}
	}
	node.Parent = n
	n.Children = append(n.Children, node)
}

func (n *TreeNode) nodeAdd(index int, key string, child *TreeNode) {
	n.search(index, key).addChild(child)
}

func (n *TreeNode) getChildKeys(key string, node *[]string) {
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
			//return
		} else {
			v.getChildKeys(key+v.Key, node)
		}
	}
}

func (n *TreeNode) checkExpire() {
	for _, v := range n.Children {
		if v.ExTime < time.Now().Unix() && v.Value != nil {
			v.del()
		}
		if len(v.Children) > 0 {
			v.checkExpire()
		}
	}
}

func (n *TreeNode) set(value interface{}, ex time.Duration) {
	n.Mu.Lock()
	defer n.Mu.Unlock()
	n.Value = value
	n.ExTime = time.Now().Unix() + int64(ex.Seconds())
}

func (n *TreeNode) del() {
	n.Mu.Lock()
	defer n.Mu.Unlock()
	n.Value = nil
	n.ExTime = 0
	if len(n.Children) <= 0 {
		n = nil
	}
}
