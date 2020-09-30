package store

const (
	IPFS = "ipfs"
	OS   = "os"
)

type Config struct {
	// 启动的仓库类型
	Type string
	// 文件存放位置
	Path string
	// 网关地址
	Gateway string
	// 插件
	Plugin
}

type Plugin struct {
	// 如果提供，os仓库可以进行远程添加文件
	AdderAddr string
	// 如果提供，用户上传文件时需要进行身份核对
	CheckerAddr string
}
