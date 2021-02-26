package main

import (
	"context"
	"github.com/shiningacg/filestore/cluster"
	"github.com/shiningacg/filestore/cluster/follower"
	"github.com/shiningacg/filestore/store"
	"github.com/shiningacg/filestore/store/remote"
	"go.etcd.io/etcd/clientv3"
)

const (
	CheckerAddr         = "192.168.1.14:8002"
	GatewayAddr         = "0.0.0.0:8005"
	GrpcAddr            = "0.0.0.0:8006"
	AnnounceGatewayAddr = "192.168.1.14:8005"
	AnnounceGrpcAddr    = "192.168.1.14:8006"
)

var config = follower.Config{
	Store: store.Config{
		Type:    "ipfs",
		Path:    "/",
		Gateway: GatewayAddr,
		Plugin: store.Plugin{
			AdderAddr:   "",
			CheckerAddr: CheckerAddr,
		},
	},
	Etcd: []string{"192.168.1.14:2379"},
	Service: cluster.Service{
		Name: "svc.file",
		Id:   "test",
		TTL:  3,
	},
	Data: cluster.Data{
		MetaData: cluster.MetaData{
			Id:      "test",
			Host:    []string{AnnounceGrpcAddr},
			Tag:     "",
			Weight:  0,
			Version: 0,
		},
		GatewayAddr: AnnounceGatewayAddr,
		Entry:       true,
		Exit:        true,
		Cap:         1024 * 1024 * 1024 * 5,
	},
}

// ipfs store
func main() {
	st, err := store.NewStore(config.ToStoreConfig())
	if err != nil {
		panic(err)
	}
	adder := remote.MockAdder{}
	// 启动grpc服务
	remote.NewStoreGRPCServer(GrpcAddr, adder, st)
	// 连接集群
	etcd := ConnectEtcd()
	//
	app := follower.NewFollower(etcd, config.Data, config.Service)
	go st.Run(context.Background())
	app.Run()
}

func ConnectEtcd() *clientv3.Client {
	cl, err := clientv3.New(clientv3.Config{
		Endpoints: []string{"192.168.1.14:2379"},
	})
	if err != nil {
		panic(err)
	}
	return cl
}
