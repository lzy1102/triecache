package main

import (
	"fmt"
	"github.com/lzy1102/triecache"
	"time"
)

func main() {
	key := "TcpStream&&Down&&192.168.1.104&&139.155.249.48&&8181(intermapper)&&55660&&1206525595&&2961302792&&1412&&1660286309000021196"
	key2 := "TcpStream&&Down&&192.168.1.104&&139.155.249.48&&8181(intermapper)&&55660&&1206525595&&2961304204&&1412&&1660286309000021319"
	value := 10
	var c triecache.Cache
	c = triecache.New(time.Minute, time.Second*5)
	err := c.Set(key, value, time.Second*10)
	if err != nil {
		panic(err)
	}
	c.Set(key2, value, time.Second)

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
	keys, err := c.Keys("TcpStream&&Down&&192.168.1.104&&139.155.249.48&&8181(intermapper)&&55660&&1206525595&&*")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("keys ", keys)
	for _, s := range keys {
		v, _ := c.Get(s)
		fmt.Println(s, v)
	}
	time.Sleep(time.Second * 5)
}
