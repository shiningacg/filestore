package gateway

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"io"
	"net/http"
	"os"
)

const (
	MaxUploadSize = 1024 * 1024 * 1024 * 10
	BufferSize    = 1024 * 1024 * 128
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
	file, err := h.api.Get(fid)
	if err != nil {
		writeError(w, 400, ErrFileNotFound)
		return
	}
	// 设置head为attachment
	w.Header().Set("Content-Disposition", "attachment; filename="+file.FileName())
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
	//info := h.getUploadInfo(token)
	// 尝试读取文件
	file := h.getFile2(r)
	if file == nil {
		writeError(w, 400, ErrReadFormFile)
		return
	}
	// 创建临时文件
	f, err := os.Create(token)
	if err != nil {
		h.log.Println(err)
		writeError(w, 500, ErrInternalServer)
		return
	}
	_, err = h.copyWithLimit(MaxUploadSize, &Record{RequestID: token, FileID: token}, f, file)
	if err != nil {
		writeError(w, 400, ErrReadSocket)
	}
	// 重置位置
	f.Seek(0, io.SeekStart)
	err = h.api.Add(&OSFile{name: file.Filename, uuid: token, File: f})
	if err != nil {
		h.log.Println(err)
		writeError(w, 400, ErrReadFormFile)
		return
	}
	// 删除缓存文件
	f.Close()
	err = os.Remove(token)
	if err != nil {
		h.log.Println(err)
	}
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
	h.log.Println(err)
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

// getFile 从request中
func (h *HttpServer) getFile(r *http.Request) *PartFile {
	// 尝试读取文件,只读第一部分
	part, err := r.MultipartReader()
	file, err := part.NextPart()
	// TODO: file可能在err为nil的情况下为nil
	if err != nil {
		h.log.Println(err)
		return nil
	}
	if file.FileName() == "" {
		return nil
	}
	return &PartFile{
		Part: file,
	}
}

func (h *HttpServer) getFile2(r *http.Request) *File {
	file, header, err := r.FormFile("file")
	if err != nil {
		h.log.Println(err)
		return nil
	}
	return &File{
		File:       file,
		FileHeader: header,
	}
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
