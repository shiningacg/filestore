package ipfs

import (
	fs "github.com/shiningacg/filestore"
	ipfs "github.com/shiningacg/sn-ipfs"
	"log"
)

// TODO: 允许外部传入配置
func NewCore(logger *log.Logger) (*Store, error) {
	s, err := ipfs.NewStore("127.0.0.1:5001", "127.0.0.1:8080")
	if err != nil {
		return nil, err
	}
	st := &Store{
		ipfs: s,
	}
	return st, nil
}

type Store struct {
	ipfs ipfs.Store
	log  *log.Logger
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
