package cluster

import (
	"encoding/json"
	"strings"
	"time"
)

type Service struct {
	Name string
	Id   string
	Host []string
	TTL  time.Duration
}

func (s Service) ToKey() string {
	return s.ToPath() + s.Id
}

func (s Service) ToPath() string {
	return "/" + strings.Join(strings.Split(s.Name, "."), "/") + "/"
}

// Data 存放在etcd，描述服务的一些信息
type Data struct {
	MetaData
	GatewayAddr string
	IsEntry     bool
	IsExit      bool
	Cap         uint64
}

func (d *Data) Encode() []byte {
	data, _ := json.Marshal(d)
	return data
}

func (d *Data) Decode(data []byte) error {
	return json.Unmarshal(data, d)
}

func (d *Data) Equal(data *Data) bool {
	var equal = true
	helper := func(stats bool) {
		if !stats {
			equal = stats
		}
	}
	helper(d.GatewayAddr == data.GatewayAddr)
	helper(d.IsExit == data.IsExit)
	helper(d.IsEntry == data.IsEntry)
	helper(d.Cap == data.Cap)
	helper(d.MetaData.Equal(data.MetaData))

	return equal
}

// 可重用的数据，打算以后抽离出一个tool包出来
type MetaData struct {
	Id      string
	Host    []string
	Tag     string
	Weight  uint8
	Version uint64
}

func (d MetaData) Equal(data MetaData) bool {
	var equal = true
	helper := func(stats bool) {
		if !stats {
			equal = stats
		}
	}
	helper(d.Tag == data.Tag)
	helper(d.Weight == data.Weight)
	helper(!d.IsHostChange(data))
	return equal
}

func (d MetaData) IsHostChange(data MetaData) bool {
	if len(d.Host) != len(data.Host) {
		return true
	}
	for i, addr := range d.Host {
		if addr != data.Host[i] {
			return true
		}
	}
	return false
}
