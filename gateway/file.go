package gateway

import "mime/multipart"

type PartFile struct {
	size uint64
	UUID string
	*multipart.Part
}

func (f *PartFile) Seek(offset int64, whence int) (int64, error) {
	panic("implement me")
}

func (f *PartFile) ID() string {
	return f.UUID
}

func (f *PartFile) Url() string {
	panic("implement me")
}

type File struct {
	size uint64
	UUID string
	multipart.File
	*multipart.FileHeader
}

func (f *File) FileName() string {
	return f.FileHeader.Filename
}

func (f *File) ID() string {
	return f.UUID
}

func (f *File) Url() string {
	panic("implement me")
}
