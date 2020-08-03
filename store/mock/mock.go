package mock

import (
	"bytes"
	store "github.com/shiningacg/filestore"
	"github.com/shiningacg/filestore/gateway"
	"io/ioutil"
	"log"
)

func NewStore(g *gateway.Gateway) *Store {
	go func() {
		err := g.Run()
		if err != nil {
			panic(err)
		}
	}()
	return &Store{g: g}
}

type Store struct {
	g *gateway.Gateway
}

func (s *Store) Get(uuid string) (store.ReadableFile, error) {
	var bs = &store.BaseFileStruct{}
	log.Printf("从仓库取出文件：%v", uuid)
	data := []byte("测试数据")
	f := bytes.NewReader(data)
	bs.SetUUID(uuid)
	bs.SetUrl(s.g.GetUrl(bs.UUID()))
	bs.SetSize(uint64(len(data)))
	bs.SetName("test.txt")
	return store.NewReadableFile(bs, f), nil
}

func (s *Store) Add(file store.ReadableFile) error {
	data, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	log.Printf("添加文件到仓库：%v %v %v %v", file.Name(), file.UUID(), file.Size(), file.Url())
	log.Println(string(data))
	file.SetUUID("test")
	file.SetUrl(s.g.GetUrl(file.UUID()))
	return nil
}

func (s *Store) Remove(uuid string) error {
	log.Printf("从仓库删除文件：%v", uuid)
	return nil
}

func (s *Store) Space() *store.Space {
	return &store.Space{
		Cap:   111,
		Total: 222,
		Free:  200,
		Used:  22,
	}
}

func (s *Store) Network() *store.Network {
	return &store.Network{
		Upload:   1000,
		Download: 2000,
	}
}

func (s *Store) Gateway() *store.Bandwidth {
	return s.g.BandWidth()
}
