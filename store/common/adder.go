package common

import (
	store "github.com/shiningacg/filestore"
	ipfs "github.com/shiningacg/sn-ipfs"
	"net/http"
)

// Adder通过查库id来获取到文件，可以下载
type Adder interface {
	Find(uuid string) store.File
}

type IPFSFile struct {
	name string
	url  string
	ipfs.File
}

func (f *IPFSFile) Close() error {
	return nil
}

func (f *IPFSFile) FileName() string {
	return f.name
}

func (f *IPFSFile) ID() string {
	return f.Cid()
}

func (f *IPFSFile) Url() string {
	return f.url
}

// 通过http去下载文件
type IPFSAdder struct {
	ipfs.Store
}

func (a *IPFSAdder) Find(uuid string) store.File {
	node := a.Get(uuid)
	file, err := node.ToFile()
	if err != nil {
		return nil
	}
	return &IPFSFile{
		name: "",
		url:  "",
		File: file,
	}
}

type HttpFile struct {
	*http.Response
	name string
	url  string
	id   string
}

func (h HttpFile) Read(p []byte) (n int, err error) {
	return h.Body.Read(p)
}

func (h HttpFile) Seek(offset int64, whence int) (int64, error) {
	return 0, nil
}

func (h HttpFile) Close() error {
	panic("implement me")
}

func (h HttpFile) FileName() string {
	return h.name
}

func (h HttpFile) ID() string {
	return h.id
}

func (h HttpFile) Url() string {
	return h.url
}

func (h HttpFile) Size() uint64 {
	panic("implement me")
}

// 通过ipfs去下载文件
type HttpAdder struct{}

func (I *HttpAdder) Find(uuid string) store.File {
	gatewayAddr := ""
	// 通过主网关去查找文件
	rsp, err := http.Get(gatewayAddr)
	if err != nil {
		return nil
	}
	// 包装reader
	file := HttpFile{
		Response: rsp,
		name:     "",
		url:      "",
		id:       "",
	}
	return file
}
