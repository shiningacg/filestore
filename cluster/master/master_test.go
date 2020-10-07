package master

import (
	"context"
	"fmt"
	"github.com/shiningacg/filestore/cluster"
	"github.com/shiningacg/filestore/store/remote"
	"go.etcd.io/etcd/clientv3"
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
	store, err := remote.NewRemoteStore("127.0.0.1:8002")
	if err != nil {
		panic(err)
	}
	fmt.Println(store.Network())
	defer store.Close()
}
