package remote

import (
	"fmt"
	"github.com/shiningacg/filestore/gateway"
	"github.com/shiningacg/filestore/store/common"
	"github.com/shiningacg/filestore/store/mock"
	"log"
	"os"
	"testing"
	"time"
)

func TestNewStoreServer(t *testing.T) {
	g := gateway.NewGateway(":8888", log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile))
	etcdConf := &common.EtcdConfig{EndPoint: []string{"127.0.0.1:2379"}}
	store := mock.NewStore(g)
	g.SetStore(store)
	NewStoreGRPCServer("127.0.0.1:6666", MockAdder{}, store, common.NewReporter(etcdConf))
	for {
		fmt.Println(g.BandWidth())
		time.Sleep(time.Second * 10)
	}
}
