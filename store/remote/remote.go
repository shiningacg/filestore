package remote

import (
	store "github.com/shiningacg/filestore"
	rpc "github.com/shiningacg/filestore/store/remote/rpc"
)

type Store struct {
	rpc.RemoteStoreClient
}

// 从etcd中获取
func (s Store) Stats() store.Stats {
	panic("implement me")
}

func (s Store) API() store.API {
	panic("implement me")
}
