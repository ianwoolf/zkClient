package command

import (
	"fmt"

	u "github.com/ianwoolf/zkClient/util"
	zk "go.intra.xiaojukeji.com/golang/go-zookeeper/zk"
)

func Get(zh *u.ZH, path string) (content string, stat zk.Stat, err error) {
	var read []byte
	read, stat, err = zh.Get(path)
	if err == nil {
		content = string(read)
	}
	return
}

func Children(zh *u.ZH, path string) (paths []string, err error) {
	paths, err = zh.Children(path)
	return
}

func Set(zh *u.ZH, path string, data []byte, version int32) (stat zk.Stat, err error) {
	stat, err = zh.Set(path, data, version)
	return
}

func GetLock(zh *u.ZH, path string) *zk.Lock {
	return zh.NewLock(path)
}

func Create(zh *u.ZH, path string, data []byte, flags int32) (string, error) {
	return zh.Create(path, data, flags)
}

func WatchExist(zh *u.ZH, path string) (ok bool, event <-chan zk.Event) {
	ok, event = zh.ExistsW(path)
	return
}

func WatchChildren(zh *u.ZH, path string) (children []string, event <-chan zk.Event, err error) {
	children, event, err = zh.ChildrenW(path)
	return
}

func DelChildNode(zh *u.ZH, path string) (err error) {
	var childNodes []string
	var stat zk.Stat
	childNodes, err = zh.Children(path)
	if err != nil {
		fmt.Println("get node %s clindren error: %s", path, err.Error())
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
