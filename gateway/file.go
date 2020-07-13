package gateway

import "mime/multipart"

type File struct {
	size uint64
	UUID string
	multipart.File
	*multipart.FileHeader
}

func (f *File) FileName() string {
	return f.Filename
}

func (f *File) ID() string {
	return f.UUID
}

func (f *File) Url() string {
	panic("implement me")
}
