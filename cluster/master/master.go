package master

import (
	"context"
	"github.com/coreos/etcd/clientv3"
	"github.com/shiningacg/filestore/cluster"
	"log"
)

func NewMaster(ctx context.Context, client *clientv3.Client, service cluster.Service) *Master {
	nodes := make(cluster.Nodes, 0, 5)
	evt := make(chan cluster.Event, 5)
	watcher := NewWatcher(client, service.ToPath())
	watcher.Events(evt)
	master := &Master{
		ctx:     ctx,
		Watcher: watcher,
		nodes:   nodes,
		evt:     evt,
	}
	go master.Watcher.Watch(ctx)
	go master.watch()
	return master
}

// 主要负责负载均衡
type Master struct {
	cluster.Service
	ctx   context.Context
	nodes cluster.Nodes
	Watcher
	evt chan cluster.Event
}

// Node 根据id查找一个节点，如果没有找到则返回nil
func (m *Master) Node(id string) cluster.Node {
	return m.nodes.Node(id)
}

// 获取所有的节点
func (m *Master) Nodes() cluster.Nodes {
	return m.nodes
}

// Entries 获取所有入口节点
func (m *Master) Entries() cluster.Nodes {
	nodes := make(cluster.Nodes, 0, len(m.nodes))
	for _, v := range nodes {
		if v.Data().IsEntry() {
			nodes = append(nodes, v)
		}
	}
	return nodes
}

// Exits 获取所有出口节点
func (m *Master) Exits() cluster.Nodes {
	nodes := make(cluster.Nodes, 0, len(m.nodes))
	for _, v := range nodes {
		if v.Data().IsExit() {
			nodes = append(nodes, v)
		}
	}
	return nodes
}

// TODO：现在的Best判断都只通过网络状况判断，如果在一秒内发生了多个请求而缓存数据没有刷新，那么可能会被分配到同一个服务器中
// BestEntry 找到最佳的上传节点
func (m *Master) BestEntry(size uint64) cluster.Node {
	nodes := m.Entries()
	if size == 0 {
		return nodes.BestUpload()
	}
	// TODO: 多线程查询，使用ctx限制超时时间
	for {
		if len(nodes) == 0 {
			break
		}
		node := nodes.BestDownload()
		if node.Space().Free < size {
			nodes = nodes.Delete(node.ID())
			continue
		}
		return node
	}
	return nil
}

// BestExit 找出最佳的出口节点
func (m *Master) BestExit(fid string) cluster.Node {
	nodes := m.Entries()
	if fid == "" {
		return nodes.BestUpload()
	}
	// TODO: 多线程查询，使用ctx限制超时时间
	for {
		if len(nodes) == 0 {
			break
		}
		node := nodes.BestDownload()
		_, err := node.Get(fid)
		if err != nil {
			nodes = nodes.Delete(node.ID())
			continue
		}
		return node
	}
	return nil
}

func (m *Master) Watch(repo chan<- cluster.Event) {
	m.Watcher.Events(repo)
}

// watch 监听etcd中发生的事件，对节点进行更新
func (m *Master) watch() {
	go func() {
		err := m.Watcher.UpdateAll()
		if err != nil {
			log.Println(err)
		}
	}()
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
					data := node.Data()
					if data.Version != evt.Version {
						err := node.Update(evt.Data)
						// 更新信息失败，节点暂时不可用，进行删除
						// TODO：让node感知到错误的发生从而进行一次回滚？
						if err != nil {
							log.Println(err)
							m.nodes.Delete(node.ID())
						}
					}
				} else { // 节点是新加入的
					node, err := NewNodeFromData(evt.Data)
					if err != nil {
						log.Println(err)
					} else {
						m.nodes = append(m.nodes, node)
						log.Printf("新节点上线：%v", evt.Id)
					}
				}
			case cluster.DEL:
				if m.nodes.Node(evt.Id) != nil {
					log.Println("节点离线:", evt.Id)
				}
				m.nodes.Delete(evt.Id)
			}
		}
	}
}
