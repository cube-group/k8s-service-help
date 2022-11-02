package main

import (
	"fmt"
	"github.com/go-zookeeper/zk"
	"strings"
	"time"
)

var Ip string
var Path = "/a"

func main() {
	fmt.Println("ip", Ip)
	c, _, err := zk.Connect(strings.Split(Ip, ","), time.Second, zk.WithEventCallback(callback)) //*10)
	if err != nil {
		panic(err)
	}
	fmt.Println("server", c.Server(), c.State().String())
	res, stat, err := c.Get(Path)
	fmt.Println("get", res, stat.DataLength, err)
	if err != nil {
		res, err := c.Create(Path, []byte("hello"), 0, zk.WorldACL(zk.PermAll))
		fmt.Println("create", res, err)
	} else {
		stat, err := c.Set(Path, []byte("hello"), 0)
		fmt.Println("set", stat.DataLength, err)
	}

	children, stat, ch, err := c.ChildrenW("/")
	fmt.Println("childrenW", children, stat.DataLength, ch, err)

	fmt.Printf("%+v %+v\n", children, stat)
	e := <-ch
	fmt.Printf("%+v\n", e)
}

// zk watch 回调函数
func callback(event zk.Event) {
	// zk.EventNodeCreated
	// zk.EventNodeDeleted
	fmt.Println("###########################")
	fmt.Println("path: ", event.Path)
	fmt.Println("type: ", event.Type.String())
	fmt.Println("state: ", event.State.String())
	fmt.Println("---------------------------")
}
