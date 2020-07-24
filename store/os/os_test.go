package os

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"
)

func testOpenStore() *Store {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	return NewOStore(&StoreConfig{
		GatewayAddr: ":8888",
		StorePath:   ".",
	}, logger)
}

func TestNewOStore(t *testing.T) {
	store := testOpenStore()
	go func() {
		for {
			time.Sleep(time.Second * 2)
			fmt.Println(store.Stats().Gateway())
		}
	}()
	err := store.gateway.Run()
	if err != nil {
		panic(err)
	}
}

func TestAPI_Add(t *testing.T) {
	f, _ := os.Open("./aaa")
	store := testOpenStore()
	err := store.API().Add(&File{
		name: "aaa",
		id:   "aaac",
		File: f,
	})
	if err != nil {
		panic(err)
	}
}
