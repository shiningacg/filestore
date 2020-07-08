package os

import store "github.com/shiningacg/filestore"

type StoreManager interface {
	GetStorePath(file store.File) string
}

type DefaultManager struct {
	storePath string
}

func (d DefaultManager) GetStorePath(file store.File) string {
	return d.storePath + file.ID() + "-" + file.FileName()
}
