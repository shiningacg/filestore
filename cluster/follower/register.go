package follower

import (
	"context"
	"errors"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/shiningacg/filestore/cluster"
	"log"
)

type Register interface {
	Register() error
	Deregister()
	GetData() (*cluster.Data, error)
}

func NewRegister(ctx context.Context, client *clientv3.Client, data *cluster.Data, service cluster.Service) Register {
	return &register{
		ctx:     ctx,
		Data:    data,
		Client:  client,
		Service: service,
	}
}

type register struct {
	ctx context.Context
	*clientv3.Client

	cluster.Service
	*cluster.Data

	leaseID   clientv3.LeaseID
	lctx      context.Context
	lcf       func()
	lrespChan <-chan *clientv3.LeaseKeepAliveResponse

	wctx context.Context
	wcf  func()
}

func (r *register) Register() error {
	kv := clientv3.NewKV(r.Client)
	if _, err := r.GetData(); err == nil {
		return fmt.Errorf("id重复出现")
	}

	lease := clientv3.NewLease(r.Client)
	leaseResp, err := lease.Grant(r.ctx, int64(r.Service.TTL.Seconds()))
	if err != nil {
		return fmt.Errorf("无法创建租约")
	}
	r.leaseID = leaseResp.ID

	_, err = kv.Put(r.ctx, r.ToKey(), string(r.Data.Encode()), clientv3.WithLease(r.leaseID))
	if err != nil {
		return fmt.Errorf("无法写入数据")
	}

	r.lctx, r.lcf = context.WithCancel(r.ctx)
	r.lrespChan, err = lease.KeepAlive(r.lctx, r.leaseID)
	if err != nil {
		return fmt.Errorf("续租失败")
	}
	go func() {
		for {
			select {
			case r := <-r.lrespChan:
				// avoid dead loop when channel was closed
				if r == nil {
					return
				}
			case <-r.lctx.Done():
				return
			}
		}
	}()
	return nil
}

func (r *register) Watch(ch chan<- struct{}) {
	r.wctx, r.wcf = context.WithCancel(r.ctx)
	watcher := clientv3.NewWatcher(r.Client)

	wch := watcher.Watch(r.wctx, r.Service.ToKey())
	for wr := range wch {
		if wr.Canceled {
			return
		}
		err := r.updateData()
		if err != nil {
			log.Printf("无法更新数据：%v", err)
			continue
		}
		ch <- struct{}{}
	}
}

func (r *register) GetData() (*cluster.Data, error) {
	var data = &cluster.Data{}
	resp, err := r.KV.Get(r.ctx, r.ToKey())
	if err != nil {
		return nil, err
	}
	// 如果没有key会怎样？index会不会越界
	if len(resp.Kvs) == 0 {
		return nil, errors.New("没有数据")
	}
	return data, data.Decode(resp.Kvs[0].Value)
}

func (r *register) updateData() error {
	data, err := r.GetData()
	if err != nil {
		return err
	}
	r.Data = data
	return nil
}

func (r *register) Deregister() {
	if r.lcf != nil {
		r.lcf()
	}
	if r.wcf != nil {
		r.wcf()
	}
}
