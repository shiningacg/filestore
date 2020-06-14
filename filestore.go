package filestore

type StoreType uint8

const (
	OS StoreType = iota
	IPFS
)

type Store interface {
	Stats() Stats
	API() API
}
