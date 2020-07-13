package gateway

import (
	"errors"
	store "github.com/shiningacg/filestore"
	"io"
	"net/http"
)

const (
	MaxUploadSize = 1024 * 1024 * 1024 * 10
	BufferSize    = 1024 * 1024 * 128
)

var (
	ErrReadFormFile   = errors.New("无法读取发送的文件")
	ErrInternalServer = errors.New("服务器错误")
	ErrReadSocket     = errors.New("传输失败")
)

type HttpServer Gateway

// Upload 处理用户的上传请求
//func (h *HttpServer) Upload(w http.ResponseWriter, r *http.Request) {
//	// 尝试获取文件
//	// 判断token是否有效，同时获取最大上传限制
//	//token := getToken(r)
//	//info := h.getUploadInfo(token)
//	// 尝试读取文件
//	file := getFile(r)
//	if file == nil {
//		writeError(w, 400, ErrReadFormFile)
//		return
//	}
//	// 创建临时文件
//	f, err := os.Create(randomString())
//	if err != nil {
//		// 打印日志
//		writeError(w, 500, ErrInternalServer)
//		return
//	}
//	// 开始传输文件
//	n, err := h.copyWithLimit(MaxUploadSize, f, file)
//	if err != nil {
//		// 判断socket是否关闭
//		// 打日志
//		writeError(w, 400, ErrReadSocket)
//		return
//	}
//	// 创建记录
//	// 回复消息
//}

// Download 处理用户的下载请求
//func (h *HttpServer) Download(w http.ResponseWriter, r *http.Request) {
//	// 尝试获取文件
//	// 判断token是否有效，同时获取最大上传限制
//	//token := getToken(r)
//	//info := h.getUploadInfo(token)
//	// 尝试读取文件
//	file := getFile(r)
//	if file == nil {
//		writeError(w, 400, ErrReadFormFile)
//		return
//	}
//	// 创建临时文件
//	f, err := os.Create(randomString())
//	if err != nil {
//		// 打印日志
//		writeError(w, 500, ErrInternalServer)
//		return
//	}
//	// 开始传输文件
//	n, err := h.copyWithLimit(MaxUploadSize, f, file)
//	if err != nil {
//		// 判断socket是否关闭
//		// 打日志
//		writeError(w, 400, ErrReadSocket)
//		return
//	}
//	// 创建记录
//	// 回复消息
//}

func (h *HttpServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	//	url := request.RequestURI
	//	action := getAction(url)
	//	switch action {
	//	case "upload":
	//		h.Upload(writer, request)
	//	case "download":
	//		h.Download(writer, request)
	//	default:
	//		writeErrorResponse(writer, err)
	//	}
}

// copyWithLimit 复制内容，但大小不会超过MaxUploadSize
func (h *HttpServer) copyWithLimit(maxSize uint64, dst io.Writer, src io.Reader) (uint64, error) {
	var (
		total, n int
		err      error
	)
	// 创建缓存
	var buffer = make([]byte, BufferSize)
	for {
		n, err = src.Read(buffer)
		if err != nil {
			break
		}
		total += n
	}
	// 出现错误
	if err != nil && err != io.EOF {
		return 0, err
	}
	return uint64(total), nil
}

// getAction 尝试通过url的前缀判断用户想要进行的操作
func getAction(url string) string {
	// /post/ssss && /get/xxxx
	if url[1:4] == "get" {
		return "download"
	} else if url[1:5] == "post" {
		return "upload"
	}
	return ""
}

// getFile 从request中
func getFile(r http.Request) store.File {
	// 尝试读取文件,只读第一部分
	file, header, err := r.FormFile("file")
	if err != nil {
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
	w.Write([]byte(err.Error()))
}

// randomString 生成随机字符串，通过uuid生成
func randomString() string {
	return "random"
}

// getToken 获取用户请求的token
func getToken(r http.Request) string {
	return ""
	//var buffer = make([]byte,1024)

}
