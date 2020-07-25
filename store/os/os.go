package os

import (
	"github.com/shiningacg/filestore/gateway"
	"github.com/shiningacg/filestore/store/common"
	"log"
)

func NewOStore(config *StoreConfig, logger *log.Logger) *Store {
	s := &Store{
		gateway:      nil,
		storeManager: NewDefaultManager(config.StorePath),
		logger:       logger,
		db:           OpenBoltDB(config.StorePath+"/store.dat", logger),
	}
	g := gateway.NewGateway(config.GatewayAddr, s, logger)
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
	common.Adder
}
