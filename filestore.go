package filestore

import "io"

// 此包为业务和底层实现的缓冲层
type FileStore interface {
	// 存贮文件
	Store(writer io.Reader) (File, error)
	// 删除文件
	Delete(uuid string) error
	// 获取文件
	Get(uuid string) File
}
