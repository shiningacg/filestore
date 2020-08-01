package os

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/boltdb/bolt"
	fs "github.com/shiningacg/filestore"
	"github.com/shiningacg/mygin-frame-libs/log"
	"time"
)

var (
	ErrDataNotFound   = errors.New("无法找到数据")
	ErrBucketNotFound = errors.New("无法找到指定bucket")
	ErrDBInfoNotFound = errors.New("无法找到储存信息")
	ErrDumpInfo       = errors.New("重复储存文件信息")
)

type DBInfo struct {
	// 总共存放的文件数量
	InitTime uint64
	FileNum  uint64
	UsedSize uint64
	FreeSize uint64
	MaxSize  uint64
}

func (info *DBInfo) Json() []byte {
	b, _ := json.Marshal(info)
	return b
}

func (info *DBInfo) FromJson(data []byte) error {
	return json.Unmarshal(data, info)
}

func (info *DBInfo) AddFile(file *DBFile) error {
	if info.FreeSize < file.Size {
		return errors.New("空间不足")
	}
	info.FreeSize -= file.Size
	info.UsedSize += file.Size
	return nil
}

func (info *DBInfo) DeleteFIle(file *DBFile) {
	info.UsedSize -= file.Size
	info.FreeSize += file.Size
}

// 这里是存放文件信息的结构体
type DBFile struct {
	Name    string
	Size    uint64
	Path    string
	UUID    string
	Deleted bool
}

func (f *DBFile) Json() []byte {
	b, _ := json.Marshal(f)
	return b
}

func (f *DBFile) FromJson(b []byte) error {
	return json.Unmarshal(b, f)
}

func (f *DBFile) FromStoreFile(file fs.BaseFile) {
	f.Name = file.Name()
	f.UUID = file.UUID()
}

var (
	// 默认bucket名称
	DefaultBucket = []byte("DefaultBucket")
)

// 尝试打开数据库
func OpenBoltDB(path string, logger *log.Logger) *BoltDB {
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		logger.Fatal(err)
		panic(err)
	}
	DB := &BoltDB{
		log: logger,
		db:  db,
	}
	// 尝试初始化数据库
	err = DB.init(1024 * 1024 * 1024)
	if err != nil {
		logger.Fatal(err)
		panic(err)
	}
	return DB
}

type BoltDB struct {
	log *log.Logger
	db  *bolt.DB
}

func (b *BoltDB) Add(file *DBFile) error {
	if file := b.Get(file.UUID); file != nil && !file.Deleted {
		fmt.Println(file)
		return ErrDumpInfo
	}
	info := b.Info()
	if info == nil {
		return ErrDBInfoNotFound
	}
	err := info.AddFile(file)
	if err != nil {
		return err
	}
	return b.set(file)
}

func (b *BoltDB) Get(uuid string) *DBFile {
	var file = &DBFile{}
	err := b.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(DefaultBucket)
		if bucket == nil {
			panic(ErrBucketNotFound)
		}
		data := bucket.Get([]byte(uuid))
		if data == nil {
			file = nil
			return ErrDataNotFound
		}
		return file.FromJson(data)
	})
	if err != nil && err != ErrDataNotFound {
		b.log.Fatal(err)
		return nil
	}
	if file == nil {
		return nil
	}
	if file.Deleted {
		return nil
	}
	return file
}

func (b *BoltDB) Update(file *DBFile) error {
	if file.UUID == "" {
		return errors.New("uuid不能为空")
	}
	f := b.Get(file.UUID)
	// 只允许修改name和path
	if file.Name != "" {
		f.Name = file.Name
	}
	if file.Path != "" {
		f.Path = file.Path
	}
	return b.set(f)
}

// 并不会真的删除文件，只是标记删除而已
func (b *BoltDB) Delete(uuid string) error {
	file := b.Get(uuid)
	if file == nil {
		return nil
	}
	file.Deleted = true
	info := b.Info()
	info.DeleteFIle(file)
	return b.set(file)
}

func (b *BoltDB) Close() error {
	return b.db.Close()
}

// 初始化信息
func (b *BoltDB) init(maxSize uint64) error {
	info := &DBInfo{}
	info.MaxSize = maxSize
	// 初始化bucket
	err := b.db.Update(func(tx *bolt.Tx) error {
		var err error
		bucket := tx.Bucket(DefaultBucket)
		if bucket == nil {
			bucket, err = tx.CreateBucket(DefaultBucket)
			if err != nil {
				return errors.New("无法初始化数据库：" + err.Error())
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	info.InitTime = uint64(time.Now().Unix())
	info.FreeSize = info.MaxSize
	return b.setInfo(info)
}

// 对文件信息的操作
func (b *BoltDB) Info() *DBInfo {
	var info = &DBInfo{}
	err := b.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(DefaultBucket)
		if bucket == nil {
			panic("无法找到指定bucket")
		}
		data := bucket.Get([]byte("info"))
		err := info.FromJson(data)
		if err != nil {
			err = errors.New("获取存储信息失败：" + err.Error())
			return err
		}
		return nil
	})
	if err != nil {
		b.log.Fatal(err)
		return nil
	}
	return info
}

func (b *BoltDB) setInfo(info *DBInfo) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(DefaultBucket)
		if bucket == nil {
			err := errors.New("无法找到指定bucket")
			b.log.Fatal(err.Error())
			panic(err)
		}
		return bucket.Put([]byte("info"), info.Json())
	})
}

// 覆盖写入
func (b *BoltDB) set(file *DBFile) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(DefaultBucket)
		if bucket == nil {
			err := errors.New("无法找到指定bucket")
			b.log.Fatal(err.Error())
			panic(err)
		}
		return bucket.Put([]byte(file.UUID), file.Json())
	})
}
