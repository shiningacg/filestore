package os

import (
	store "github.com/shiningacg/filestore"
	"github.com/shiningacg/filestore/gateway"
	"log"
)

func NewOStore(config *StoreConfig, logger *log.Logger) *Store {
	s := &Store{
		gateway:      nil,
		storeManager: NewDefaultManager(config.StorePath),
		logger:       logger,
		db:           OpenBoltDB(config.StorePath, logger),
	}
	g := gateway.NewGateway(config.GatewayAddr, s.API())
	s.gateway = g
	return s
}

type StoreConfig struct {
	GatewayAddr string
	StorePath   string
}

type Store struct {
	gateway      *gateway.Gateway
	storeManager StoreManager
	logger       *log.Logger
	db           *BoltDB
}

func (s *Store) Stats() store.Stats {
	panic("implement me")
}

func (s *Store) API() store.API {
	return (*API)(s)
}
