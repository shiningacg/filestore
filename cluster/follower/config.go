package follower

import (
	"github.com/shiningacg/filestore/cluster"
	"github.com/shiningacg/filestore/store"
)

type Config struct {
	Store store.Config
	// 集群地址
	Etcd []string `json:"etcd"`
	// 集群相关信息
	cluster.Service
	cluster.Data
}

func (c *Config) ToStoreConfig() store.Config {
	return c.Store
}
