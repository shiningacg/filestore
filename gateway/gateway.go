package gateway

import (
	"context"
	"fmt"
	store "github.com/shiningacg/filestore"
	"log"
	"net/http"
)

func NewGateway(addr string, api store.API, logger *log.Logger) *Gateway {
	return &Gateway{
		log:     logger,
		addr:    addr,
		api:     api,
		monitor: NewMonitor(context.TODO()),
	}
}

type Gateway struct {
	// 负责日志控制
	log *log.Logger
	// 负责数据统计
	monitor *DefaultMonitor
	addr    string
	// 存放文件的仓库，能够通过id存放和获取文件
	api store.API
}

/*
	可供外部调用的方法
*/

// 获取统计信息
func (g *Gateway) BandWidth() *store.Gateway {
	return g.monitor.Bandwidth()
}

func (g *Gateway) Run() error {

	go g.monitor.Run()
	return http.ListenAndServe(g.addr, (*HttpServer)(g))
}

// 传入一个uuid，返回下载地址
func (g *Gateway) GetUrl(uuid string) string {
	return fmt.Sprintf("http://%v/get/%v", g.addr, uuid)
}

// 传入一个uuid，获取上传地址
func (g *Gateway) PostUrl(uuid string) string {
	return fmt.Sprintf("http://%v/post/%v", g.addr, uuid)
}
