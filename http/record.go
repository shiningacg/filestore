package http

import (
	"context"
	"github.com/shiningacg/filestore"
	"time"
)

// GatewayMonitor是通用的模块，用来处理http网关返回的信息
type GatewayMonitor struct {
	// 输入
	ctx context.Context
	// 流量
	bandwidth        uint64
	currentBandwidth uint64
	// 访问数
	visit        uint64
	currentVisit uint64
	requests     []*Record
	input        chan *Record
}

func NewMonitor(ctx context.Context) *GatewayMonitor {
	input := make(chan *Record, 100)
	return &GatewayMonitor{ctx: ctx, input: input}
}

func (b *GatewayMonitor) Chan() chan<- *Record {
	return b.input
}

func (b *GatewayMonitor) Run() {
	for {
		select {
		case r := <-b.input:
			b.addRecord(r)
		case <-b.ctx.Done():
			return
		}
	}
}

func (b *GatewayMonitor) Gateway() filestore.Bandwidth {
	b.delTimeout()
	return filestore.Bandwidth{
		Visit:        b.visit,
		DayVisit:     b.currentVisit,
		Bandwidth:    b.bandwidth,
		DayBandwidth: b.currentBandwidth,
	}
}

func (b *GatewayMonitor) addRecord(r *Record) {
	b.requests = append(b.requests, r)
	b.visit += 1
	b.currentVisit += 1
	b.bandwidth += r.Bandwidth
	b.currentBandwidth += r.Bandwidth
}

// 删除过期的信息
func (b *GatewayMonitor) delTimeout() {
	var updateTime int64 = 24 * 60 * 60
	var index int
	timeline := time.Now().Unix() - updateTime

	for i, r := range b.requests {
		if r.EndTime > uint64(timeline) {
			break
		}
		b.currentBandwidth -= r.Bandwidth
		b.currentVisit -= 1
		index = i
	}
	b.requests = b.requests[index+1:]
}
