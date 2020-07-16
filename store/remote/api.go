package remote

import (
	"context"
	store "github.com/shiningacg/filestore"
	"github.com/shiningacg/filestore/store/remote/rpc"
)

type API Store

func (a API) Get(uuid string) (file store.File, err error) {
	pbFile, err := a.RemoteStoreClient.Get(context.TODO(), &rpc.UUID{UUID: uuid})
	if err != nil {
		return nil, err
	}
	return wrapPBFile(pbFile), nil
}

func (a API) Add(file store.File) error {
	_, err := a.RemoteStoreClient.Add(context.TODO(), toPBFile(file))
	return err
}

func (a API) Remove(uuid string) error {
	_, err := a.RemoteStoreClient.Remove(context.TODO(), &rpc.UUID{UUID: uuid})
	return err
}
