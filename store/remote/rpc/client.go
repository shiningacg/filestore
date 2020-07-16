package rpc

import "google.golang.org/grpc"

func NewStoreClient(addr string) (RemoteStoreClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return NewRemoteStoreClient(conn), nil
}
