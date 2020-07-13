package remote

import store "github.com/shiningacg/filestore"

type Store struct {
}

// 从etcd中获取
func (s Store) Stats() store.Stats {
	panic("implement me")
}

func (s Store) API() store.API {
	panic("implement me")
}
