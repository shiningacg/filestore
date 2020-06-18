package filestore

import "io"

// 所有文件都要实现该接口
type File interface {
	io.ReadSeeker
	io.Closer
	FileName() string
	ID() string
	// 重定向跳转链接
	Url() string
}
