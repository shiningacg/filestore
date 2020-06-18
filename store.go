package filestore

type StoreType uint8

const (
	// 通过系统文件系统实现
	OS StoreType = iota
	// 通过ipfs文件系统实现
	IPFS
)

type Store interface {
	Stats() Stats
	API() API
}
