package master

import (
	"errors"
	"fmt"
	fs "github.com/shiningacg/filestore"
	"github.com/shiningacg/filestore/store/common"
	"github.com/shiningacg/filestore/store/remote"
)

func NewMaster(master *common.Master) *Master {
	return &Master{stores: make(map[string]fs.InfoStore), etcd: master}
}

type Master struct {
	stores map[string]fs.InfoStore
	etcd   *common.Master
}

func (m *Master) Get(uuid string) (fs.BaseFile, error) {
	var store fs.InfoFS
	var network = &fs.Network{}
	for _, s := range m.stores {
		n := s.Network()
		if n.Upload <= network.Upload {
			store = s
		}
	}
	if store == nil {
		return nil, errors.New("没有存储节点在线")
	}
	return store.Get(uuid)
}

func (m *Master) Add(file fs.BaseFile) error {
	// 进行高级处理,当前只是核心仓库存储一次
	store, has := m.stores["center"]
	if has {
		return store.Add(file)
	}
	return errors.New("中心仓库节点不在线")
}

func (m *Master) Remove(uuid string) error {
	var e error
	for _, store := range m.stores {
		err := store.Remove(uuid)
		if err != nil {
			e = err
		}
	}
	return e
}

// 同理可得
func (m *Master) Space() *fs.Space {
	panic("implement me")
}

// 同理可得
func (m *Master) Network() *fs.Network {
	panic("implement me")
}

// 统计总流量
func (m *Master) Gateway() *fs.Bandwidth {
	var bd = &fs.Bandwidth{}
	for _, store := range m.stores {
		sbd := store.Gateway()
		bd.Bandwidth += sbd.Bandwidth
		bd.DayBandwidth += sbd.DayBandwidth
		bd.HourBandwidth += sbd.HourBandwidth
		bd.Visit += sbd.Visit
		bd.DayVisit += sbd.DayVisit
		bd.HourVisit += sbd.HourVisit
	}
	return bd
}

func (m *Master) Offline(info *common.NodeInfo) {
	// TODO: 如何关闭连接
	delete(m.stores, info.NodeId)
}

func (m *Master) Online(info *common.NodeInfo) {
	m.Offline(info)
	store, err := remote.NewRemoteStore(info.GRPCAddr)
	if err != nil {
		fmt.Println(err)
	}
	m.stores[info.NodeId] = store
}
