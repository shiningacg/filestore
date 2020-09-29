package main

import (
	"container/ring"
	"context"
	"fmt"
	netStats "github.com/shirou/gopsutil/net"
	"time"
)

type SpeedStats struct {
	Minute  uint64
	FiveSec uint64
	Second  uint64
}

type NetworkStatus struct {
	Upload   SpeedStats
	Download SpeedStats
}

type Network interface {
	Stats() *NetworkStatus
}

// 每一次查询时系统的总收发数据数量
type speedRecord struct {
	Upload   uint64
	Download uint64
}

func (s *speedRecord) Reset() {
	s.Upload = 0
	s.Download = 0
}

func (s *speedRecord) Empty() bool {
	return s.Upload == 0 && s.Download == 0
}

func NewDefaultNetwork(ctx context.Context) *DefaultNetwork {
	maxRecordTime := 60
	network := &DefaultNetwork{
		Ring: ring.New(maxRecordTime + 1),
		len:  maxRecordTime + 1,
	}
	// 计划任务，每秒钟收集信息
	go network.run(ctx)
	return network
}

type DefaultNetwork struct {
	*ring.Ring
	len int
}

func (d *DefaultNetwork) Stats() *NetworkStatus {
	var stats = &NetworkStatus{}
	// 还没有数据
	if d.Ring.Value == nil {
		return stats
	}
	// 当前数据
	cur := d.get(0)
	if cur.Empty() {
		return stats
	}
	// 获取一秒前的数据
	sec := d.get(-1)
	fivSec := d.get(-5)
	min := d.get(-60)
	if !sec.Empty() {
		stats.Upload.Second = cur.Upload - sec.Upload
		stats.Download.Second = cur.Download - sec.Download
	}
	if !fivSec.Empty() {
		stats.Upload.FiveSec = cur.Upload - fivSec.Upload
		stats.Download.FiveSec = cur.Download - fivSec.Download
	}
	if !min.Empty() {
		stats.Upload.Minute = cur.Upload - min.Upload
		stats.Download.Minute = cur.Download - min.Download
	}
	return stats
}

// 收集一次网络信息
func (d *DefaultNetwork) collect() error {
	rcd := d.next()
	items, err := netStats.IOCounters(false)
	if err != nil {
		return err
	}
	for _, item := range items {
		rcd.Upload += item.BytesSent
		rcd.Download += item.BytesRecv
	}
	return nil
}

func (d *DefaultNetwork) run(ctx context.Context) {
	t := time.NewTimer(time.Second)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			t.Reset(time.Second)
			err := d.collect()
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

func (d *DefaultNetwork) next() *speedRecord {
	d.Ring = d.Ring.Next()
	if d.Ring.Value == nil {
		d.Ring.Value = &speedRecord{}
	}
	rcd := d.Ring.Value.(*speedRecord)
	rcd.Reset()
	return rcd
}

// TODO： 加锁
func (d *DefaultNetwork) get(n int) *speedRecord {
	r := d.Ring.Move(n)
	if r.Value == nil {
		r.Value = &speedRecord{}
	}
	rcd := r.Value.(*speedRecord)
	return rcd
}
