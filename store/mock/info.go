package mock

import (
	"bytes"
	store "github.com/shiningacg/filestore"
	l "log"
)

func NewInfoStore() store.InfoStore {
	return &InfoStore{}
}

type InfoStore struct {
}

func (s *InfoStore) Get(uuid string) (store.BaseFile, error) {
	var bs = &store.BaseFileStruct{}
	l.Printf("从仓库取出文件：%v", uuid)
	data := []byte("测试数据")
	f := bytes.NewReader(data)
	bs.SetUUID(uuid)
	bs.SetSize(uint64(len(data)))
	bs.SetName("test.txt")
	return store.NewReadableFile(bs, f), nil
}

func (s *InfoStore) Add(file store.BaseFile) error {
	l.Printf("添加文件到仓库：%v %v %v", file.Name(), file.UUID(), file.Size())
	file.SetUUID("test")
	return nil
}

func (s *InfoStore) Remove(uuid string) error {
	l.Printf("从仓库删除文件：%v", uuid)
	return nil
}

func (s *InfoStore) Space() *store.Space {
	return &store.Space{
		Cap:   111,
		Total: 222,
		Free:  200,
		Used:  22,
	}
}

func (s *InfoStore) Network() *store.Network {
	return &store.Network{
		Upload:   1000,
		Download: 2000,
	}
}

func (s *InfoStore) Gateway() *store.Bandwidth {
	return &store.Bandwidth{
		Visit:         100,
		DayVisit:      10,
		HourVisit:     1,
		Bandwidth:     1000,
		DayBandwidth:  100,
		HourBandwidth: 10,
	}
}
