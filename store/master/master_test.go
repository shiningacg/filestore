package master

import (
	"context"
	"fmt"
	"github.com/shiningacg/filestore/store/common"
	"testing"
	"time"
)

func TestNewMaster(t *testing.T) {
	etcdConf := &common.EtcdConfig{
		EndPoint: []string{"127.0.0.1:2379"},
	}
	etcd := common.NewMaster(etcdConf, "/store/")
	master := NewMaster(etcd)
	etcd.SetHandler(master)
	go etcd.Run(context.TODO())
	for {
		time.Sleep(time.Second * 10)
		for id, _ := range master.stores {
			fmt.Println("store:", id)
		}
	}
}
