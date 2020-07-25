package remote

import (
	"fmt"
	"github.com/shiningacg/filestore/gateway"
	"github.com/shiningacg/filestore/store/mock"
	"log"
	"os"
	"testing"
	"time"
)

func TestNewStoreServer(t *testing.T) {
	g := gateway.NewGateway(":8888", log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile))
	store := mock.NewStore(g)
	g.SetStore(store)
	NewStoreGRPCServer(":6666", MockAdder{}, store)
	for {
		fmt.Println(g.BandWidth())
		time.Sleep(time.Second * 10)
	}
}
