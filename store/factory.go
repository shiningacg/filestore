package store

import (
	"context"
	fs "github.com/shiningacg/filestore"
	"github.com/shiningacg/filestore/gateway"
	"github.com/shiningacg/filestore/gateway/checker"
	"github.com/shiningacg/mygin-frame-libs/log"
)

type FactoryStore interface {
	fs.FileStore
	SetGateway(gateway gateway.Gateway)
	SetLogger(logger *log.Logger)
	GetGateway() gateway.Gateway
	GetLogger() *log.Logger
}

func NewStore(store FactoryStore, gatewayAddr string, checker checker.Checker, logger *log.Logger) fs.FileStore {
	gtw := gateway.NewMyginGateway(gatewayAddr, checker)
	store.SetGateway(gtw)
	gtw.SetStore(store)
	go func() {
		for {
			err := gtw.Run(context.TODO())
			if err != nil {
				panic(err)
			}
		}
	}()
	return store
}
