package master

import (
	"context"
	"errors"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/shiningacg/filestore/cluster"
	"time"
)

func NewWatcher(client *clientv3.Client, path string) Watcher {
	return &watcher{
		path:   path,
		Client: client,
	}
}

type Watcher interface {
	Watch(ctx context.Context)
	Events(chan<- cluster.Event)
	UpdateAll()
}

type watcher struct {
	path string

	repo []chan<- cluster.Event
	*clientv3.Client
}

func (w *watcher) UpdateAll() {
	kv := clientv3.NewKV(w.Client)
	// 获取所有已经在连接的节点
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	resp, err := kv.Get(ctx, w.path, clientv3.WithPrefix())
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

func (w *watcher) Watch(ctx context.Context) {
	watcher := clientv3.NewWatcher(w.Client)

	wch := watcher.Watch(ctx, w.path, clientv3.WithPrefix())
	for wr := range wch {
		if wr.Canceled {
			return
		}
		for _, evt := range wr.Events {
			data := &cluster.Data{}
			switch evt.Type {
			case mvccpb.PUT:
				err := data.Decode(evt.Kv.Value)
				if err != nil {
					continue
				}
				w.sendEvent(cluster.NewEvent(data, cluster.PUT))
			case mvccpb.DELETE:
				id, err := w.idFromKey(string(evt.Kv.Key))
				if err != nil {
					continue
				}
				data.Id = id
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

func (w *watcher) idFromKey(key string) (string, error) {
	if len(key) < len(w.path) {
		return "", errors.New("无效的key")
	}
	return key[:len(w.path)], nil
}
