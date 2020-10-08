package cluster

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
	"log"
	"strings"
	"time"
)

// watcher负责监听节点配置中发生的事件
type Watcher interface {
	Watch(ctx context.Context)
	// 注册监听管道
	Events(chan<- *Event)
	Exist() ([]*Event, error)
}

func NewWatcher(client *clientv3.Client, path string, prefix bool) *watcher {
	return &watcher{
		client: client,
		path:   path,
		prefix: prefix,
	}
}

type watcher struct {
	client *clientv3.Client
	ctx    context.Context

	path   string
	prefix bool

	repo []chan<- *Event
}

func (w *watcher) Watch(ctx context.Context) {
	if w.ctx != nil {
		return
	}
	watcher := clientv3.NewWatcher(w.client)
	w.ctx = ctx

	var wch clientv3.WatchChan
	switch w.prefix {
	case true:
		wch = watcher.Watch(ctx, w.path, clientv3.WithPrefix())
	case false:
		wch = watcher.Watch(ctx, w.path)
	}

	for {
		select {
		case <-w.ctx.Done():
			return
		case wr := <-wch:
			if wr.Canceled {
				return
			}
			for _, evt := range wr.Events {
				data := &Data{}
				switch evt.Type {
				case mvccpb.PUT:
					err := data.Decode(evt.Kv.Value)
					if err != nil {
						// TODO:使用统一的log
						log.Println(err)
						continue
					}
					w.Send(NewEvent(data, PUT))
				case mvccpb.DELETE:
					id := w.idFromKey(string(evt.Kv.Key))
					if id == "" {
						log.Println("无效的id")
						continue
					}
					data.Id = id
					w.Send(NewEvent(data, DEL))
				}
			}
		}
	}
}

func (w *watcher) Events(repo chan<- *Event) {
	w.repo = append(w.repo, repo)
}

func (w watcher) idFromKey(key string) string {
	if w.prefix {
		if len(key) < len(w.path) {
			return ""
		}
		return key[len(w.path):]
	}
	temps := strings.Split("/", key)
	return temps[len(temps)-1]
}

// LoadExist
func (w watcher) Exist() ([]*Event, error) {
	kv := clientv3.NewKV(w.client)
	// 获取所有已经在连接的节点
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	resp, err := kv.Get(ctx, w.path, clientv3.WithPrefix())
	if err != nil {
		return nil, fmt.Errorf("初始化失败：%v", err)
	}
	if len(resp.Kvs) == 0 {
		return nil, nil
	}
	evts := make([]*Event, 0, len(resp.Kvs))
	for _, v := range resp.Kvs {
		data := &Data{}
		id := w.idFromKey(string(v.Key))
		err := data.Decode(v.Value)
		if err != nil || id == "" {
			fmt.Printf("无法加载节点信息：%v", id)
		}
		// 创建event
		evts = append(evts, NewEvent(data, PUT))
	}
	return evts, nil
}

func (w watcher) Send(events ...*Event) {
	for _, evt := range events {
		for _, c := range w.repo {
			c <- evt
		}
	}
}

type copyWatch struct {
	ctx context.Context

	recv chan *Event
	repo []chan<- *Event
}

func (w *copyWatch) Exist() ([]*Event, error) {
	panic("implement me")
}

func (w *copyWatch) Watch(ctx context.Context) {
	if w.ctx != nil {
		return
	}
	w.ctx = ctx

	for {
		select {
		case <-ctx.Done():
			return
		case evt := <-w.recv:
			// TODO: 判断管道是否阻塞了
			for _, repo := range w.repo {
				repo <- evt
			}
		}
	}
}

// 添加监听器
func (w *copyWatch) Events(repo chan<- *Event) {
	if w.repo == nil {
		w.repo = make([]chan<- *Event, 0, 3)
	}
	w.repo = append(w.repo, repo)
}

// 添加事件源
func (w *copyWatch) Source(source <-chan *Event) {
	go w.copy(source)
}

func (w *copyWatch) copy(recv <-chan *Event) {
	for {
		select {
		case <-w.ctx.Done():
			return
		case evt := <-recv:
			// 关闭管道信息
			if evt == nil {
				return
			}
			w.recv <- evt
		}
	}
}
