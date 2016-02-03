package util

import (
	// "sort"
	"time"

	zk "go.intra.xiaojukeji.com/golang/go-zookeeper/zk"
)

type ZH struct {
	conn *zk.Conn
}
type (
// Stat zk.Stat
// Lock zk.Lock
)

var (
	ZK_ACL = zk.WorldACL(zk.PermAll)
)

func NewZH(servers []string, timeout int) (zh *ZH, err error) {
	_conn, _, err := zk.Connect(servers, time.Duration(timeout)*time.Second)
	if err != nil {
		return
	}

	zh = &ZH{
		conn: _conn,
	}

	return
}

func (z *ZH) Close() {
	z.conn.Close()
}

func (z *ZH) isRecoverable(err error) bool {
	return err == zk.ErrConnectionClosed || err == zk.ErrNoServer || err == zk.ErrSessionExpired
}

func (z *ZH) Create(path string, data []byte, flags int32) (string, error) {
	for {
		path, err := z.conn.Create(path, data, flags, ZK_ACL)
		if z.isRecoverable(err) {
			continue
		}
		return path, err
	}
}

func (z *ZH) Delete(path string, version int32) error {
	for {
		err := z.conn.Delete(path, version)
		if z.isRecoverable(err) {
			continue
		}
		return err
	}
}

func (z *ZH) Exists(path string) bool {
	for {
		exist, _, err := z.conn.Exists(path)
		if err != nil {
			continue
		}
		return exist
	}
}

func (z *ZH) Get(path string) ([]byte, zk.Stat, error) {
	for {
		data, stat, err := z.conn.Get(path)
		if z.isRecoverable(err) {
			continue
		}
		return data, stat, err
	}
}

func (z *ZH) Set(path string, data []byte, version int32) (zk.Stat, error) {
	for {
		stat, err := z.conn.Set(path, data, version)
		if z.isRecoverable(err) {
			continue
		}
		return stat, err
	}
}

func (z *ZH) Children(path string) ([]string, error) {
	for {
		paths, _, err := z.conn.Children(path)
		if z.isRecoverable(err) {
			continue
		}
		return paths, err
	}
}

func (z *ZH) ExistsW(path string) (bool, <-chan zk.Event) {
	for {
		exist, _, event, err := z.conn.ExistsW(path)
		if err != nil {
			continue
		}
		return exist, event
	}
}

func (z *ZH) GetW(path string) ([]byte, <-chan zk.Event, error) {
	for {
		data, _, event, err := z.conn.GetW(path)
		if z.isRecoverable(err) {
			continue
		}
		return data, event, err
	}
}

func (z *ZH) ChildrenW(path string) ([]string, <-chan zk.Event, error) {
	for {
		paths, _, events, err := z.conn.ChildrenW(path)
		if z.isRecoverable(err) {
			continue
		}
		return paths, events, err
	}
}

func (z *ZH) NewLock(path string) *zk.Lock {
	for {
		lock := zk.NewLock(z.conn, path, ZK_ACL)
		return lock
	}
	// func (l *Lock) Lock() error
	// It will wait to return until the lock is acquired or an error occurs.
	// If this instance already has the lock then ErrDeadlock is returned.

	// func (l *Lock) Unlock() error
	// Unlock releases an acquired lock.
	// If the lock is not currently acquired by this Lock instance than ErrNotLocked is returned.
}
