package follower

import (
	"context"
	"errors"
	"github.com/shiningacg/filestore/cluster"
	"go.etcd.io/etcd/clientv3"
	"log"
)

func NewFollower(client *clientv3.Client, data cluster.Data, service cluster.Service) *Follower {
	ctx, cf := context.WithCancel(context.Background())
	register := NewRegister(ctx, client, &data, service)
	watcher := cluster.NewWatcher(client, service.ToKey(), false)

	return &Follower{
		data:     data,
		service:  service,
		Register: register,
		Watcher:  watcher,
		cf:       cf,
		ctx:      ctx,
	}
}

type Follower struct {
	ctx context.Context
	cf  func()

	data    cluster.Data
	service cluster.Service
	Register
	cluster.Watcher
}

func (f *Follower) Run() error {
	// 检查连接是否正常，id是否冲突
	evt, err := f.Watcher.Exist()
	if err != nil {
		return err
	}
	if evt != nil && evt[0].Id != "" {
		return errors.New("节点ID重复")
	}
	// 注册节点信息
	err = f.Register.Register()
	if err != nil {
		return err
	}
	log.Println("节点注册成功！")
	ch := make(chan *cluster.Event, 10)
	f.Watcher.Events(ch)
	go f.Watcher.Watch(f.ctx)
	for {
		select {
		case <-f.ctx.Done():
			return nil
		case evt := <-ch:
			switch evt.Action {
			case cluster.DEL:
				f.close()
				return nil
			case cluster.PUT:
				if evt.Version > f.data.Version {
					f.ChangeData(evt.Data)
				}
			}
		}
	}
}

// TODO: 自动化的重新加载gateway
// 对改变对配置进行响应
func (f *Follower) ChangeData(data *cluster.Data) {
	log.Printf("更新配置：Version-%v", data.Version)
	// 判断那些内容发生了改变进行响应的重置
	// 最后更新数据
	f.data = *data
	log.Println("更新配置成功")
}

func (f *Follower) GetData() cluster.Data {
	return f.data
}

func (f *Follower) GetService() cluster.Service {
	return f.service
}

func (f *Follower) close() {
	log.Println("服务器下发关闭节点信息")
	f.Stop()
}

func (f *Follower) Stop() {
	log.Println("正在关闭节点...")
	f.cf()
}
