package filestore

import (
	"errors"
	"io"
)

// 应该被处理的错误
var (
	ErrSeekNotSupport = errors.New("不支持seek方法")
)

// 只包含信息的文件
type BaseFile interface {
	// 全局文件id
	UUID() string
	SetUUID(uuid string)
	// 文件名
	Name() string
	SetName(name string)
	// 文件大小
	Size() uint64
	SetSize(size uint64)
}

// 可以直接读取的文件
type ReadableFile interface {
	BaseFile
	io.ReadSeeker
	io.Closer
}

/*
	接口的简单实现：
	  包含 BaseFileStruct readableFile
*/

// BaseFile接口的简单实现
type BaseFileStruct struct {
	uuid string
	name string
	size uint64
}

func (b *BaseFileStruct) SetSize(size uint64) {
	b.size = size
}

func (b *BaseFileStruct) SetUUID(uuid string) {
	b.uuid = uuid
}

func (b *BaseFileStruct) SetName(name string) {
	b.name = name
}

func (b *BaseFileStruct) Name() string {
	return b.name
}

func (b *BaseFileStruct) Size() uint64 {
	return b.size
}

func (b *BaseFileStruct) UUID() string {
	return b.uuid
}

func NewReadableFile(base BaseFile, reader io.Reader) *readableFile {
	return &readableFile{
		Reader:   reader,
		BaseFile: base,
	}
}

type readableFile struct {
	io.Reader
	BaseFile
}

func (r *readableFile) Seek(offset int64, whence int) (int64, error) {
	if seeker, ok := r.Reader.(io.Seeker); ok {
		return seeker.Seek(offset, whence)
	}
	return 0, ErrSeekNotSupport
}

func (r *readableFile) Close() error {
	if closer, ok := r.Reader.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}
