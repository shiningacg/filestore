package mock

import (
	"github.com/shiningacg/filestore/cluster"
	"time"
)

func NewMaster() *Master {
	mock := NewNode()
	nodes := make(cluster.Nodes, 1)
	nodes[0] = mock
	repo := make([]chan<- *cluster.Event, 0, 3)
	master := &Master{node: mock, nodes: nodes, repo: repo}
	go master.watch()
	return master
}

type Master struct {
	node  *Node
	nodes cluster.Nodes
	repo  []chan<- *cluster.Event
}

func (m *Master) Node(id string) cluster.Node {
	return m.node
}

func (m *Master) Nodes() cluster.Nodes {
	return m.nodes
}

func (m *Master) Entries() cluster.Nodes {
	return m.nodes
}

func (m *Master) Exits() cluster.Nodes {
	return m.nodes
}

func (m *Master) BestEntry(size uint64) cluster.Node {
	return m.node
}

func (m *Master) BestExit(fid string) cluster.Node {
	return m.node
}

func (m *Master) Watch(repo chan<- *cluster.Event) {
	m.repo = append(m.repo, repo)
}

func (m *Master) watch() {
	for {
		time.Sleep(time.Second * 5)
		for _, ch := range m.repo {
			data := m.node.Data()
			ch <- cluster.NewEvent(&data, cluster.PUT)
			time.Sleep(time.Second * 5)
			ch <- cluster.NewEvent(&data, cluster.DEL)
		}
	}
}
