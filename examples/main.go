package main

import (
	"fmt"
	"github.com/lzy1102/triecache"
	"time"
)

func main() {
	key := "hellofdafas123"
	key2 := "helloerqwr456"
	value := 10
	var c triecache.Cache
	c = triecache.New(time.Minute*5, time.Second*10)
	err := c.Set(key, value, time.Second*10)
	if err != nil {
		panic(err)
	}
	c.Set(key2, value, time.Second)
	time.Sleep(time.Second * 2)
	get, err := c.Get(key)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("key ", key, "  value ", get)

	get, err = c.Get(key2)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("key2 ", key2, "  value2 ", get)
	keys, err := c.Keys("hello*")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("keys ", keys)
	for _, s := range keys {
		v, _ := c.Get(s)
		fmt.Println(s, v)
	}
}
