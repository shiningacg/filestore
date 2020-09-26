package mock

import (
	"context"
	"fmt"
	"github.com/shiningacg/filestore/gateway"
	"github.com/shiningacg/filestore/gateway/checker"
	"github.com/shiningacg/mygin-frame-libs/log"
	"testing"
	"time"
)

func MewStore() *FileStore {
	log.OpenLog(&log.Config{})
	g := gateway.NewDefaultGateway(":8888", checker.MockChecker{}, log.Default())
	go func() {
		for {
			fmt.Println(g.BandWidth())
			time.Sleep(time.Second * 10)
		}
	}()
	store := &FileStore{g: g}
	g.SetStore(store)
	go func() {
		err := g.Run(context.TODO())
		if err != nil {
			panic(err)
		}
	}()
	return store
}

func NewStoreWithMygin() *FileStore {
	g := gateway.NewMyginGateway(":8888", checker.MockChecker{})
	store := &FileStore{g: g}
	g.SetStore(store)
	go func() {
		err := g.Run(context.TODO())
		if err != nil {
			panic(err)
		}
	}()
	return store
}

func TestApp(t *testing.T) {
	g := NewStoreWithMygin()
	for {
		time.Sleep(time.Second * 4)
		fmt.Println(g.Gateway())
	}
}
