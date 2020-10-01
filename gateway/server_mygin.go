package gateway

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	fs "github.com/shiningacg/filestore"
	"github.com/shiningacg/filestore/gateway/checker"
	monitor2 "github.com/shiningacg/filestore/gateway/monitor"
	"github.com/shiningacg/mygin"
	"io"
	l "log"
	"mime/multipart"
	"net/http"
	"os"
)

const (
	GRPC string = "GRPC"
	HTTP string = "HTTP"
)

type MyginGateway struct {
	ctx context.Context
	cf  func()
	*mygin.Engine
	// 上传监控器
	checker checker.Checker
	// 流量监控器
	monitor *monitor2.DefaultMonitor
	// 监听地址
	addr string
	// 存放文件的仓库，能够通过id存放和获取文件
	fs fs.FileFS
}

func NewMyginGateway(addr string, checkerType, checkerAddr, checkKey string) (*MyginGateway, error) {
	var (
		ck  checker.Checker
		err error
	)
	switch checkerType {
	case GRPC:
		ck, err = checker.NewGrpcChecker(checkerAddr, checkKey)
		if err != nil {
			return nil, err
		}
	}
	return DesignMyginGateway(addr, ck), nil
}

func DesignMyginGateway(addr string, checker checker.Checker) *MyginGateway {
	hs := &MyginGateway{
		addr:    addr,
		checker: checker,
		Engine:  mygin.New(),
		monitor: monitor2.NewMonitor(),
	}
	hs.LoadRouter(hs.Engine)
	return hs
}

// 获取统计信息
func (g *MyginGateway) BandWidth() *fs.Bandwidth {
	return g.monitor.Bandwidth()
}

func (g *MyginGateway) SetStore(store fs.FileFS) {
	g.fs = store
}

func (g *MyginGateway) Host() string {
	if g.addr[0] == ':' {
		return "0.0.0.0" + g.addr
	}
	return g.addr
}

func (g *MyginGateway) Reset(addr string, checker checker.Checker) error {
	if !g.Closed() {
		return errors.New("服务还未停止")
	}
	if checker != nil {
		g.checker = checker
	}
	g.monitor = monitor2.NewMonitor()
	g.ctx = nil
	g.addr = addr
	return nil
}

func (g *MyginGateway) Closed() bool {
	var closed bool
	select {
	case <-g.ctx.Done():
		closed = true
	default:
		closed = false
	}
	return closed
}

func (g *MyginGateway) Run(ctx context.Context) error {
	if g.ctx != nil {
		return errors.New("服务已经在运行")
	}
	g.ctx = ctx
	go g.monitor.Run(ctx)
	err := http.ListenAndServe(g.addr, g.Engine)
	if err != nil {
		return err
	}
	return nil
}

func (g *MyginGateway) LoadRouter(engine *mygin.Engine) {
	r := engine.Router()
	r.Use(g.RequestID)
	r.Get("/download/:fid").Do(g.Download)
	r.Post("/upload/:token").Use().Do(g.Upload)
}

// TODO：如果有session则进行记录
func (g *MyginGateway) Download(ctx *mygin.Context) {
	// 尝试获取uuid
	requestID := ctx.Value("RequestID").(string)
	fid := ctx.RouterValue("fid")
	file, err := g.fs.Get(fid)
	if err != nil {
		fmt.Println(err)
		ctx.Status(404)
		return
	}
	// 设置为原生操作
	ctx.Proto()
	// 设置head为attachment
	writer := ctx.Write
	writer.WriteHeader(200)
	writer.Header().Set("Content-Disposition", "attachment; filename="+file.Name())
	// 开始传输文件
	_, err = g.copyWithoutLimit(&monitor2.Record{RequestID: requestID, FileID: fid}, writer, file)
	if err != nil {
		// 判断socket是否关闭
		// 打日志
		l.Fatal(err)
		return
	}
}

// TODO： 没有checker的情况下允许所有的请求
// Download 处理用户的下载请求
func (g *MyginGateway) Upload(ctx *mygin.Context) {
	// 判断token是否有效，同时获取最大上传限制
	token := ctx.RouterValue("token")
	checkResult, err := g.checker.Get(token)
	if err != nil {
		ctx.Status(400)
		return
	}
	// 尝试读取文件
	file, header, err := g.getFile(ctx.Request)
	if err != nil {
		ctx.Status(400)
		return
	}
	// 创建临时文件，可以考虑弄一个函数
	f, err := os.Create(token)
	if err != nil {
		ctx.Status(500)
		return
	}
	// 删除缓存文件
	defer func() {
		f.Close()
		err = os.Remove(token)
		if err != nil {
			l.Fatal(err)
		}
	}()
	// 开始读取
	size, err := g.copyWithLimit(checkResult.Size, &monitor2.Record{RequestID: token, FileID: token}, f, file)
	if err == monitor2.ErrReachMaxSize {
		// writeError(w, 400, err)
		ctx.Status(400)
		return
	}
	if err != nil {
		// writeError(w, 400, ErrReadSocket)
		ctx.Status(400)
		return
	}
	if size != checkResult.Size && checkResult.Size != 0 {
		// 文件大小不对
		ctx.Status(400)
		return
	}
	// 重置文件的读取位置
	f.Seek(0, io.SeekStart)
	// 记录信息
	bs := &fs.BaseFileStruct{}
	bs.SetUUID(token)
	bs.SetName(header.Filename)
	bs.SetSize(size)
	// 放入仓库中
	rf := fs.NewReadableFile(bs, f)
	err = g.fs.Add(rf)
	if err != nil {
		l.Fatal(err)
		// writeError(w, 400, ErrReadFormFile)
		ctx.Status(400)
		return
	}
	err = g.checker.Set(checkResult.Checked(rf.UUID()))
	if err != nil {
		l.Fatal(err)
		// writeError(w, 500, ErrInternalServer)
		ctx.Status(500)
		return
	}
	// 写入回复
	ctx.Body([]byte(`{"uuid":"` + rf.UUID() + `"}`))
	ctx.Status(200)
	return
}

func (g *MyginGateway) RequestID(ctx *mygin.Context) {
	rid := uuid.New().String()
	ctx.Set("RequestID", rid)
	ctx.Next()
	// 写入到数据库中
}

// copyWithLimit 复制内容，但大小不会超过MaxUploadSize
func (g *MyginGateway) copyWithLimit(maxSize uint64, r *monitor2.Record, dst io.Writer, src io.Reader) (uint64, error) {
	return g.monitor.Copy(maxSize, r, dst, src)
}

// copyWithLimit 复制内容，大小不受限制
func (g *MyginGateway) copyWithoutLimit(r *monitor2.Record, dst io.Writer, src io.Reader) (uint64, error) {
	return g.monitor.Copy(0, r, dst, src)
}

func (g *MyginGateway) getFile(r *http.Request) (multipart.File, *multipart.FileHeader, error) {
	file, header, err := r.FormFile("file")
	if err != nil {
		return nil, nil, err
	}
	return file, header, nil
}