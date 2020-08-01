package os

import (
	"fmt"
	fs "github.com/shiningacg/filestore"

	"github.com/shiningacg/mygin-frame-libs/log"
	"os"
	"testing"
	"time"
)

func testOpenStore() *Store {
	log.OpenLog(&log.Config{})
	logger := log.Default()
	return NewOStore(&StoreConfig{
		GatewayAddr: ":8887",
		StorePath:   ".",
	}, logger)
}

func TestNewOStore(t *testing.T) {
	store := testOpenStore()
	go func() {
		for {
			time.Sleep(time.Second * 2)
			fmt.Println(store.Gateway())
		}
	}()
	err := store.gateway.Run()
	if err != nil {
		panic(err)
	}
}

func TestAPI_Add(t *testing.T) {
	var bs = &fs.BaseFileStruct{}
	f, _ := os.Open("./aaa")
	bs.SetName("aaa")
	bs.SetSize(100)
	store := testOpenStore()
	err := store.Add(fs.NewReadableFile(bs, f))
	if err != nil {
		panic(err)
	}
}
