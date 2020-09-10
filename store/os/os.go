package os

import (
	"errors"
	fs "github.com/shiningacg/filestore"
	"github.com/shiningacg/filestore/gateway"
	"github.com/shiningacg/filestore/store/common"
	"github.com/shiningacg/filestore/store/remote"
	"github.com/shiningacg/mygin-frame-libs/log"
	"io"
	"os"
)

var (
	ErrEmptyID      = errors.New("添加到os仓库时id不能为空")
	ErrFileNotFound = errors.New("无法找到文件")
)

func NewOStore(config *StoreConfig, checker gateway.Checker, logger *log.Logger) *Store {
	g := gateway.NewGateway(config.GatewayAddr, checker, logger)
	s := &Store{
		gateway:      g,
		storeManager: NewDefaultManager(config.StorePath),
		logger:       logger,
		db:           OpenBoltDB(config.StorePath+"/store.dat", logger),
	}
	g.SetStore(s)
	go g.Run()
	return s
}

type StoreConfig struct {
	GatewayAddr string
	StorePath   string
}

type Store struct {
	gateway      *gateway.Gateway
	storeManager StoreManager
	logger       *log.Logger
	db           *BoltDB
	remote.Adder
}

func (s *Store) Get(uuid string) (fs.ReadableFile, error) {
	dbFile := s.db.Get(uuid)
	if dbFile == nil {
		return nil, errors.New("没有找到文件：" + uuid)
	}
	file := s.fromDBFile(dbFile)
	if file == nil {
		return nil, errors.New("文件丢失")
	}
	return file, nil
}

// 不嫩使用这里的file的size方法
func (s *Store) Add(file fs.ReadableFile) error {
	// 添加到os到文件一定要有id，没有则报错
	if file.UUID() == "" {
		return ErrEmptyID
	}
	// 测试是否可读,如果不可读，则调用adder去创建一个可读到reader
	if false {
		f := s.Find(file)
		if f == nil {
			return ErrFileNotFound
		}
		file = f
	}
	dbFile := s.storeFileToDBFile(file)
	f, err := os.Create(dbFile.Path)
	if err != nil {
		err = errors.New("无法创建文件：" + err.Error())
		s.logger.Fatal(err)
		return err
	}
	n, err := io.Copy(f, file)
	if err != nil {
		err = errors.New("写入文件错误：" + err.Error())
		s.logger.Fatal(err)
	}
	dbFile.Size = uint64(n)
	return s.db.Add(dbFile)
}

func (s *Store) storeFileToDBFile(file fs.ReadableFile) *DBFile {
	dbFile := &DBFile{
		UUID: file.UUID(),
		Name: file.Name(),
	}
	dbFile.Path = s.storeManager.GetStorePath(file)
	return dbFile
}

func (s *Store) Remove(uuid string) error {
	file := s.db.Get(uuid)
	if file == nil {
		return nil
	}
	err := os.Remove(file.Path)
	if err != nil {
		err = errors.New("删除文件错误：" + err.Error())
		s.logger.Fatal(err)
	}
	return s.db.Delete(file.UUID)
}

func (s *Store) fromDBFile(file *DBFile) fs.ReadableFile {
	var bs = &fs.BaseFileStruct{}
	f, err := os.Open(file.Path)
	if err != nil {
		s.logger.Fatal(err)
		return nil
	}
	bs.SetUUID(file.UUID)
	bs.SetName(file.Name)
	bs.SetSize(file.Size)
	return fs.NewReadableFile(bs, f)
}

func (s *Store) Space() *fs.Space {
	space := &fs.Space{}
	dbInfo := s.db.Info()
	if dbInfo != nil {
		space.Total = dbInfo.MaxSize
		space.Used = dbInfo.UsedSize
		space.Free = dbInfo.FreeSize
	}
	diskStats := common.DiskUsage(s.storeManager.GetBasePath())
	if diskStats != nil {
		space.Cap = diskStats.Total - diskStats.Used
	}
	return space
}

func (s *Store) Network() *fs.Network {
	panic("implement me")
}

func (s *Store) Gateway() *fs.Bandwidth {
	return s.gateway.BandWidth()
}
