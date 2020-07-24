package gateway

import (
	fs "github.com/shiningacg/filestore"
	"mime/multipart"
	"os"
)

func NewMultipartFile(file multipart.File, header *multipart.FileHeader) fs.ReadableFile {
	var bs = &fs.BaseFileStruct{}
	bs.SetName(header.Filename)
	return fs.NewReadableFile(bs, file)
}

func NewOSFile(fileStruct *fs.BaseFileStruct, file *os.File) fs.ReadableFile {
	return fs.NewReadableFile(fileStruct, file)
}
