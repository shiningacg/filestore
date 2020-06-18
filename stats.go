package filestore

type Stats interface {
	// 查询空间信息
	Space() Space
	// 查询网络状况
	Bandwidth() Bandwidth
	// 查询http网关信息
	Gateway() Gateway
}

// 储存空间情况数据结构
type Space struct {
	Total uint64
	Free  uint64
	Used  uint64
}

// 网络状况的数据结构
type Bandwidth struct {
	Upload        uint64
	Download      uint64
	TotalUpload   uint64
	TotalDownload uint64
}

type Gateway struct {
	// 总共访问量
	Visit uint64
	// 日访问量
	DayVisit uint64
	// 总共流出流量
	Bandwidth uint64
	// 日流出流量
	DayBandwidth uint64
}
