package master

import (
	"fmt"
	fs "github.com/shiningacg/filestore"
	"github.com/shiningacg/filestore/cluster"
	"github.com/shiningacg/filestore/store/remote"
	"time"
)

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
	node.data.MetaData.Id = data.Id
	return node, node.Update(data)
}

// Node 实际操作的节点对象
type Node struct {
	*remote.Store
	data cluster.Data
	// 缓存信息
	network    *fs.Network
	lastUpdate time.Time
}

func (n *Node) ID() string {
	return n.data.Id
}

// Update 更新节点的信息，如果地址发送了改变那么会重新建立grpc连接
func (n *Node) Update(node *cluster.Data) error {
	// 如果监听地址变化了，那么就需要重新加载
	data := n.data
	if data.IsHostChange(node.MetaData) {
		for i, addr := range node.Host {
			fmt.Println(addr)
			store, err := remote.NewRemoteStore(addr)
			// 所有地址都无法连接
			if err != nil && i == len(node.Host)-1 {
				return fmt.Errorf("更新节点信息失败：无法连接到 %v", node.Id)
			}
			n.Store = store
		}
	}
	data.Exit = node.Entry
	data.Exit = node.Exit
	data.Cap = node.Cap
	data.MetaData.Update(node.MetaData)
	return nil
}

func (n *Node) Data() cluster.Data {
	return n.data
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

func (n *Node) Entry(token string) string {
	return n.data.GatewayAddr + "/upload/" + token
}

func (n *Node) Exit(fid string) string {
	return n.data.GatewayAddr + "/download/" + fid
}
