package remote

import (
	"fmt"
	"github.com/shiningacg/filestore/gateway"
	"github.com/shiningacg/filestore/gateway/checker"
	"github.com/shiningacg/filestore/store/mock"
	"github.com/shiningacg/mygin-frame-libs/log"
	"testing"
	"time"
)

func TestNewStoreServer(t *testing.T) {
	log.OpenLog(&log.Config{})
	g := gateway.NewMyginGateway(":8888", checker.MockChecker{})
	store := mock.NewFileStore(g)
	g.SetStore(store)
	NewStoreGRPCServer("127.0.0.1:5060", MockAdder{}, store)
	for {
		fmt.Println(g.BandWidth())
		time.Sleep(time.Second * 10)
	}
}

// TODO： 重写！！！
func TestNewStoreServerWithRedisChecker(t *testing.T) {
	log.OpenLog(&log.Config{})
	checker, err := checker.NewRedisChecker("127.0.0.1:6379", "")
	if err != nil {
		panic(err)
	}
	g := gateway.NewMyginGateway(":8888", checker)
	store := mock.NewFileStore(g)
	g.SetStore(store)
	NewStoreGRPCServer("127.0.0.1:6666", MockAdder{}, store)
	for {
		fmt.Println(g.BandWidth())
		time.Sleep(time.Second * 10)
	}
}
