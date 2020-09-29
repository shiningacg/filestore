package remote

import (
	"fmt"
	fs "github.com/shiningacg/filestore"
	"testing"
)

func testNewRemoteStore() *Store {
	s, err := NewRemoteStore("127.0.0.1:8002")
	if err != nil {
		panic(err)
	}
	return s
}

func TestStore_Add(t *testing.T) {
	var bf = &fs.BaseFileStruct{}
	bf.SetName("test.txt")
	bf.SetUUID("aaa")
	bf.SetSize(12)
	err := testNewRemoteStore().Add(bf)
	if err != nil {
		panic(err)
	}
}

func TestStore_Remove(t *testing.T) {
	var bf = &fs.BaseFileStruct{}
	bf.SetName("test.txt")
	bf.SetUUID("aaa")
	bf.SetSize(12)
	err := testNewRemoteStore().Remove(bf.UUID())
	if err != nil {
		panic(err)
	}
}

func TestStore_Get(t *testing.T) {
	var bf = &fs.BaseFileStruct{}
	bf.SetName("test.txt")
	bf.SetUUID("aaa")
	bf.SetSize(12)
	file, err := testNewRemoteStore().Get(bf.UUID())
	if err != nil {
		panic(err)
	}
	if file.UUID() != bf.UUID() {
		panic("not the same file")
	}
}

func TestStore_Network(t *testing.T) {
	network := testNewRemoteStore().Network()
	fmt.Println(network)
}

func TestStore_Gateway(t *testing.T) {
	network := testNewRemoteStore().Gateway()
	fmt.Println(network)
}
