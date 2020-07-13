package remote

import (
	store "github.com/shiningacg/filestore"
)

type API Store

func (A API) Get(uuid string) (file store.File, err error) {
	// 查询数据库，查看是否存在这样的文件
	// 获取到信息，创建file文件
	// 通过该id去寻找合适的节点获取url
	// 返回
	return
}

// remote不会调用这个方法
func (A API) Add(file store.File) error {
	return nil
}

//
func (A API) Remove(uuid string) error {
	// 在数据库中查找，获取信息
	// 标记删除
	// 下发请求
	return nil
}
