package mock

import (
	store "github.com/shiningacg/filestore"
	"github.com/shiningacg/filestore/cluster"
	mc "github.com/shiningacg/filestore/store/mock"
)

func NewNode() *Node {
	return &Node{
		data:      cluster.Data{},
		InfoStore: mc.NewInfoStore(),
	}
}

type Node struct {
	data cluster.Data
	store.InfoStore
}

func (n *Node) ID() string {
	return "node"
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
