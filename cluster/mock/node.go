package mock

import (
	store "github.com/shiningacg/filestore"
	"github.com/shiningacg/filestore/cluster"
	mock2 "github.com/shiningacg/filestore/store/mock"
)

func NewNode() *Node {
	return &Node{
		data: cluster.Data{
			MetaData: cluster.MetaData{
				Id:      "mock",
				Host:    []string{"127.0.0.1:6666"},
				Tag:     "",
				Weight:  0,
				Version: 0,
			},
			GatewayAddr: "http://mockstore",
			Entry:       true,
			Exit:        true,
			Cap:         1024,
		},
		InfoStore: mock2.NewInfoStore(),
	}
}

type Node struct {
	data cluster.Data
	store.InfoStore
}

func (n *Node) ID() string {
	return n.data.Id
}

func (n *Node) Update(data *cluster.Data) error {
	n.data = *data
	return nil
}

func (n *Node) Data() cluster.Data {
	return n.data
}

func (n *Node) Entry(token string) string {
	return "http://mockstore/upload/" + token
}

func (n *Node) Exit(fid string) string {
	return "http://mockstore/download/" + fid
}
