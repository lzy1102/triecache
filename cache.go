package triecache

import (
	"fmt"
	"strconv"
	"time"
)

type cache struct {
	root    *treeNode
	extime  time.Duration
	polling time.Duration
}

func New(extime time.Duration, polling time.Duration) *cache {
	c := new(cache)
	c.root = newRootNode()
	c.extime = extime
	c.polling = polling
	go c.ticker()
	return c
}

func (c *cache) ticker() {
	tic := time.NewTicker(c.polling)
	for range tic.C {
		c.root.checkExpire()
	}
}

func (c *cache) find(key string) *treeNode {
	var node *treeNode
	node = c.root

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

func (c *cache) Set(key string, value interface{}, ex time.Duration) error {
	if &ex == nil {
		ex = c.extime
	}
	node := c.find(key)
	if node == nil {
		return fmt.Errorf("set node nil")
	}
	node.set(value, ex)
	return nil
}

func (c *cache) Get(key string) (interface{}, error) {
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

func (c *cache) Keys(pattern string) ([]string, error) {
	if pattern[len(pattern)-1] != '*' {
		return nil, fmt.Errorf("not *")
	}
	node := c.find(pattern[:len(pattern)-1])
	if node == nil {
		return nil, fmt.Errorf("not key")
	}
	nodes := node.getChildKeys(pattern[:len(pattern)-1])
	return nodes, nil
}

func (c *cache) Expire(key string, ex time.Duration) error {
	if &ex == nil {
		ex = c.extime
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

func (c *cache) TTL(key string) (int64, error) {
	node := c.find(key)
	if node == nil {
		return 0, fmt.Errorf("not key")
	}
	if node.Value == nil {
		return 0, fmt.Errorf("not value")
	}
	return node.ExTime - time.Now().Unix(), nil
}

func (c *cache) GetInt64(key string) (int64, error) {
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

func (c *cache) GetFloat64(key string) (float64, error) {
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

func (c *cache) Incr(key string, ex time.Duration) (int64, error) {
	if &ex == nil {
		ex = c.extime
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

func (c *cache) IncrBy(key string, value int64, ex time.Duration) (int64, error) {
	if &ex == nil {
		ex = c.extime
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

func (c *cache) Decr(key string, ex time.Duration) (int64, error) {
	if &ex == nil {
		ex = c.extime
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

func (c *cache) DecrBy(key string, value int64, ex time.Duration) (int64, error) {
	if &ex == nil {
		ex = c.extime
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
