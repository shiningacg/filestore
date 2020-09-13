package remote

import (
	"fmt"
	"github.com/shiningacg/filestore/gateway"
	"github.com/shiningacg/filestore/store/common"
	"github.com/shiningacg/filestore/store/mock"
	"github.com/shiningacg/mygin-frame-libs/log"
	"testing"
	"time"
)

func TestNewStoreServer(t *testing.T) {
	log.OpenLog(&log.Config{})
	g := gateway.NewGateway(":8888", gateway.MockChecker{}, log.Default())
	etcdConf := &common.EtcdConfig{EndPoint: []string{"127.0.0.1:2379"}}
	store := mock.NewFileStore(g)
	g.SetStore(store)
	NewStoreGRPCServer("127.0.0.1:5060", g, MockAdder{}, store, common.NewReporter(etcdConf))
	for {
		fmt.Println(g.BandWidth())
		time.Sleep(time.Second * 10)
	}
}

func TestNewStoreServerWithRedisChecker(t *testing.T) {
	log.OpenLog(&log.Config{})
	checker, err := gateway.NewRedisChecker("127.0.0.1:6379", "")
	if err != nil {
		panic(err)
	}
	g := gateway.NewGateway(":8888", checker, log.Default())
	etcdConf := &common.EtcdConfig{EndPoint: []string{"127.0.0.1:2379"}}
	store := mock.NewFileStore(g)
	g.SetStore(store)
	NewStoreGRPCServer("127.0.0.1:6666", g, MockAdder{}, store, common.NewReporter(etcdConf))
	for {
		fmt.Println(g.BandWidth())
		time.Sleep(time.Second * 10)
	}
}
