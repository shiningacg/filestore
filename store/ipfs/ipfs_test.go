package ipfs

import (
	"fmt"
	"github.com/shiningacg/filestore/gateway/checker"
	store2 "github.com/shiningacg/filestore/store"
	"github.com/shiningacg/filestore/store/common"
	"github.com/shiningacg/filestore/store/remote"
	"github.com/shiningacg/mygin-frame-libs/log"
	"testing"
	"time"
)

func TestNewStore(t *testing.T) {
	log.OpenLog(&log.Config{})
	store, err := NewStore(":8888", checker.MockChecker{}, log.Default())
	if err != nil {
		panic(err)
	}
	for {
		time.Sleep(time.Second * 3)
		fmt.Println(store.Gateway())
	}
}

func TestNewRemoteStoreServer(t *testing.T) {
	etcdConf := &common.EtcdConfig{
		EndPoint: []string{"127.0.0.1:2379"},
	}
	checker, err := checker.NewRedisChecker("127.0.0.1:6379", "")
	if err != nil {
		panic(err)
	}
	log.OpenLog(&log.Config{})
	store, err := NewStore(":8888", checker, log.Default())
	if err != nil {
		panic(err)
	}
	ss := remote.NewStoreGRPCServer("127.0.0.1:6666", store.(store2.FactoryStore).GetGateway(), remote.MockAdder{}, store, common.NewReporter(etcdConf))
	for {
		time.Sleep(time.Second * 3)
		fmt.Println(ss.Gateway())
	}
}
