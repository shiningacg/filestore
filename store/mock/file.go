package mock

import (
	"bytes"
	store "github.com/shiningacg/filestore"
	"github.com/shiningacg/filestore/gateway"
	"io/ioutil"
	"log"
)

func NewFileStore(g *gateway.Gateway) store.FileStore {
	store := &FileStore{g: g}
	g.SetStore(store)
	go func() {
		err := g.Run()
		if err != nil {
			panic(err)
		}
	}()
	return store
}

func NewFileStoreWithoutWeb() store.FileStore {
	return &FileStore{}
}

type FileStore struct {
	g *gateway.Gateway
}

func (s *FileStore) Get(uuid string) (store.ReadableFile, error) {
	var bs = &store.BaseFileStruct{}
	log.Printf("从仓库取出文件：%v", uuid)
	data := []byte("测试数据")
	f := bytes.NewReader(data)
	bs.SetUUID(uuid)
	bs.SetSize(uint64(len(data)))
	bs.SetName("test.txt")
	return store.NewReadableFile(bs, f), nil
}

func (s *FileStore) Add(file store.ReadableFile) error {
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

func (s *FileStore) Space() *store.Space {
	return &store.Space{
		Cap:   111,
		Total: 222,
		Free:  200,
		Used:  22,
	}
}

func (s *FileStore) Network() *store.Network {
	return &store.Network{
		Upload:   1000,
		Download: 2000,
	}
}

func (s *FileStore) Gateway() *store.Bandwidth {
	if s.g != nil {
		return s.g.BandWidth()
	}
	return &store.Bandwidth{
		Visit:         100,
		DayVisit:      10,
		HourVisit:     1,
		Bandwidth:     1000,
		DayBandwidth:  100,
		HourBandwidth: 10,
	}
}
