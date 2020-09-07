package master

import (
	"context"
	"fmt"
	"github.com/shiningacg/filestore/cluster"
	"go.etcd.io/etcd/clientv3"
	"log"
)

func NewMaster(ctx context.Context, client *clientv3.Client) *Master {
	nodes := make(Nodes, 0, 5)
	evt := make(chan cluster.Event, 5)
	watcher := NewWatcher(ctx, client, "svc.file")
	watcher.Events(evt)
	master := &Master{
		ctx:     ctx,
		Watcher: watcher,
		nodes:   nodes,
		evt:     evt,
	}
	go master.watch()
	return master
}

// 主要负责负载均衡
type Master struct {
	ctx   context.Context
	nodes Nodes
	Watcher
	evt chan cluster.Event
}

// Node 根据id查找一个节点，如果没有找到则返回nil
func (m *Master) Node(id string) *Node {
	return m.nodes.Node(id)
}

// 获取所有的节点
func (m *Master) Nodes() Nodes {
	return m.nodes
}

// Entries 获取所有入口节点
func (m *Master) Entries() Nodes {
	nodes := make(Nodes, 0, len(m.nodes))
	for _, v := range nodes {
		if v.IsEntry {
			nodes = append(nodes, v)
		}
	}
	return nodes
}

// Exits 获取所有出口节点
func (m *Master) Exits() Nodes {
	nodes := make(Nodes, 0, len(m.nodes))
	for _, v := range nodes {
		if v.IsExit {
			nodes = append(nodes, v)
		}
	}
	return nodes
}

// TODO：现在的Best判断都只通过网络状况判断，如果在一秒内发生了多个请求而缓存数据没有刷新，那么可能会被分配到同一个服务器中
// BestEntry 找到最佳的上传节点
func (m *Master) BestEntry() *Node {
	return m.Entries().SortBest(func(n1, n2 *Node) bool {
		return n1.Network().Upload < n2.Network().Upload
	})
}

// BestExit 找出最佳的出口节点
func (m *Master) BestExit() *Node {
	return m.Exits().SortBest(func(n1, n2 *Node) bool {
		return n1.Network().Download < n1.Network().Download
	})
}

// watch 监听etcd中发生的事件，对节点进行更新
func (m *Master) watch() {
	go m.Watcher.UpdateAll()
	for {
		select {
		case <-m.ctx.Done():
			return
		case evt := <-m.evt:
			switch evt.Action {
			case cluster.PUT:
				// 节点是否已经存在过了
				if node := m.nodes.Node(evt.Id); node != nil {
					// 版本号不同，那么更新node
					if node.Version != evt.Version {
						err := node.Update(evt.Data)
						// 更新信息失败，节点暂时不可用，进行删除
						// TODO：让node感知到错误的发生从而进行一次回滚？
						if err != nil {
							fmt.Println(err)
							m.nodes.Delete(node.Id)
						}
					}
				} else { // 节点是新加入的
					node, err := NewNodeFromData(evt.Data)
					if err != nil {
						log.Println(err)
					} else {
						m.nodes = append(m.nodes, node)
					}
				}
			case cluster.DEL:
				m.nodes.Delete(evt.Id)
			}
		}
	}
}
