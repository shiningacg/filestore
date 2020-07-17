package filestore

// 文件交互接口
type API interface {
	// 通过全局文件标识获取文件
	Get(uuid string) (File, error)
	// 向储存库添加文件,底层实现os和ipfs不一样。
	Add(file File) error
	// 通过全局文件标识删除储存库内的文件
	Remove(uuid string) error
}
