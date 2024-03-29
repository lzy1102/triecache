package triecache

import (
	"fmt"
	"strconv"
	"time"
)

type cacheInfo struct {
	Root    *treeNode     `json:"root"  bson:"root"`
	Extime  time.Duration `json:"extime"  bson:"extime" jsonapi:"extime"`
	Polling time.Duration `json:"polling" bson:"polling"`
}

func New(extime time.Duration, polling time.Duration) Cache {
	c := new(cacheInfo)
	c.Root = newRootNode()
	c.Extime = extime
	c.Polling = polling
	go c.ticker()
	return c
}

func (c *cacheInfo) ticker() {
	tic := time.NewTicker(c.Polling)
	for range tic.C {
		c.Root.checkExpire()
	}
}

func (c *cacheInfo) find(key string) *treeNode {
	var node *treeNode
	node = c.Root
	for i := 0; i < len(key); i++ {
		tmp := node.search(i, string(key[i]))
		if tmp == nil {
			tmp = newNode(string(key[i]), nil, i)
			node.addChild(tmp)
		}
		node = tmp
		if i == len(key)-1 {
			break
		}
	}

	return node
}

func (c *cacheInfo) Set(key string, value interface{}, ex time.Duration) error {
	if &ex == nil {
		ex = c.Extime
	}
	node := c.find(key)
	if node == nil {
		return fmt.Errorf("set node nil")
	}
	node.set(value, ex)
	return nil
}

func (c *cacheInfo) Get(key string) (interface{}, error) {
	node := c.find(key)
	if node == nil {
		return nil, fmt.Errorf("not key")
	}
	if node.Value == nil {
		return nil, fmt.Errorf("not key value")
	}
	if node.ExTime < time.Now().Unix() {
		node.del()
		return nil, fmt.Errorf("key extime")
	}
	return node.Value, nil
}

func (c *cacheInfo) Delete(key string) error {
	node := c.find(key)
	if node != nil {
		node.del()
	}
	return nil
}

func (c *cacheInfo) Keys(pattern string) ([]string, error) {
	if pattern[len(pattern)-1] != '*' {
		return nil, fmt.Errorf("not *")
	}

	node := c.find(pattern[:len(pattern)-1])
	if node == nil {
		return nil, fmt.Errorf("not key")
	}
	var nodes []string
	node.getChildKeys(pattern[:len(pattern)-1], &nodes)
	return nodes, nil
}

func (c *cacheInfo) Expire(key string, ex time.Duration) error {
	if &ex == nil {
		ex = c.Extime
	}
	node := c.find(key)
	if node == nil {
		return fmt.Errorf("not key")
	}
	if node.Value == nil {
		return fmt.Errorf("not value")
	}
	node.ExTime = time.Now().Unix() + int64(ex.Seconds())
	return nil
}

func (c *cacheInfo) TTL(key string) (int64, error) {
	node := c.find(key)
	if node == nil {
		return 0, fmt.Errorf("not key")
	}
	if node.Value == nil {
		return 0, fmt.Errorf("not value")
	}
	return node.ExTime - time.Now().Unix(), nil
}

func (c *cacheInfo) GetInt64(key string) (int64, error) {
	node := c.find(key)
	if node == nil {
		return 0, fmt.Errorf("not key")
	}
	if node.Value == nil {
		return 0, fmt.Errorf("not key value")
	}
	if node.ExTime < time.Now().Unix() {
		node.del()
		return 0, fmt.Errorf("key extime")
	}
	result, err := strconv.ParseInt(fmt.Sprintf("%v", node.Value), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("value not to int64")
	}
	return result, nil
}

func (c *cacheInfo) GetFloat64(key string) (float64, error) {
	node := c.find(key)
	if node == nil {
		return 0, fmt.Errorf("not key")
	}
	if node.Value == nil {
		return 0, fmt.Errorf("not key value")
	}
	if node.ExTime < time.Now().Unix() {
		node.del()
		return 0, fmt.Errorf("key extime")
	}
	result, err := strconv.ParseFloat(fmt.Sprintf("%v", node.Value), 64)
	if err != nil {
		return 0, fmt.Errorf("value not to int64")
	}
	return result, nil
}

func (c *cacheInfo) Incr(key string, ex time.Duration) (int64, error) {
	if &ex == nil {
		ex = c.Extime
	}
	node := c.find(key)
	if node == nil {
		return 0, fmt.Errorf("set node nil")
	}
	if node.Value == nil {
		node.set(int64(1), ex)
		return 1, nil
	}
	result, err := strconv.ParseInt(fmt.Sprintf("%v", node.Value), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("value not to int64")
	}
	node.set(result+1, ex)
	return result + 1, nil
}

func (c *cacheInfo) IncrBy(key string, value int64, ex time.Duration) (int64, error) {
	if &ex == nil {
		ex = c.Extime
	}
	node := c.find(key)
	if node == nil {
		return 0, fmt.Errorf("set node nil")
	}
	if node.Value == nil {
		node.set(value, ex)
		return value, nil
	}
	result, err := strconv.ParseInt(fmt.Sprintf("%v", node.Value), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("value not to int64")
	}
	node.set(result+value, ex)
	return result + value, nil
}

func (c *cacheInfo) Decr(key string, ex time.Duration) (int64, error) {
	if &ex == nil {
		ex = c.Extime
	}
	node := c.find(key)
	if node == nil {
		return 0, fmt.Errorf("set node nil")
	}
	if node.Value == nil {
		node.set(int64(-1), ex)
		return -1, nil
	}
	result, err := strconv.ParseInt(fmt.Sprintf("%v", node.Value), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("value not to int64")
	}
	node.set(result-1, ex)
	return result - 1, nil
}

func (c *cacheInfo) DecrBy(key string, value int64, ex time.Duration) (int64, error) {
	if &ex == nil {
		ex = c.Extime
	}
	node := c.find(key)
	if node == nil {
		return 0, fmt.Errorf("set node nil")
	}
	if node.Value == nil {
		node.set(-value, ex)
		return -value, nil
	}
	result, err := strconv.ParseInt(fmt.Sprintf("%v", node.Value), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("value not to int64")
	}
	node.set(result-value, ex)
	return result - value, nil
}
