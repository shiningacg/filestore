package cluster

import (
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"testing"
)

func TestWatcher_Exist(t *testing.T) {
	cl, err := clientv3.New(clientv3.Config{
		Endpoints: []string{"127.0.0.1:2379"},
	})
	if err != nil {
		panic(err)
	}
	watcher := NewWatcher(cl, "/svc/file", true)
	evt, err := watcher.Exist()
	if err != nil {
		panic(err)
	}
	fmt.Println(evt)
}
