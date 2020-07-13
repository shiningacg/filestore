package rpc

import (
	"context"
	"errors"
	store "github.com/shiningacg/filestore"
	"github.com/shiningacg/filestore/gateway"
)

type StoreServer struct {
	*gateway.Gateway
	store.API
}

func (s *StoreServer) Delete(ctx context.Context, request *DeleteRequest) (*DeleteReply, error) {
	if request == nil || request.UUID == "" {
		return nil, errors.New("无效数据")
	}
	return nil, s.API.Remove(request.UUID)
}

func (s *StoreServer) GetUrl(ctx context.Context, request *GetUrlRequest) (*GetUrlReply, error) {
	if request == nil || request.UUID == "" {
		return nil, errors.New("空的请求")
	}
	file, err := s.Get(request.UUID)
	if err != nil {
		return nil, err
	}
	return &GetUrlReply{Url: file.Url()}, nil
}

func (s *StoreServer) PostUrl(ctx context.Context, request *PostUrlRequest) (*PostUrlReply, error) {
	if request == nil || request.UUID == "" {
		return nil, errors.New("空的请求")
	}
	return &PostUrlReply{Url: s.Gateway.PostUrl(request.UUID)}, nil
}
