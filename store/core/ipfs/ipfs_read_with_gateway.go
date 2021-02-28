package ipfs

import (
	"errors"
	"fmt"
	fs "github.com/shiningacg/filestore"
	ipfs "github.com/shiningacg/sn-ipfs"
	"io"
	"log"
	"net/http"
)

func NewReadWithGatewayStore(logger *log.Logger) (*ReadWithGatewayStore, error) {
	st, err := NewCore(logger)
	if err != nil {
		return nil, err
	}
	return &ReadWithGatewayStore{Store: st}, nil
}

type ReadWithGatewayStore struct {
	*Store
}

func (s *ReadWithGatewayStore) Get(uuid string) (fs.ReadableFile, error) {
	node := s.Store.ipfs.Get(uuid)
	return &HttpBodyReader{Node: node}, nil
}

type HttpBodyReader struct {
	ipfs.Node
	index int64
	*http.Response
}

func (h *HttpBodyReader) open() error {
	if h.Response != nil {
		return errors.New("原链接没有关闭")
	}
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://127.0.0.1:8080/ipfs/%v", h.Node.Cid()), nil)
	if err != nil {
		return err
	}
	req.Header.Add("Range", fmt.Sprintf("bytes=%v-", h.index))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	h.Response = resp
	return nil
}

func (h *HttpBodyReader) WriteTo(w io.Writer) (n int64, err error) {
	if h.Response == nil {
		err = h.open()
	}
	if err != nil {
		return 0, err
	}
	return io.Copy(w, h.Response.Body)
}

func (h *HttpBodyReader) UUID() string {
	return h.Cid()
}

func (h *HttpBodyReader) Read(p []byte) (n int, err error) {
	if h.Response == nil {
		h.open()
	}
	return h.Response.Body.Read(p)
}

func (h *HttpBodyReader) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case io.SeekStart:
		if uint64(offset) < h.Node.Size() {
			h.index = offset
			return offset, h.Close()
		}
		return 0, errors.New("无效的区域")
	}
	panic("未实现的seek功能被调用")
}

func (h *HttpBodyReader) Close() error {
	if h.Response != nil {
		return h.Response.Body.Close()
		h.Response = nil
	}
	return nil
}

func (h *HttpBodyReader) SetUUID(uuid string) {
	panic("implement me")
}

func (h *HttpBodyReader) SetName(name string) {
	panic("implement me")
}

func (h *HttpBodyReader) SetSize(size uint64) {
	panic("implement me")
}
