package remote

import store "github.com/shiningacg/filestore"

type Stats Store

func (s Stats) Space() *store.Space {
	panic("implement me")
}

func (s Stats) Network() *store.Network {
	panic("implement me")
}

func (s Stats) Bandwidth() *store.Gateway {
	panic("implement me")
}
