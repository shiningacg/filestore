package cluster

type Master interface {
	// Node 根据id查找一个节点，如果没有找到则返回nil
	Node(id string) Node
	// 获取所有的节点
	Nodes() Nodes
	// Entries 获取所有入口节点
	Entries() Nodes
	// Exits 获取所有出口节点
	Exits() Nodes
	// BestEntry 找到最佳的上传节点
	BestEntry(size uint64) Node
	// BestExit 找出最佳的出口节点
	BestExit(fid string) Node
	// 提供监控的功能
	Watch(chan<- Event)
}
