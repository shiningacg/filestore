package common

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"go.etcd.io/etcd/clientv3"
	"strings"
	"time"
)

// 用于汇报节点信息
type NodeInfo struct {
	NodeId      string
	NodeType    string
	GRPCAddr    string
	GatewayAddr string
}

func (n *NodeInfo) Key() string {
	return fmt.Sprintf("/%v/%v", n.NodeType, n.NodeId)
}

func (n *NodeInfo) LoadByKey(key string) {
	args := strings.Split(key, "/")
	if len(args) != 3 {
		return
	}
	n.NodeType = args[1]
	n.NodeId = args[2]
}

func (n *NodeInfo) Encode() string {
	b, _ := json.Marshal(n)
	return string(b)
}

func (n *NodeInfo) Decode(data []byte) error {
	return json.Unmarshal(data, n)
}

type EtcdConfig struct {
	EndPoint []string
	Username string
	Password string
}

func NewMaster(config *EtcdConfig, key string) *Master {
	cl, err := clientv3.New(translateConfig(config))
	if err != nil {
		panic(fmt.Errorf("无法连接etcd:%v", err))
	}
	return &Master{
		w:   clientv3.NewWatcher(cl),
		key: key,
	}
}

type MasterHandler interface {
	Online(info *NodeInfo)
	Offline(info *NodeInfo)
}

type Master struct {
	key string
	w   clientv3.Watcher
	MasterHandler
}

func (m *Master) Run(ctx context.Context) {
	m.watch(ctx)
}

func (m *Master) SetHandler(handler MasterHandler) {
	m.MasterHandler = handler
}

func (m *Master) watch(ctx context.Context) {
	watchRespChan := m.w.Watch(ctx, m.key, clientv3.WithPrefix())
	for watchResp := range watchRespChan {
		for _, event := range watchResp.Events {
			var nodeInfo = &NodeInfo{}
			// TODO: 进行数据有效性检查
			nodeInfo.Decode(event.Kv.Value)
			switch event.Type {
			case mvccpb.PUT:
				m.Online(nodeInfo)
			case mvccpb.DELETE:
				nodeInfo.LoadByKey(string(event.Kv.Key))
				m.Offline(nodeInfo)
			}
		}
	}
}

func NewReporter(config *EtcdConfig) *Reporter {
	cl, err := clientv3.New(translateConfig(config))
	if err != nil {
		panic(fmt.Errorf("无法连接etcd:%v", err))
	}
	return &Reporter{
		client:  cl,
		lease:   clientv3.NewLease(cl),
		errChan: make(chan error, 1),
	}
}

type Reporter struct {
	client  *clientv3.Client
	lease   clientv3.Lease
	leaseId clientv3.LeaseID
	info    *NodeInfo
	errChan chan error
}

func (r *Reporter) KeepAlive(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case err := <-r.errChan:
			return err
		default:
			err := r.keepAlive()
			if err != nil {
				return err
			}
			time.Sleep(2 * time.Second)
		}
	}
}

func (r *Reporter) UpdateInfo(info *NodeInfo) {
	leaseGrantResp, err := r.lease.Grant(context.TODO(), 3)
	if err != nil {
		r.errChan <- fmt.Errorf("无法创建租约:%v", err)
	}
	leaseId := leaseGrantResp.ID
	_, err = r.client.Put(context.TODO(), info.Key(), info.Encode(), clientv3.WithLease(leaseId))
	if err != nil {
		r.errChan <- err
	}
	r.leaseId = leaseId
}

func (r *Reporter) keepAlive() error {
	if r.leaseId == 0 {
		return nil
	}
	_, err := r.lease.KeepAliveOnce(context.TODO(), r.leaseId)
	if err != nil {
		r.leaseId = 0
	}
	return err
}

func translateConfig(config *EtcdConfig) clientv3.Config {
	if len(config.EndPoint) == 0 {
		panic("etcd没有指定服务器地址")
	}
	return clientv3.Config{
		Endpoints:   config.EndPoint,
		Username:    config.Username,
		Password:    config.Password,
		DialTimeout: 2 * time.Second,
	}
}
