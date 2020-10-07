package follower

import (
	"context"
	"fmt"
	"github.com/shiningacg/filestore/cluster"
	"go.etcd.io/etcd/clientv3"
	"testing"
	"time"
)

var client *clientv3.Client

var data = cluster.Data{
	MetaData: cluster.MetaData{
		Id:      "center",
		Host:    []string{"127.0.0.1:5060"},
		Tag:     "lala",
		Weight:  10,
		Version: 1,
	},
	Entry: true,
	Exit:  true,
	Cap:   0,
}

var service = cluster.Service{
	Name: "svc.file",
	Id:   "center",
	TTL:  time.Second * 3,
}

func testInit() {
	c, err := clientv3.New(clientv3.Config{
		Endpoints: []string{"127.0.0.1:2379"},
	})
	if err != nil {
		panic(err)
	}
	client = c
}

func TestNewRegister(t *testing.T) {
	testInit()
	register := NewRegister(context.Background(), client, &data, service)
	err := register.Register()
	if err != nil {
		panic(err)
	}
	data, err := register.GetData()
	if err != nil {
		panic(err)
	}
	fmt.Println(data)
	time.Sleep(time.Second * 10)
	register.Deregister()
}

func TestNewFollower(t *testing.T) {
	testInit()
	follower := NewFollower(client, data, service)
	follower.Run()
}
