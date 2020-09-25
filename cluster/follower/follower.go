package follower

import (
	"context"
	"github.com/coreos/etcd/clientv3"
	"github.com/shiningacg/filestore/cluster"
)

func NewFollower(ctx context.Context, client *clientv3.Client, data cluster.Data, service cluster.Service) (*Follower, error) {
	register := NewRegister(ctx, client, &data, service)
	watcher := cluster.NewWatcher(client, service.ToKey(), false)
	err := register.Register()
	if err != nil {
		return nil, err
	}
	watcher.Watch(ctx)
	return &Follower{
		Data:     data,
		Service:  service,
		Register: register,
		Watcher:  watcher,
	}, nil
}

type Follower struct {
	cluster.Data
	cluster.Service
	Register
	cluster.Watcher
}
