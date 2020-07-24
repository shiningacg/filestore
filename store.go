package filestore

// 进行信息的交互，无法直接进行读取
type InfoStore interface {
	InfoFS
	Stats
}

// 可以进行读取
type FileStore interface {
	FileFS
	Stats
}
