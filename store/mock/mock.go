package mock

import (
	"fmt"
	store "github.com/shiningacg/filestore"
	"github.com/shiningacg/filestore/gateway"
	"io/ioutil"
	"os"
)

type Store struct {
	g *gateway.Gateway
}

func (s *Store) Get(uuid string) (store.ReadableFile, error) {
	var bs = &store.BaseFileStruct{}
	f, _ := os.Open("mock.txt")
	stats, _ := f.Stat()
	bs.SetUUID("aa")
	bs.SetUrl(s.g.GetUrl(bs.UUID()))
	bs.SetSize(uint64(stats.Size()))
	bs.SetName(stats.Name())
	return store.NewReadableFile(bs, f), nil
}

func (s *Store) Add(file store.ReadableFile) error {
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bytes))
	return nil
}

func (s *Store) Remove(uuid string) error {
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
	return &store.Bandwidth{
		Visit:         3,
		DayVisit:      2,
		HourVisit:     1,
		Bandwidth:     1000,
		DayBandwidth:  200,
		HourBandwidth: 100,
	}
}
