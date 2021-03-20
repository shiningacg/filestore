package monitor

import (
	"errors"
	"github.com/google/uuid"
	"io"
	"time"
)

const (
	DAY  uint64 = 60 * 60 * 24
	HOUR uint64 = 60 * 60
)

var (
	ErrReachMaxSize = errors.New("超过数据大小限制")
)

// Monitor 负责监控流量，通过包装reader或者writer来监控流量情况
type Monitor interface {
	// In 监控流入流量
	In(reader io.Reader, id string, limit int64) io.ReadCloser
	// Out 监控流出流量
	Out(writer io.Writer, id string, limit int64) io.WriteCloser
	// 总体数据
	Stats() *Stats
}

type Record struct {
	Id        string
	Size      uint64
	In        bool
	StartTime time.Time
	EndTime   time.Time
}

type Stats struct {
	In  Bandwidth
	Out Bandwidth
}

type Bandwidth struct {
	Total uint64
	// 日流出流量
	Day uint64
	// 小时流出流量
	Hour uint64
	// 分钟
	Min uint64
	// 实时数据
	Now uint64
}

func NewMonitor() *monitor {
	m := &monitor{
		rcds:    make(map[string]*Record),
		secIn:   make([]int, 10),
		secOut:  make([]int, 10),
		minIn:   make([]int, 60),
		minOut:  make([]int, 60),
		hourIn:  make([]int, 60),
		hourOut: make([]int, 60),
		dayIn:   make([]int, 24),
		dayOut:  make([]int, 24),
	}
	go m.timer()
	return m
}

type monitor struct {
	rcds map[string]*Record
	// 总传输
	totalIn  uint64
	totalOut uint64
	// 默认长度为10
	curSec uint8
	secIn  []int
	secOut []int
	// 默认长度为60
	curMin uint8
	minIn  []int
	minOut []int
	// 默认长度为60
	curHour uint8
	hourIn  []int
	hourOut []int
	// 默认长度为24
	curDay uint8
	dayIn  []int
	dayOut []int
}

func (m *monitor) timer() {
	t := time.NewTicker(time.Microsecond * 100)
	for {
		select {
		case <-t.C:
			go m.doAfterSleep()
		}
	}
}

// 进行100ms后的操作
func (m *monitor) doAfterSleep() {
	// 前进并且清空上一秒的数据
	m.curSec += 1
	// 如果满位，则进位
	if m.curSec >= 10 {
		m.curSec = 0
		m.curMin += 1
		if m.curMin >= 60 {
			m.curMin = 0
			m.curHour += 1
			if m.curHour >= 60 {
				m.curHour = 0
				m.curDay += 1
				if m.curDay >= 24 {
					m.curDay = 0
				}
				m.dayIn[m.curDay] = 0
				m.dayOut[m.curDay] = 0
			}
			m.hourIn[m.curHour] = 0
			m.hourOut[m.curHour] = 0
		}
		m.minIn[m.curMin] = 0
		m.minOut[m.curMin] = 0
	}
	m.secIn[m.curSec] = 0
	m.secOut[m.curSec] = 0
}

func (m *monitor) In(reader io.Reader, id string, limit int64) io.ReadCloser {
	if id == "" {
		id = uuid.New().String()
	}
	return &rw{
		Reader: reader,
		limit:  limit,
		r:      m.new(id, true),
		m:      m,
	}
}

func (m *monitor) Out(writer io.Writer, id string, limit int64) io.WriteCloser {
	if id == "" {
		id = uuid.New().String()
	}
	return &rw{
		Writer: writer,
		limit:  limit,
		r:      m.new(id, false),
		m:      m,
	}
}

func (m *monitor) add(size int, in bool) {
	// 操作流量记录
	if in {
		m.totalIn += uint64(size)
		m.secIn[m.curSec] += size
		m.minIn[m.curSec] += size
		m.hourIn[m.curSec] += size
		return
	}
	m.totalOut += uint64(size)
	m.secOut[m.curSec] += size
	m.minOut[m.curSec] += size
	m.hourOut[m.curSec] += size
}

func (m *monitor) new(id string, in bool) *Record {
	var rcd = &Record{
		Id:        id,
		In:        in,
		Size:      0,
		StartTime: time.Now(),
		EndTime:   time.Time{},
	}
	m.rcds[id] = rcd
	return rcd
}

func (m *monitor) Delete(id string) {
	delete(m.rcds, id)
}

func (m *monitor) Stats() *Stats {
	var secIn, secOut, minIn, minOut, hourIn, hourOut, dayIn, dayOut uint64
	for _, v := range m.secIn {
		secIn += uint64(v)
	}
	for _, v := range m.secOut {
		secOut += uint64(v)
	}
	for _, v := range m.minIn {
		minIn += uint64(v)
	}
	for _, v := range m.minOut {
		minOut += uint64(v)
	}
	for _, v := range m.hourIn {
		hourIn += uint64(v)
	}
	for _, v := range m.hourOut {
		hourOut += uint64(v)
	}
	for _, v := range m.dayIn {
		dayIn += uint64(v)
	}
	for _, v := range m.dayOut {
		dayOut += uint64(v)
	}
	return &Stats{
		In: Bandwidth{
			Total: m.totalIn,
			Day:   dayIn,
			Hour:  hourIn,
			Min:   minIn,
			Now:   secIn,
		},
		Out: Bandwidth{
			Total: m.totalOut,
			Day:   dayOut,
			Hour:  hourOut,
			Min:   minOut,
			Now:   secOut,
		},
	}
}

type rw struct {
	limit int64
	cache []byte
	io.Reader
	io.Writer
	r *Record
	m *monitor
}

func (r *rw) Read(p []byte) (n int, err error) {
	if r.cache == nil {
		r.cache = r.makeBestCache(p)
	}
	n, err = r.Reader.Read(r.cache)
	if err != nil {
		return 0, err
	}
	// 限制大小
	if r.limit != 0 && uint64(r.limit) < r.r.Size+uint64(n) {
		return 0, ErrReachMaxSize
	}
	n = copy(p, r.cache[:n])
	r.add(int64(n))
	return n, nil
}

func (r *rw) Write(p []byte) (n int, err error) {
	if r.cache == nil {
		r.cache = r.makeBestCache(p)
	}
	// 剩下的内容比预期更多
	if r.limit != 0 && len(p) > int(r.limit)-int(r.r.Size) {
		return 0, ErrReachMaxSize
	}
	n, err = r.Writer.Write(p)
	if err != nil {
		return 0, err
	}
	r.m.add(n, true)
	return n, nil
}

func (r *rw) add(size int64) {
	r.r.Size += uint64(size)
	r.m.add(int(size), r.r.In)
}

func (r *rw) makeBestCache(dst []byte) []byte {
	if size := len(dst); size < 2048*32*2 {
		return make([]byte, size)
	}
	return make([]byte, 2048*32)
}

func (r *rw) Close() error {
	if r.Reader != nil {
		if closer, ok := r.Reader.(io.Closer); ok {
			return closer.Close()
		}
	}
	if r.Writer != nil {
		if closer, ok := r.Writer.(io.Closer); ok {
			return closer.Close()
		}
	}
	r.r.EndTime = time.Now()
	return nil
}
