package os

import (
	store "github.com/shiningacg/filestore"
	"io"
)

type API struct {
}

func (A API) Get(uuid string) (store.File, error) {
	panic("implement me")
}

func (A API) Add(reader io.Reader) error {
	panic("implement me")
}

func (A API) Remove(uuid string) error {
	panic("implement me")
}
