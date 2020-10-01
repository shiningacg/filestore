package store

import (
	"context"
	"github.com/shiningacg/filestore/store/remote"
	"testing"
)

const (
	CheckerAddr         = "127.0.0.1:5040"
	GatewayAddr         = "0.0.0.0:8001"
	GrpcAddr            = "0.0.0.0:8002"
	AnnounceGatewayAddr = "127.0.0.1:8001"
	AnnounceGrpcAddr    = "127.0.0.1:8002"
)

var config = Config{
	Type:    "ipfs",
	Path:    "/",
	Gateway: GatewayAddr,
	Plugin: Plugin{
		AdderAddr:   "",
		CheckerAddr: CheckerAddr,
	},
}

func TestNewStore(t *testing.T) {
	store, err := NewStore(config)
	if err != nil {
		panic(err)
	}
	store = store
	store.Run(context.Background())
}

func TestNewRemoteStore(t *testing.T) {
	store, err := NewStore(config)
	if err != nil {
		panic(err)
	}
	remote.NewStoreGRPCServer(GrpcAddr, remote.MockAdder{}, store)
	store.Run(context.Background())
}
