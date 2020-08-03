package mock

import (
	"fmt"
	"github.com/shiningacg/filestore/gateway"
	"github.com/shiningacg/mygin-frame-libs/log"
	"testing"
	"time"
)

func MewStore() *Store {
	log.OpenLog(&log.Config{})
	g := gateway.NewGateway(":8888", gateway.MockChecker{}, log.Default())
	go func() {
		for {
			fmt.Println(g.BandWidth())
			time.Sleep(time.Second * 10)
		}
	}()
	store := &Store{g: g}
	g.SetStore(store)
	go func() {
		err := g.Run()
		if err != nil {
			panic(err)
		}
	}()
	return store
}

func TestApp(t *testing.T) {
	log.OpenLog(&log.Config{})
	g := gateway.NewGateway(":1111", gateway.MockChecker{}, log.Default())
	go func() {
		for {
			fmt.Println(g.BandWidth())
			time.Sleep(time.Second * 10)
		}
	}()
	store := &Store{g: g}
	g.SetStore(store)
	err := g.Run()
	if err != nil {
		panic(err)
	}
}
