package mock

import (
	"bytes"
	fs "github.com/shiningacg/filestore"
	"github.com/shiningacg/filestore/store"
	"io/ioutil"
	"log"
)

func NewCore() store.Core {
	store := &FileStore{}
	return store
}

type FileStore struct{}

func (s *FileStore) Get(uuid string) (fs.ReadableFile, error) {
	var bs = &fs.BaseFileStruct{}
	log.Printf("从仓库取出文件：%v", uuid)
	data := []byte("测试数据")
	f := bytes.NewReader(data)
	bs.SetUUID(uuid)
	bs.SetSize(uint64(len(data)))
	bs.SetName("test.txt")
	return fs.NewReadableFile(bs, f), nil
}

func (s *FileStore) Add(file fs.ReadableFile) error {
	data, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	log.Printf("添加文件到仓库：%v %v %v", file.Name(), file.UUID(), file.Size())
	log.Println(string(data))
	file.SetUUID("test")
	return nil
}

func (s *FileStore) Remove(uuid string) error {
	log.Printf("从仓库删除文件：%v", uuid)
	return nil
}
