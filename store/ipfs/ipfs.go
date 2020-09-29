package ipfs

import (
	fs "github.com/shiningacg/filestore"
	"github.com/shiningacg/filestore/gateway"
	"github.com/shiningacg/filestore/gateway/checker"
	"github.com/shiningacg/filestore/store"
	"github.com/shiningacg/mygin-frame-libs/log"
	ipfs "github.com/shiningacg/sn-ipfs"
)

// TODO: 允许外部传入配置
func NewStore(gatewayAddr string, checker checker.Checker, logger *log.Logger) (fs.FileStore, error) {
	s, err := ipfs.NewStore("127.0.0.1:5001", "127.0.0.1:8080")
	if err != nil {
		return nil, err
	}
	st := &Store{
		ipfs: s,
	}
	return store.NewStore(st, gatewayAddr, checker, logger), nil
}

type Store struct {
	ipfs ipfs.Store
	g    gateway.Gateway
	log  *log.Logger
}

// 工厂方法
func (s *Store) SetGateway(g gateway.Gateway) {
	s.g = g
}

func (s *Store) GetGateway() gateway.Gateway {
	return s.g
}

func (s *Store) GetLogger() *log.Logger {
	return s.log
}

func (s *Store) SetLogger(l *log.Logger) {
	s.log = l
}

func (s *Store) Add(file fs.ReadableFile) error {
	f, err := s.ipfs.AddFromReader(file)
	if err != nil {
		return err
	}
	err = s.ipfs.PinMany(f.Blocks())
	if err != nil {
		return err
	}
	file.SetUUID(f.Cid())
	return nil
}

func (s *Store) Get(uuid string) (fs.ReadableFile, error) {
	n := s.ipfs.Get(uuid)
	f, err := n.ToFile()
	if err != nil {
		return nil, err
	}
	var bf = &fs.BaseFileStruct{}
	bf.SetSize(f.Size())
	bf.SetUUID(f.Cid())
	bf.SetName(f.Cid())
	return fs.NewReadableFile(bf, f), nil
}

func (s *Store) Remove(uuid string) error {
	n := s.ipfs.Get(uuid)
	f, err := n.ToFile()
	if err != nil {
		return err
	}
	return s.ipfs.UnpinMany(f.Blocks())
}

func (s *Store) Space() *fs.Space {
	return nil
}

func (s *Store) Network() *fs.Network {
	return nil
}

func (s *Store) Gateway() *fs.Bandwidth {
	return s.g.BandWidth()
}
