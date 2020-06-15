package http

import (
	"context"
	"time"
)

// httpBandwidth是通用的模块，用来处理http返回的信息
type Bandwidth struct {
	// 输入
	ctx      context.Context
	total    uint64
	current  uint64
	lastTime time.Time
	input    chan *Record
}

func NewHttpBandwidth(ctx context.Context) *Bandwidth {
	input := make(chan *Record, 100)
	return &Bandwidth{ctx: ctx, input: input}
}

func (b *Bandwidth) Chan() chan<- *Record {
	return b.input
}

func (b *Bandwidth) Run() {
	for {
		select {
		case r := <-b.input:
			b.total += r.Bandwidth
			b.setCurrent(r)
		case <-b.ctx.Done():
			return
		}
	}
}

func (b *Bandwidth) setCurrent(r *Record) {
	now := time.Now()
	if now.Unix() != b.lastTime.Unix() {
		b.lastTime = now
	}
	b.current += r.Bandwidth
}
