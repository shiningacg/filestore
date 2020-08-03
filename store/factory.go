package store

import (
	fs "github.com/shiningacg/filestore"
	"github.com/shiningacg/filestore/gateway"
	"github.com/shiningacg/mygin-frame-libs/log"
)

type FactoryStore interface {
	fs.FileStore
	SetGateway(gateway *gateway.Gateway)
	SetLogger(logger *log.Logger)
	GetGateway() *gateway.Gateway
	GetLogger() *log.Logger
}

func NewStore(store FactoryStore, gatewayAddr string, checker gateway.Checker, logger *log.Logger) fs.FileStore {
	gtw := gateway.NewGateway(gatewayAddr, checker, logger)
	store.SetGateway(gtw)
	gtw.SetStore(store)
	go func() {
		for {
			err := gtw.Run()
			if err != nil {
				panic(err)
			}
		}
	}()
	return store
}
