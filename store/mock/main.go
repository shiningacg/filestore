package mock

import store "github.com/shiningacg/filestore"

type Store struct{}

func (s *Store) Stats() store.Stats {
	return (*Stats)(s)
}

func (s *Store) API() store.API {
	return (*API)(s)
}
