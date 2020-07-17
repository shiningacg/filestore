package os

import (
	"errors"
	store "github.com/shiningacg/filestore"
	"io"
	"os"
)

var (
	ErrEmptyID      = errors.New("添加到os仓库时id不能为空")
	ErrFileNotFound = errors.New("无法找到文件")
)

type API Store

func (a *API) Get(uuid string) (store.File, error) {
	dbFile := a.db.Get(uuid)
	if dbFile == nil {
		return nil, errors.New("没有找到文件：" + uuid)
	}
	file := a.fromDBFile(dbFile)
	if file == nil {
		return file, errors.New("文件丢失")
	}
	return file, nil
}

// 不嫩使用这里的file的size方法
func (a *API) Add(file store.File) error {
	// 添加到os到文件一定要有id，没有则报错
	if file.ID() == "" {
		return ErrEmptyID
	}
	// 测试是否可读,如果不可读，则调用adder去创建一个可读到reader
	if false {
		f := a.Find(file)
		if f == nil {
			return ErrFileNotFound
		}
		file = f
	}
	dbFile := a.storeFileToDBFile(file)
	f, err := os.Create(dbFile.Path)
	if err != nil {
		err = errors.New("无法创建文件：" + err.Error())
		a.logger.Println(err)
		return err
	}
	n, err := io.Copy(f, file)
	if err != nil {
		err = errors.New("写入文件错误：" + err.Error())
		a.logger.Println(err)
	}
	dbFile.Size = uint64(n)
	return a.db.Add(dbFile)
}

func (a *API) storeFileToDBFile(file store.File) *DBFile {
	dbFile := &DBFile{
		UUID: file.ID(),
		Name: file.FileName(),
	}
	dbFile.Path = a.storeManager.GetStorePath(file)
	return dbFile
}

func (a *API) Remove(uuid string) error {
	file := a.db.Get(uuid)
	err := os.Remove(file.Path)
	if err != nil {
		err = errors.New("删除文件错误：" + err.Error())
		a.logger.Println(err)
	}
	return a.db.Delete(file.UUID)
}

func (a *API) fromDBFile(file *DBFile) *File {
	f, err := fromDBFile(file)
	if err != nil {
		a.logger.Println(err)
		return nil
	}
	f.url = a.gateway.GetUrl(file.UUID)
	return f
}
