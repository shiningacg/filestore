package os

import store "github.com/shiningacg/filestore"

type StoreManager interface {
	GetStorePath(file store.BaseFile) string
	GetBasePath() string
}

func NewDefaultManager(storePath string) *DefaultManager {
	if storePath[len(storePath)-1:] != "/" {
		storePath = storePath + "/"
	}
	return &DefaultManager{storePath: storePath}
}

type DefaultManager struct {
	storePath string
}

func (d *DefaultManager) GetStorePath(file store.BaseFile) string {
	return d.storePath + file.UUID() + "-" + file.Name()
}

func (d *DefaultManager) GetBasePath() string {
	return d.storePath
}
