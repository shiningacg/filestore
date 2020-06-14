package filestore

import "io"

type File interface {
	io.ReadSeeker
	io.Closer
	FileName() string
	ID() string
	// 重定向跳转链接
	Url() string
}
