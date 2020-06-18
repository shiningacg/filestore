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
	StartTime uint64
	EndTime   uint64
}

type Gateway struct {
	// 反馈chan
	*GatewayMonitor
	host string
	md   *Global
}

// 新建gateway
func NewGateway(ctx context.Context, host string, s store.Store) *Gateway {
	bd := NewMonitor(ctx)
	g := &Gateway{
		GatewayMonitor: bd,
		host:           host,
		md:             NewMiddleware(s, bd.Chan()),
	}
	apicore.AddMiddleware(func() apicore.MiddleWare {
		return g.md
	})
	return g
}

func (g *Gateway) RunHttp(host string) error {
	return apicore.Run(host)
}
