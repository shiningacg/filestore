package remote

import (
	"filesys"
)

type RemoteStore struct {
}

func (r RemoteStore) Stats() filestore.Stats {
	panic("implement me")
}

func (r RemoteStore) API() filestore.API {
	panic("implement me")
}
