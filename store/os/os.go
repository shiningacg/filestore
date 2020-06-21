package os

import store "github.com/shiningacg/filestore"

type OSStore struct {
	*OSStoreStats
}

func (O OSStore) Stats() store.Stats {
	return O.OSStoreStats
}

func (O OSStore) API() store.API {
	panic("implement me")
}

type OSStoreStats struct{}

func (O OSStoreStats) Space() store.Space {
	panic("implement me")
}

func (O OSStoreStats) Bandwidth() store.Bandwidth {
	panic("implement me")
}

func (O OSStoreStats) Gateway() store.Gateway {
	panic("implement me")
}
