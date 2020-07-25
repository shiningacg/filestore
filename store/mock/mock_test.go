package mock

import (
	"fmt"
	"github.com/shiningacg/filestore/gateway"
	"log"
	"os"
	"testing"
	"time"
)

func MewStore() *Store {
	g := gateway.NewGateway(":8888", log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile))
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
	g := gateway.NewGateway(":8888", log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile))
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
