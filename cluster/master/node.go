package master

import (
	"fmt"
	fs "github.com/shiningacg/filestore"
	"github.com/shiningacg/filestore/cluster"
	"github.com/shiningacg/filestore/store/remote"
	"time"
)

type Nodes []*Node

func (n Nodes) SortBest(less func(n1, n2 *Node) bool) *Node {
	var best *Node
	if len(n) == 0 {
		return nil
	}
	best = n[0]
	for i := 1; i < len(n)-1; i++ {
		if !less(best, n[i]) {
			best = n[i]
		}
	}
	return best
}

func (n Nodes) Node(id string) *Node {
	for _, v := range n {
		if v.Id == id {
			return v
		}
	}
	return nil
}

// TODO：在删除前进行grpc断开的处理
func (n Nodes) Delete(id string) Nodes {
	for i, v := range n {
		if v.Id == id {
			return append(n[:i], n[i+1:]...)
		}
	}
	return n
}

// NewNode 通过给定的地址创建一个node
func NewNode(host string) (*Node, error) {
	store, err := remote.NewRemoteStore(host)
	if err != nil {
		return nil, err
	}
	return &Node{
		Store: store,
	}, nil
}

// NewNodeFromData 通过Data更新一个node
func NewNodeFromData(data *cluster.Data) (*Node, error) {
	var node = &Node{}
	return node, node.Update(data)
}

// Node 实际操作的节点对象
type Node struct {
	*remote.Store
	cluster.Data
	// 缓存信息
	network    *fs.Network
	lastUpdate time.Time
}

// Update 更新节点的信息，如果地址发送了改变那么会重新建立grpc连接
func (n *Node) Update(node *cluster.Data) error {
	// 如果监听地址变化了，那么就需要重新加载
	if n.IsHostChange(node.MetaData) {
		for i, addr := range node.Host {
			store, err := remote.NewRemoteStore(addr)
			// 所有地址都无法连接
			if err != nil && i == len(node.Host)-1 {
				return fmt.Errorf("更新节点信息失败：无法连接到 %v", node.Id)
			}
			n.Store = store
		}
	}
	n.IsEntry = node.IsEntry
	n.IsExit = node.IsExit
	n.Cap = node.Cap
	n.Version = node.Version
	return nil
}

// Network 获取当前节点的网络状况，缓存一秒
func (n *Node) Network() *fs.Network {
	now := time.Now()
	if n.lastUpdate.Add(time.Second * 10).Before(now) {
		n.network = n.Store.Network()
	}
	n.lastUpdate = now
	return n.network
}
