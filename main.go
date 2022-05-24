package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"
)

func main() {
	key := "lizhiyong123"
	key2 := "lizhiyong456"
	value := 10
	c := NewCache()
	err := c.Set(key, value, time.Second*10)
	if err != nil {
		return
	}
	c.Set(key2, value, time.Second*10)
	get, err := c.Get(key)
	if err != nil {
		return
	}
	fmt.Println("key ", key, "  value ", get)

	get, err = c.Get(key2)
	if err != nil {
		return
	}
	fmt.Println("key2 ", key2, "  value2 ", get)
	scan, err := c.Scan("lizhiyong*")
	if err != nil {
		return
	}
	fmt.Println(scan)
	for _, s := range scan {
		v, _ := c.Get(s)
		fmt.Println(s, v)
	}
	marshal, err := json.Marshal(c.root)
	if err != nil {
		return
	}
	ioutil.WriteFile("root.json", marshal, 777)
}
