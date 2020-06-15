package http

import (
	"context"
	"github.com/shiningacg/apicore"
	store "github.com/shiningacg/filestore"
)

// 单次反馈记录
type Record struct {
	RequestID string
	Ip        string
	FileID    string
	Bandwidth uint64
	Finish    bool
}

type Gateway struct {
	// 反馈chan
	*Bandwidth
	host string
	md   *Global
}

func NewGateway(host string, s store.Store) *Gateway {
	bd := NewHttpBandwidth(context.TODO())
	g := &Gateway{
		Bandwidth: bd,
		host:      host,
		md:        NewMiddleware(s, bd.Chan()),
	}
	apicore.AddMiddleware(func() apicore.MiddleWare {
		return g.md
	})
	return g
}

func (g *Gateway) Upload() uint64 {
	return g.Bandwidth.current
}

func (g *Gateway) TotalUpload() uint64 {
	return g.Bandwidth.total
}

func (g *Gateway) RunHttp(host string) error {
	return apicore.Run(host)
}
