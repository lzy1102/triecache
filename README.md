# triecache

triecache是一个内存中的 key  -  value 缓存，用法类似redis，对于一些想要临时存储一些数据，但是没必要接入redis的场景适用；原理是通过树形结构实现的，相较于根据map[string]interface{} 的  go-cache 等，实现了 redis 中 keys的功能



#### **安装** 



`go get -u github.com/lzy1102/triecache`



### 用法



```
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

```

