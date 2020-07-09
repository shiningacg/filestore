package os

import store "github.com/shiningacg/filestore"

type StoreManager interface {
	GetStorePath(file store.File) string
}

func NewDefaultManager(storePath string) *DefaultManager {
	return &DefaultManager{storePath: storePath}
}

type DefaultManager struct {
	storePath string
}

func (d *DefaultManager) GetStorePath(file store.File) string {
	return d.storePath + file.ID() + "-" + file.FileName()
}
