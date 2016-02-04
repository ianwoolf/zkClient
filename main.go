package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	c "github.com/ianwoolf/zkClient/command"
	u "github.com/ianwoolf/zkClient/util"
)

type FlagParam []string

func (f *FlagParam) String() string {
	return "string method"
}

func (f *FlagParam) Set(value string) error {
	*f = strings.Split(value, ",")
	return nil
}

var (
	zh      *u.ZH
	servers FlagParam
	command string
	path    string
	// todo interface
	data     string
	version  int
	flags    int
	usageMsg string = `Usage of %s:zkClient -c [param] -p [param] (-d -v -flags) 
command list: get/set/child/creat/watchExist/watchChildren/delChildNode

 `
)

func parasFlag() []string {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, usageMsg, os.Args[0])
		flag.PrintDefaults()
	}

	flag.Var(&servers, "s", "zk server list. e.g: 127.0.0.1:2181,127.0.0.2:2181.")
	flag.StringVar(&command, "c", "", "command, such as: get/set/child/creat/watch/del")
	flag.StringVar(&path, "p", "/mynode/test", "node path")
	flag.StringVar(&data, "d", "test set2", "string data")
	flag.IntVar(&version, "v", 0, "data version")
	flag.IntVar(&flags, "flags", 0, "flag: 0-Permanent 1 2-sequence")
	flag.Parse()
	return flag.Args()
}

func initZk(servers []string, timeout int) (err error) {
	zh, err = u.NewZH(servers, timeout)
	return
}

func main() {
	var timeout int = 3
	parasFlag()

	err := initZk(servers, timeout)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
	defer zh.Close()

	switch command {
	case "get":
		content, stat, err := c.Get(zh, path)
		if err != nil {
			fmt.Println("get zk node fail:", path, err.Error())
		} else {
			fmt.Println("content:", content)
			fmt.Printf("version: %d, time: %v", stat.Version(), stat.MTime())
		}
	case "child":
		paths, err := c.Children(zh, path)
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println(paths)
	case "set":
		nodeToSet := path
		_, errSet := c.Set(zh, nodeToSet, []byte(data), int32(version))
		if errSet != nil {
			fmt.Println(errSet.Error())
		}
		content, _, err := c.Get(zh, nodeToSet)
		if err != nil {
			fmt.Println("get node %s fail when set: %s", nodeToSet, err.Error())
			return
		}
		fmt.Println("after set:", content)
	case "lock":
		nodeToSet := path + "/test1"
		fmt.Println("begin to lock")
		lock := c.GetLock(zh, nodeToSet)
		lock.Lock()
		content, stat, err := c.Get(zh, nodeToSet)
		if err != nil {
			fmt.Println("get zk node fail when set:", err.Error())
			return
		}
		fmt.Println("before set begin lock and set:", content)
		Verl := stat.Version()

		_, errSl := zh.Set(nodeToSet, []byte("test set in lock"), int32(Verl))
		if errSl != nil {
			fmt.Println(errSl.Error())
		}
		content, _, err = c.Get(zh, nodeToSet)
		if err != nil {
			fmt.Println("get zk node fail in lock:", err.Error())
			return
		}
		fmt.Println("after set in lock:", content)
		lock.Unlock()

	case "create":
		//flags: 0- Permanent node  2- sequence node
		fmt.Println(path, data, flags)
		c.Create(zh, path, []byte(data), int32(flags))
		paths, cerr := zh.Children(path)
		if cerr != nil {
			fmt.Println(cerr.Error())
		}
		fmt.Println("after create:", paths)

	case "watch":
		// todo: watch by param type: exist/child
		fmt.Println("check watch on path:")
		existOk, Eevent := c.WatchExist(zh, path)
		if existOk {
			fmt.Println(path, "exist and watch")
		}
		// _, Gevent, GWerr := zh.GetW(path)
		// if GWerr == nil {
		// 	fmt.Println(path, "get zk node success and watch")
		// }
		_, Cevent, CWerr := c.WatchChildren(zh, path)
		if CWerr != nil {
			fmt.Println(CWerr.Error())
		}
		select {
		case a := <-Eevent:
			fmt.Println("Exist event")
			fmt.Println(a.Type.String(), a.Path)
		case c := <-Cevent:
			fmt.Println("Children event")
			fmt.Println(c)
		}
	case "del":
		// delete child nodes
		c.DelChildNode(zh, path)
	}
}
