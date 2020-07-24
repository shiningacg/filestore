package remote

import (
	"context"
	store "github.com/shiningacg/filestore"
	"github.com/shiningacg/filestore/store/remote/rpc"
)

type Stats Store

func (s Stats) Space() *store.Space {
	info, err := s.RemoteStoreClient.Space(context.TODO(), &rpc.Empty{})
	if err != nil {
		return &store.Space{}
	}
	return toStoreSpace(info)
}

func (s Stats) Network() *store.Network {
	info, err := s.RemoteStoreClient.Network(context.TODO(), &rpc.Empty{})
	if err != nil {
		return &store.Network{}
	}
	return toStoreNetwork(info)
}

func (s Stats) Gateway() *store.Bandwidth {
	info, err := s.RemoteStoreClient.Bandwidth(context.TODO(), &rpc.Empty{})
	if err != nil {
		return &store.Bandwidth{}
	}
	return toStoreBandwidth(info)
}
