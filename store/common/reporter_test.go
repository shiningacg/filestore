package common

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestNewEtcdReporter(t *testing.T) {
	reporter := NewReporter(&EtcdConfig{
		EndPoint: []string{"127.0.0.1:2379"},
	})
	info := &NodeInfo{
		NodeId:   "aaa",
		NodeType: "store",
		GRPCAddr: "11",
	}
	reporter.UpdateInfo(info)
	go reporter.KeepAlive(context.TODO())
	time.Sleep(time.Second)
	reporter.UpdateInfo(info)
	select {}
}

type MockMaster struct{}

func (m *MockMaster) Online(info *NodeInfo) {
	fmt.Println("online:", info)
}

func (m *MockMaster) Offline(info *NodeInfo) {
	fmt.Println("offline", info)
}

func TestNewMaster(t *testing.T) {
	master := NewMaster(&EtcdConfig{
		EndPoint: []string{"127.0.0.1:2379"},
	}, "/store/")
	master.SetHandler(&MockMaster{})
	master.Run(context.TODO())
}
