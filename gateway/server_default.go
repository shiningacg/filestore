package gateway

import (
	"context"
	"errors"
	"github.com/google/uuid"
	fs "github.com/shiningacg/filestore"
	"github.com/shiningacg/filestore/gateway/checker"
	monitor2 "github.com/shiningacg/filestore/gateway/monitor"
	"github.com/shiningacg/mygin-frame-libs/log"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

func NewDefaultGateway(addr string, checker checker.Checker, logger *log.Logger) *DefaultGateway {
	return &DefaultGateway{
		log:     logger,
		addr:    addr,
		checker: checker,
		monitor: monitor2.NewMonitor(),
	}
}

type DefaultGateway struct {
	checker checker.Checker
	// 负责日志控制
	log *log.Logger
	// 负责数据统计
	monitor *monitor2.DefaultMonitor
	addr    string
	// 存放文件的仓库，能够通过id存放和获取文件
	fs fs.FileFS
}

// 获取统计信息
func (g *DefaultGateway) BandWidth() *fs.Bandwidth {
	return g.monitor.Bandwidth()
}

func (g *DefaultGateway) Run(ctx context.Context) error {
	if g.fs == nil {
		panic("空的仓库")
	}
	go g.monitor.Run(ctx)
	return http.ListenAndServe(g.addr, g)
}

func (g *DefaultGateway) SetStore(store fs.FileStore) {
	g.fs = store
}

// 传入一个uuid，返回下载地址
func (g *DefaultGateway) Host() string {
	if g.addr[0] == ':' {
		return "0.0.0.0" + g.addr
	}
	return g.addr
}

const (
	BufferSize = 1024 * 1024 * 128
)

var (
	ErrAction         = errors.New("未知操作")
	ErrFileNotFound   = errors.New("没有找到文件")
	ErrReadFormFile   = errors.New("无法读取发送的文件")
	ErrInternalServer = errors.New("服务器错误")
	ErrReadSocket     = errors.New("传输失败")
	ErrInvalidToken   = errors.New("无效的token")
)

// Upload 处理用户的上传请求
func (g *DefaultGateway) Download(w http.ResponseWriter, r *http.Request) {
	// 尝试获取uuid
	requestID := g.generateRequestID()
	fid := g.getUUID(r)
	if fid == "" {
		g.writeError(w, 400, ErrFileNotFound)
		return
	}
	file, err := g.fs.Get(fid)
	if err != nil {
		g.writeError(w, 400, ErrFileNotFound)
		return
	}
	// 设置head为attachment
	w.Header().Set("Content-Disposition", "attachment; filename="+file.Name())
	// 开始传输文件
	_, err = g.copyWithoutLimit(&monitor2.Record{RequestID: requestID, FileID: fid}, w, file)
	if err != nil {
		// 判断socket是否关闭
		// 打日志
		g.log.Fatal(err)
		g.writeError(w, 400, ErrReadSocket)
		return
	}
}

// Download 处理用户的下载请求
func (g *DefaultGateway) Upload(w http.ResponseWriter, r *http.Request) {
	// 判断token是否有效，同时获取最大上传限制
	token := g.getToken(r)
	if token == "" {
		g.writeError(w, 400, ErrInvalidToken)
		return
	}
	checkResult, err := g.checker.Get(token)
	if err != nil {
		g.writeError(w, 400, ErrInvalidToken)
		return
	}
	// 尝试读取文件
	file, header, err := g.getFile(r)
	if err != nil {
		g.writeError(w, 400, ErrReadFormFile)
		return
	}
	// 创建临时文件，可以考虑弄一个函数
	f, err := os.Create(token)
	if err != nil {
		g.log.Fatal(err)
		g.writeError(w, 500, ErrInternalServer)
		return
	}
	// 删除缓存文件
	defer func() {
		f.Close()
		err = os.Remove(token)
		if err != nil {
			g.log.Fatal(err)
		}
	}()
	// 开始读取
	size, err := g.copyWithLimit(checkResult.Size, &monitor2.Record{RequestID: token, FileID: token}, f, file)
	if err == monitor2.ErrReachMaxSize {
		g.writeError(w, 400, err)
		return
	}
	if err != nil {
		g.writeError(w, 400, ErrReadSocket)
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
		g.log.Fatal(err)
		g.writeError(w, 400, ErrReadFormFile)
		return
	}
	err = g.checker.Set(checkResult.Checked(rf.UUID()))
	if err != nil {
		g.log.Fatal(err)
		g.writeError(w, 500, ErrInternalServer)
		return
	}
	// 写入回复
	g.writeSucResponse(w)
}

func (g *DefaultGateway) generateRequestID() string {
	return uuid.New().String()
}

func (g *DefaultGateway) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	url := request.RequestURI
	action := g.getAction(url)
	switch action {
	case "upload":
		g.Upload(writer, request)
	case "download":
		g.Download(writer, request)
	default:
		g.writeErrorResponse(writer, 400, ErrAction)
	}
}

// copyWithLimit 复制内容，但大小不会超过MaxUploadSize
func (g *DefaultGateway) copyWithLimit(maxSize uint64, r *monitor2.Record, dst io.Writer, src io.Reader) (uint64, error) {
	return g.monitor.Copy(maxSize, r, dst, src)
}

// copyWithLimit 复制内容，大小不受限制
func (g *DefaultGateway) copyWithoutLimit(r *monitor2.Record, dst io.Writer, src io.Reader) (uint64, error) {
	return g.monitor.Copy(0, r, dst, src)
}

func (g *DefaultGateway) writeErrorResponse(w http.ResponseWriter, code int, err error) {
	// 输出错误日志
	g.log.Fatal(err)
	w.WriteHeader(code)
	_, _ = w.Write([]byte(err.Error()))
}

func (g *DefaultGateway) writeSucResponse(w http.ResponseWriter) {
	w.WriteHeader(200)
}

// getAction 尝试通过url的前缀判断用户想要进行的操作
func (g *DefaultGateway) getAction(url string) string {
	// /post/ssss && /get/xxxx
	if len(url) < 6 {
		return ""
	}
	if url[1:4] == "get" {
		return "download"
	} else if url[1:5] == "post" {
		return "upload"
	}
	return ""
}

func (g *DefaultGateway) getFile(r *http.Request) (multipart.File, *multipart.FileHeader, error) {
	file, header, err := r.FormFile("file")
	if err != nil {
		g.log.Fatal(err)
		return nil, nil, err
	}
	return file, header, nil
}

// writeError 快捷回复用户消息
func (g *DefaultGateway) writeError(w http.ResponseWriter, code int, err error) {
	w.WriteHeader(code)
	_, _ = w.Write([]byte(err.Error()))
}

// getToken 获取用户请求的token
func (g *DefaultGateway) getToken(r *http.Request) string {
	// 将用户的post路径后的token取出
	return r.RequestURI[6:]
}

func (g *DefaultGateway) getUUID(r *http.Request) string {
	return r.RequestURI[5:]
}
