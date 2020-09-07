package follower

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/shiningacg/filestore/cluster"
	"testing"
	"time"
)

func TestNewRegister(t *testing.T) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints: []string{"127.0.0.1:2379"},
	})
	if err != nil {
		panic(err)
	}
	data := &cluster.Data{
		MetaData: cluster.MetaData{
			Id:      "center",
			Host:    []string{"127.0.0.1:5060"},
			Tag:     "lala",
			Weight:  10,
			Version: 1,
		},
		IsEntry: true,
		IsExit:  true,
		Cap:     0,
	}
	service := cluster.Service{
		Name: "svc.file",
		Id:   "center",
		Host: []string{"127.0.0.1:5060"},
		TTL:  time.Second * 3,
	}
	register := NewRegister(context.Background(), client, data, service)
	err = register.Register()
	if err != nil {
		panic(err)
	}
	data, err = register.GetData()
	if err != nil {
		panic(err)
	}
	fmt.Println(data)
	time.Sleep(time.Second * 10)
	register.Deregister()
}
