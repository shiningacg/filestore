package master

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/shiningacg/filestore/cluster"
	"github.com/shiningacg/filestore/store/remote"
	"testing"
)

func TestNewMaster(t *testing.T) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints: []string{"127.0.0.1:2379"},
	})
	if err != nil {
		panic(err)
	}
	master := NewMaster(context.Background(), client, cluster.Service{
		Name: "svc.file",
	})
	var evts = make(chan *cluster.Event, 5)
	master.Watcher.Events(evts)
	fmt.Println("hi")
	for evt := range evts {
		fmt.Println(evt)
	}
}

func TestNewRemoteStore(t *testing.T) {
	store, err := remote.NewRemoteStore("127.0.0.1:8888")
	if err != nil {
		panic(err)
	}
	defer store.Close()
}
