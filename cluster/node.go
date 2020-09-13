package cluster

import (
	fs "github.com/shiningacg/filestore"
)

type Nodes []Node

// TODO: 添加一些集体调用的代码
func (n Nodes) SortBest(better func(n1, n2 Node) bool) Node {
	var best Node
	if len(n) == 0 {
		return nil
	}
	best = n[0]
	for i := 1; i < len(n)-1; i++ {
		if !better(best, n[i]) {
			best = n[i]
		}
	}
	return best
}

func (n Nodes) BestUpload() Node {
	return n.SortBest(func(n1, n2 Node) bool {
		return n1.Network().Upload < n2.Network().Upload
	})
}

func (n Nodes) BestDownload() Node {
	return n.SortBest(func(n1, n2 Node) bool {
		return n1.Network().Download < n2.Network().Download
	})
}

func (n Nodes) BestSpace() Node {
	return n.SortBest(func(n1, n2 Node) bool {
		return n1.Space().Free > n2.Space().Free
	})
}

func (n Nodes) Node(id string) Node {
	for _, v := range n {
		if v.ID() == id {
			return v
		}
	}
	return nil
}

func (n Nodes) Delete(id string) Nodes {
	for i, v := range n {
		if v.ID() == id {
			return append(n[:i], n[i+1:]...)
		}
	}
	return n
}

// TODO:删除node同时断开grpc
func (n Nodes) Destroy(id string) Nodes {
	// node := n.Node(id)
	return n.Delete(id)
}

// Node 实际操作的节点对象
type Node interface {
	ID() string
	Update(data *Data) error
	Data() Data
	Entry(token string) string
	Exit(fid string) string
	fs.InfoStore
}
