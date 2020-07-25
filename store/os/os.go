package os

import (
	"github.com/shiningacg/filestore/gateway"
	"github.com/shiningacg/filestore/store/remote"
	"log"
)

func NewOStore(config *StoreConfig, logger *log.Logger) *Store {
	g := gateway.NewGateway(config.GatewayAddr, logger)
	s := &Store{
		gateway:      g,
		storeManager: NewDefaultManager(config.StorePath),
		logger:       logger,
		db:           OpenBoltDB(config.StorePath+"/store.dat", logger),
	}
	g.SetStore(s)
	go g.Run()
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
	remote.Adder
}
