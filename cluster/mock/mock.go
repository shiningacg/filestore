package mock

import "github.com/shiningacg/filestore/cluster"

func NewMaster() Master {
	mock := NewNode()
	nodes := make(cluster.Nodes, 1)
	nodes[0] = mock
	return Master{node: mock, nodes: nodes}
}

type Master struct {
	node  *Node
	nodes cluster.Nodes
}

func (m Master) Node(id string) cluster.Node {
	return m.node
}

func (m Master) Nodes() cluster.Nodes {
	return m.nodes
}

func (m Master) Entries() cluster.Nodes {
	return m.nodes
}

func (m Master) Exits() cluster.Nodes {
	return m.nodes
}

func (m Master) BestEntry(size uint64) cluster.Node {
	return m.node
}

func (m Master) BestExit(fid string) cluster.Node {
	return m.node
}
