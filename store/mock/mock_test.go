package mock

import (
	"fmt"
	"github.com/shiningacg/filestore/gateway"
	"log"
	"os"
	"testing"
	"time"
)

func TestApp(t *testing.T) {
	store := &Store{}
	g := gateway.NewGateway(":8888", store, log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile))
	go func() {
		for {
			fmt.Println(g.BandWidth())
			time.Sleep(time.Second * 10)
		}
	}()
	store.g = g
	err := g.Run()
	if err != nil {
		panic(err)
	}
}
