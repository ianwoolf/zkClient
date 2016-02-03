package main

import (
	"fmt"
	// "sort"

	u "github.com/ianwoolf/zkClient/util"
	zk "go.intra.xiaojukeji.com/golang/go-zookeeper/zk"
)

func get(zh *u.ZH, path string) (content string, stat zk.Stat, err error) {
	var read []byte
	read, stat, err = zh.Get(path)
	if err == nil {
		content = string(read)
	}
	return
}

func children(zh *u.ZH, path string) (paths []string, err error) {
	paths, err = zh.Children(path)
	return
}

func set(zh *u.ZH, path string, data []byte, version int32) (stat zk.Stat, err error) {
	stat, err = zh.Set(path, data, version)
	return
}

func getLock(zh *u.ZH, path string) *zk.Lock {
	return zh.NewLock(path)
}

func create(zh *u.ZH, path string, data []byte, flags int32) (string, error) {
	return zh.Create(path, data, flags)
}

func watchExist(zh *u.ZH, path string) (ok bool, event <-chan zk.Event) {
	ok, event = zh.ExistsW(path)
	return
}

func watchChildren(zh *u.ZH, path string) (children []string, event <-chan zk.Event, err error) {
	children, event, err = zh.ChildrenW(path)
	return
}

func delChildNode(zh *u.ZH, path string) (err error) {
	var childNodes []string
	var stat zk.Stat
	childNodes, err = zh.Children(path)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("begin to del node:", childNodes)
	for _, node := range childNodes {
		nodeToDel := fmt.Sprintf("%s/%s", path, node)
		_, stat, err = zh.Get(nodeToDel)
		if err != nil {
			fmt.Println("get zk node fail:", nodeToDel, err.Error())
			return
		} else {
			fmt.Println("begin to del node:", nodeToDel)
		}
		err = zh.Delete(nodeToDel, int32(stat.Version()))
		if err != nil {
			fmt.Printf("del node %s error: %s", nodeToDel, err.Error())
			continue
		}
	}
	return
}
func main() {
	testNode := "/mynode/test"
	servers := []string{"106.186.127.250:2181"}
	timeout := 3
	zh, err := u.NewZH(servers, timeout)
	if err == nil {
		defer zh.Close()
	}

	content, stat, err := get(zh, testNode)
	if err != nil {
		fmt.Println("get zk node fail:", testNode, err.Error())
	} else {
		fmt.Println(content)
		fmt.Println(stat.Version(), stat.MTime())
	}

	// children
	paths, err := children(zh, testNode)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(paths)

	// set
	nodeToSet := testNode + "/test1"
	_, stat, err = get(zh, nodeToSet)
	if err != nil {
		fmt.Println("get zk node fail when set:", err.Error())
		return
	}
	Ver := stat.Version()
	_, errSet := set(zh, nodeToSet, []byte("test set2"), int32(Ver))
	if errSet != nil {
		fmt.Println(errSet.Error())
	}
	content, _, err = get(zh, nodeToSet)
	if err != nil {
		fmt.Println("get zk node fail when set:", err.Error())
		return
	}
	fmt.Println("after set:", content)

	// lock
	fmt.Println("begin to lock")
	lock := getLock(zh, nodeToSet)
	lock.Lock()
	content, stat, err = get(zh, nodeToSet)
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
	content, _, err = get(zh, nodeToSet)
	if err != nil {
		fmt.Println("get zk node fail in lock:", err.Error())
		return
	}
	fmt.Println("after set in lock:", content)
	lock.Unlock()

	// create
	// 0: Permanent node
	create(zh, testNode+"/test1", []byte("data"), 0)
	create(zh, testNode+"/test2", []byte("data"), 0)
	// 2: sequence node
	create(zh, testNode+"/se-job", []byte("data"), 2)
	create(zh, testNode+"/se-job", []byte("data"), 2)
	paths, cerr := zh.Children(testNode)
	if cerr != nil {
		fmt.Println(cerr.Error())
	}
	fmt.Println("after create:", paths)

	// watch
	fmt.Println("check watch on testNode:")
	existOk, Eevent := watchExist(zh, testNode)
	if existOk {
		fmt.Println(testNode, "exist and watch")
	}
	// _, Gevent, GWerr := zh.GetW(testNode)
	// if GWerr == nil {
	// 	fmt.Println(testNode, "get zk node success and watch")
	// }
	_, Cevent, CWerr := watchChildren(zh, testNode)
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

	// delete child nodes node
	delChildNode(zh, testNode)
}
