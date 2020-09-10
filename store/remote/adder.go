package remote

import (
	"bytes"
	"fmt"
	fs "github.com/shiningacg/filestore"
	ipfs "github.com/shiningacg/sn-ipfs"
	"log"
	"net/http"
)

// Adder通过查库id来获取到文件，可以下载
type Adder interface {
	Find(file fs.BaseFile) fs.ReadableFile
}

// 通过ipfs去下载文件
type IPFSAdder struct {
	ipfs.Store
}

func (a *IPFSAdder) Find(file fs.BaseFile) fs.ReadableFile {
	node := a.Get(file.UUID())
	f, err := node.ToFile()
	f = f
	if err != nil {
		return nil
	}
	return nil
}

func NewHttpAdder(gatewayAddr string) *HttpAdder {
	return &HttpAdder{gatewayAddr: gatewayAddr}
}

// 通过http去下载文件
type HttpAdder struct {
	gatewayAddr string
}

func (a *HttpAdder) Find(file fs.BaseFile) fs.ReadableFile {
	var bs = &fs.BaseFileStruct{}
	gatewayAddr := a.getUrl(file.UUID())
	// 通过主网关去查找文件
	rsp, err := http.Get(gatewayAddr)
	if err != nil {
		return nil
	}
	// 包装reader
	bs.SetName(file.Name())
	bs.SetUUID(file.UUID())
	return fs.NewReadableFile(bs, rsp.Body)
}

func (a *HttpAdder) getUrl(uuid string) string {
	return fmt.Sprintf("http://%v/get/%v", a.gatewayAddr, uuid)
}

type MockAdder struct{}

func (m MockAdder) Find(file fs.BaseFile) fs.ReadableFile {
	log.Printf("通过adder寻找文件：%v", file.UUID())
	var bt = []byte("测试数据")
	return fs.NewReadableFile(file, bytes.NewReader(bt))
}
