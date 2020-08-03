package gateway

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	fs "github.com/shiningacg/filestore"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

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

type HttpServer Gateway

// Upload 处理用户的上传请求
func (h *HttpServer) Download(w http.ResponseWriter, r *http.Request) {
	// 尝试获取uuid
	requestID := h.generateRequestID()
	fid := getUUID(r)
	if fid == "" {
		writeError(w, 400, ErrFileNotFound)
		return
	}
	file, err := h.fs.Get(fid)
	if err != nil {
		writeError(w, 400, ErrFileNotFound)
		return
	}
	// 设置head为attachment
	w.Header().Set("Content-Disposition", "attachment; filename="+file.Name())
	// 开始传输文件
	_, err = h.copyWithoutLimit(&Record{RequestID: requestID, FileID: fid}, w, file)
	if err != nil {
		// 判断socket是否关闭
		// 打日志
		writeError(w, 400, ErrReadSocket)
		return
	}
}

// Download 处理用户的下载请求
func (h *HttpServer) Upload(w http.ResponseWriter, r *http.Request) {
	// 判断token是否有效，同时获取最大上传限制
	token := getToken(r)
	if token == "" {
		writeError(w, 400, ErrInvalidToken)
		return
	}
	checkResult, err := h.checker.Get(token)
	if err != nil {
		writeError(w, 400, ErrInvalidToken)
		return
	}
	// 尝试读取文件
	file, header, err := h.getFile(r)
	if err != nil {
		writeError(w, 400, ErrReadFormFile)
		return
	}
	// 创建临时文件，可以考虑弄一个函数
	f, err := os.Create(token)
	if err != nil {
		h.log.Fatal(err)
		writeError(w, 500, ErrInternalServer)
		return
	}
	// 删除缓存文件
	defer func() {
		f.Close()
		err = os.Remove(token)
		if err != nil {
			h.log.Fatal(err)
		}
	}()
	// 开始读取
	size, err := h.copyWithLimit(checkResult.Size, &Record{RequestID: token, FileID: token}, f, file)
	if err == ErrReachMaxSize {
		writeError(w, 400, err)
		return
	}
	if err != nil {
		writeError(w, 400, ErrReadSocket)
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
	err = h.fs.Add(rf)
	if err != nil {
		h.log.Fatal(err)
		writeError(w, 400, ErrReadFormFile)
		return
	}
	err = h.checker.Set(checkResult.Checked(rf.UUID()))
	if err != nil {
		h.log.Fatal(err)
		writeError(w, 500, ErrInternalServer)
		return
	}
	// 写入回复
	h.writeSucResponse(w)
}

func (h *HttpServer) generateRequestID() string {
	return uuid.New().String()
}

func (h *HttpServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	url := request.RequestURI
	action := getAction(url)
	fmt.Println(action, url)
	switch action {
	case "upload":
		h.Upload(writer, request)
	case "download":
		h.Download(writer, request)
	default:
		h.writeErrorResponse(writer, 400, ErrAction)
	}
}

// copyWithLimit 复制内容，但大小不会超过MaxUploadSize
func (h *HttpServer) copyWithLimit(maxSize uint64, r *Record, dst io.Writer, src io.Reader) (uint64, error) {
	return h.monitor.Copy(maxSize, r, dst, src)
}

// copyWithLimit 复制内容，大小不受限制
func (h *HttpServer) copyWithoutLimit(r *Record, dst io.Writer, src io.Reader) (uint64, error) {
	return h.monitor.Copy(0, r, dst, src)
}

func (h *HttpServer) writeErrorResponse(w http.ResponseWriter, code int, err error) {
	// 输出错误日志
	h.log.Fatal(err)
	w.WriteHeader(code)
	_, _ = w.Write([]byte(err.Error()))
}

func (h *HttpServer) writeSucResponse(w http.ResponseWriter) {
	w.WriteHeader(200)
}

// getAction 尝试通过url的前缀判断用户想要进行的操作
func getAction(url string) string {
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

func (h *HttpServer) getFile(r *http.Request) (multipart.File, *multipart.FileHeader, error) {
	file, header, err := r.FormFile("file")
	if err != nil {
		h.log.Fatal(err)
		return nil, nil, err
	}
	return file, header, nil
}

// writeError 快捷回复用户消息
func writeError(w http.ResponseWriter, code int, err error) {
	w.WriteHeader(code)
	_, _ = w.Write([]byte(err.Error()))
}

// getToken 获取用户请求的token
func getToken(r *http.Request) string {
	// 将用户的post路径后的token取出
	return r.RequestURI[6:]
}

func getUUID(r *http.Request) string {
	return r.RequestURI[5:]
}
