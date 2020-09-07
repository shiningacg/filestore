package master

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/shiningacg/filestore/cluster"
	"go.etcd.io/etcd/clientv3"
)

func NewWatcher(ctx context.Context, client *clientv3.Client, name string) Watcher {
	service := cluster.Service{Name: name}
	return &watcher{
		path:   service.ToPath(),
		ctx:    ctx,
		Client: client,
	}
}

type Watcher interface {
	Events(chan<- cluster.Event)
	UpdateAll()
}

type watcher struct {
	path string
	ctx  context.Context

	cancel func()
	wctx   context.Context

	repo []chan<- cluster.Event
	*clientv3.Client
}

func (w *watcher) UpdateAll() {
	kv := clientv3.NewKV(w.Client)
	// 获取所有已经在连接的节点
	resp, err := kv.Get(w.ctx, w.path, clientv3.WithPrefix())
	if err != nil {
		fmt.Println("初始化失败：", err)
		return
	}
	for _, v := range resp.Kvs {
		data := &cluster.Data{}
		id := string(v.Key[:len(w.path)])
		err := data.Decode(v.Value)
		if err != nil {
			fmt.Printf("无法加载节点信息：%v", id)
		}
		// 创建event
		w.sendEvent(cluster.NewEvent(data, cluster.PUT))
	}
}

func (w *watcher) watch() {
	watcher := clientv3.NewWatcher(w.Client)
	w.wctx, w.cancel = context.WithCancel(w.ctx)

	wch := watcher.Watch(w.wctx, w.path, clientv3.WithPrefix())
	for wr := range wch {
		if wr.Canceled {
			return
		}
		for _, evt := range wr.Events {
			data := &cluster.Data{}
			err := data.Decode(evt.Kv.Value)
			fmt.Println(err)
			switch evt.Type {
			case mvccpb.PUT:
				w.sendEvent(cluster.NewEvent(data, cluster.PUT))
			case mvccpb.DELETE:
				w.sendEvent(cluster.NewEvent(data, cluster.DEL))
			}
		}
	}
}

func (w *watcher) Events(repo chan<- cluster.Event) {
	w.repo = append(w.repo, repo)
}

func (w *watcher) sendEvent(event cluster.Event) {
	for _, c := range w.repo {
		c <- event
	}
}
