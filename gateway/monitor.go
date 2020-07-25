package gateway

import (
	"context"
	"github.com/shiningacg/filestore"
	"io"
	"time"
)

const (
	DAY  uint64 = 60 * 60 * 24
	HOUR uint64 = 60 * 60
)

// GatewayMonitor是通用的模块，用来处理http网关返回的信息
type DefaultMonitor struct {
	// 统计总共的访问数据，每天和小时数据是计算生成的
	visit     uint64
	bandwidth uint64
	// 输入
	ctx context.Context
	// 保存请求记录
	records []*Record
	// 存入记录的chan
	input  chan *Record
	closed bool
}

func (b *DefaultMonitor) Bandwidth() *filestore.Bandwidth {
	b.delTimeout()
	calculate := func(rcds []*Record) uint64 {
		var bandwidth uint64
		for _, rcd := range rcds {
			bandwidth += rcd.Bandwidth
		}
		return bandwidth
	}
	hourRecords := b.getRecord(HOUR)
	hourBandwidth := calculate(hourRecords)
	hourVisit := len(hourRecords)
	dayVisit := len(b.records)
	dayBandwidth := calculate(b.records)
	return &filestore.Bandwidth{
		Visit:         b.visit,
		DayVisit:      uint64(dayVisit),
		HourVisit:     uint64(hourVisit),
		Bandwidth:     b.bandwidth,
		DayBandwidth:  dayBandwidth,
		HourBandwidth: hourBandwidth,
	}
}

func NewMonitor(ctx context.Context) *DefaultMonitor {
	input := make(chan *Record, 100)
	return &DefaultMonitor{
		ctx:     ctx,
		records: make([]*Record, 0, 1000),
		input:   input,
	}
}

// 协助拷贝数据，同时进行流量记录
func (b *DefaultMonitor) Copy(maxSize uint64, r *Record, dst io.Writer, src io.Reader) (uint64, error) {
	var total uint64
	b.addRecord(&Record{
		RequestID: r.RequestID,
		Ip:        r.Ip,
		FileID:    r.FileID,
		StartTime: uint64(time.Now().Unix()),
	})
	n, err := copy(dst, src, func(i int) bool {
		b.AddRecord(&Record{RequestID: r.RequestID, Bandwidth: uint64(i)})
		total += uint64(i)
		if maxSize == 0 {
			return true
		}
		if total >= maxSize {
			return false
		}
		return true
	})
	b.AddRecord(&Record{RequestID: r.RequestID, EndTime: uint64(time.Now().Unix())})
	return n, err
}

// 多线程安全添加记录
func (b *DefaultMonitor) AddRecord(record *Record) {
	if b.closed {
		return
	}
	b.input <- record
}

// 启动goroutine单线程处理记录
func (b *DefaultMonitor) Run() {
	// 开启定时任务
	for {
		select {
		case r := <-b.input:
			b.addRecord(r)
		case <-b.ctx.Done():
			close(b.input)
			b.closed = true
			return
		}
	}
}

// addRecord 把输入的record记录下并且及时更新gateway数据
func (b *DefaultMonitor) addRecord(r *Record) {
	var record *Record
	// 通过id查找是否存在过记录,从后开始查询
	for i := len(b.records) - 1; i >= 0; i-- {
		// 已经出现过，进行合并
		if b.records[i].RequestID == r.RequestID {
			record = b.records[i]
		}
	}
	if record != nil {
		record.Bandwidth += r.Bandwidth
		// 传输任务已经完成
		if r.EndTime != 0 {
			record.EndTime = r.EndTime
		}
	} else {
		// 添加访问记录
		b.visit += 1
		// 新的任务，添加记录
		b.records = append(b.records, &Record{
			RequestID: r.RequestID,
			Ip:        r.Ip,
			FileID:    r.FileID,
			Bandwidth: r.Bandwidth,
			StartTime: r.StartTime,
			EndTime:   r.EndTime,
		})
	}
	// 添加了流量记录
	b.bandwidth += r.Bandwidth
}

// delTimeout 删除过期的信息,同时更新每天数据(默认为一天)
func (b *DefaultMonitor) delTimeout() {
	// 一天
	b.records = b.getRecord(DAY)
}

// 获取截止日期前的记录
func (b *DefaultMonitor) getRecord(duration uint64) []*Record {
	var index int
	timeline := uint64(time.Now().Unix()) - duration
	for _, r := range b.records {
		if r.EndTime > timeline {
			break
		}
		index++
	}

	return b.records[index:]
}

func copy(dst io.Writer, src io.Reader, stop func(int) bool) (uint64, error) {
	var n, total int
	var err error
	// 创建缓存
	var buffer = make([]byte, BufferSize)
	for stop(n) {
		var wt, w int
		n, err = src.Read(buffer)
		if err == io.EOF {
			err = nil
			break
		}
		if err != nil {
			break
		}
		// 写入dst
		for {
			w, err = dst.Write(buffer[wt:n])
			if err != nil {
				break
			}
			wt += w
			if wt == n {
				break
			}
		}
		// 计算总和
		total += n
	}
	// 出现错误
	if err != nil {
		return 0, err
	}
	return uint64(total), nil
}

// Record 单次反馈数据
type Record struct {
	RequestID string
	Ip        string
	FileID    string
	Bandwidth uint64
	StartTime uint64
	EndTime   uint64
}
