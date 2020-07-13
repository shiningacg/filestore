package mock

import (
	"fmt"
	store "github.com/shiningacg/filestore"
	"io/ioutil"
	"os"
)

type API Store

func (A API) Get(uuid string) (store.File, error) {
	f, err := os.Open("./mock.txt")
	if err != nil {
		return nil, err
	}
	return &File{
		id:   uuid,
		name: "mock",
		File: f,
	}, nil
}

func (A API) Add(file store.File) error {
	b, err := ioutil.ReadAll(file)
	fmt.Println(string(b))
	return err
}

func (A API) Remove(uuid string) error {
	return nil
}
