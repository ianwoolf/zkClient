package main

import (
	"flag"
	"fmt"
	"os"

	c "github.com/ianwoolf/zkClient/command"
	u "github.com/ianwoolf/zkClient/util"
	// zk "go.intra.xiaojukeji.com/golang/go-zookeeper/zk"
)

// func get(zh *u.ZH, path string) (content string, stat zk.Stat, err error) {
// 	var read []byte
// 	read, stat, err = zh.Get(path)
// 	if err == nil {
// 		content = string(read)
// 	}
// 	return
// }

// func children(zh *u.ZH, path string) (paths []string, err error) {
// 	paths, err = zh.Children(path)
// 	return
// }

// func set(zh *u.ZH, path string, data []byte, version int32) (stat zk.Stat, err error) {
// 	stat, err = zh.Set(path, data, version)
// 	return
// }

// func getLock(zh *u.ZH, path string) *zk.Lock {
// 	return zh.NewLock(path)
// }

// func create(zh *u.ZH, path string, data []byte, flags int32) (string, error) {
// 	return zh.Create(path, data, flags)
// }

// func watchExist(zh *u.ZH, path string) (ok bool, event <-chan zk.Event) {
// 	ok, event = zh.ExistsW(path)
// 	return
// }

// func watchChildren(zh *u.ZH, path string) (children []string, event <-chan zk.Event, err error) {
// 	children, event, err = zh.ChildrenW(path)
// 	return
// }

// func delChildNode(zh *u.ZH, path string) (err error) {
// 	var childNodes []string
// 	var stat zk.Stat
// 	childNodes, err = zh.Children(path)
// 	if err != nil {
// 		fmt.Println(err.Error())
// 		return
// 	}
// 	fmt.Println("begin to del node:", childNodes)
// 	for _, node := range childNodes {
// 		nodeToDel := fmt.Sprintf("%s/%s", path, node)
// 		_, stat, err = zh.Get(nodeToDel)
// 		if err != nil {
// 			fmt.Println("get zk node fail:", nodeToDel, err.Error())
// 			return
// 		} else {
// 			fmt.Println("begin to del node:", nodeToDel)
// 		}
// 		err = zh.Delete(nodeToDel, int32(stat.Version()))
// 		if err != nil {
// 			fmt.Printf("del node %s error: %s", nodeToDel, err.Error())
// 			continue
// 		}
// 	}
// 	return
// }

var (
	command string
	path    string
	// todo interface
	data     string
	version  int
	flags    int
	usageMsg string = `Usage of %s:zkClient -c [param] -p [param] (-d -v -flags) 
command list: get/set/child/creat/watchExist/watchChildren/delChildNode

 `
	servers []string = []string{"106.186.127.250:2181"}
	timeout int      = 3
	zh      *u.ZH
)

func parasFlag() []string {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, usageMsg, os.Args[0])
		flag.PrintDefaults()
	}

	flag.StringVar(&command, "c", "", "command, such as: get/set/child/creat/watchExist/watchChildren/del")
	flag.StringVar(&path, "p", "/mynode/test", "node path")
	flag.StringVar(&data, "d", "test set2", "string data")
	flag.IntVar(&version, "v", 0, "data version")
	flag.IntVar(&flags, "flags", 0, "flag: 0-Permanent 1 2-sequence")
	flag.Parse()
	return flag.Args()
}

func init() {
	var err error
	zh, err = u.NewZH(servers, timeout)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
	defer zh.Close()
}

func main() {
	args := parasFlag()
	fmt.Println(args)
	fmt.Println(command)
	fmt.Println(path)
	fmt.Println(data)
	fmt.Println(version)
	fmt.Println(flags)

	// testNode := "/mynode/test"
	timeout := 3
	zh, err := u.NewZH(servers, timeout)
	if err == nil {
		defer zh.Close()
	}

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
		// children
		paths, err := c.Children(zh, path)
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println(paths)
	case "set":
		// set
		nodeToSet := path
		// _, stat, err = c.Get(zh, nodeToSet)
		// if err != nil {
		// 	fmt.Println("get zk node fail when set:", err.Error())
		// 	return
		// }
		// Ver := stat.Version()
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
		// lock
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
		// 0: Permanent node
		c.Create(zh, path+"/test1", []byte("data"), 0)
		c.Create(zh, path+"/test2", []byte("data"), 0)
		// 2: sequence node
		c.Create(zh, path+"/se-job", []byte("data"), 2)
		c.Create(zh, path+"/se-job", []byte("data"), 2)
		paths, cerr := zh.Children(path)
		if cerr != nil {
			fmt.Println(cerr.Error())
		}
		fmt.Println("after create:", paths)

	case "watch":
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
		// delete child nodes node
		c.DelChildNode(zh, path)
	}
}
