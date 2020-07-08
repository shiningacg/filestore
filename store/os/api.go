package os

import (
	"errors"
	store "github.com/shiningacg/filestore"
	"io"
	"log"
	"os"
)

type API struct {
	storeManager StoreManager
	baseGetUrl   string
	storePath    string
	logger       *log.Logger
	gatewayAddr  string
	db           *BoltDB
}

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
	dbfile := a.storeFileToDBFile(file)
	f, err := os.Create(dbfile.Path)
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
	dbfile.Size = uint64(n)
	return a.db.Add(dbfile)
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
	f.url = a.baseGetUrl + file.UUID
	return f
}
