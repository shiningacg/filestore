package os

import (
	"errors"
	fs "github.com/shiningacg/filestore"
	"io"
	"os"
)

var (
	ErrEmptyID      = errors.New("添加到os仓库时id不能为空")
	ErrFileNotFound = errors.New("无法找到文件")
)

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
		s.logger.Println(err)
		return err
	}
	n, err := io.Copy(f, file)
	if err != nil {
		err = errors.New("写入文件错误：" + err.Error())
		s.logger.Println(err)
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
		s.logger.Println(err)
	}
	return s.db.Delete(file.UUID)
}

func (s *Store) fromDBFile(file *DBFile) fs.ReadableFile {
	var bs = &fs.BaseFileStruct{}
	f, err := os.Open(file.Path)
	if err != nil {
		s.logger.Println(err)
		return nil
	}
	bs.SetUUID(file.UUID)
	bs.SetName(file.Name)
	bs.SetSize(file.Size)
	bs.SetUrl(s.gateway.GetUrl(file.UUID))
	return fs.NewReadableFile(bs, f)
}
