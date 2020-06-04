package filestore

import "io"

type File interface {
	// 提供直接读写的方法
	io.Reader
	// 获取下载链接
	GetUrl() string
	// 获取文件名
	GetName() string
	// 获取全局文件id
	GetUUID() string
}
