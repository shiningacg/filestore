package store

import (
	"context"
	fs "github.com/shiningacg/filestore"
	"github.com/shiningacg/filestore/gateway"
	"github.com/shiningacg/filestore/gateway/checker"
	"github.com/shiningacg/filestore/store/common"
	"github.com/shiningacg/mygin-frame-libs/log"
)

type FactoryStore interface {
	fs.FileStore
	SetGateway(gateway gateway.Gateway)
	SetLogger(logger *log.Logger)
	GetGateway() gateway.Gateway
	GetLogger() *log.Logger
}

type store struct {
	FactoryStore
	net common.Network
}

func (s *store) Network() *fs.Network {
	cur := s.net.Stats()
	return &fs.Network{
		Upload:   cur.Upload.FiveSec,
		Download: cur.Download.FiveSec,
	}
}

// 通过外层保证
func NewStore(st FactoryStore, gatewayAddr string, checker checker.Checker, logger *log.Logger) fs.FileStore {
	gtw := gateway.NewMyginGateway(gatewayAddr, checker)
	network := common.NewDefaultNetwork(context.Background())
	st.SetGateway(gtw)
	gtw.SetStore(st)
	// 启动gateway
	go func() {
		for {
			err := gtw.Run(context.TODO())
			if err != nil {
				panic(err)
			}
		}
	}()
	return &store{
		FactoryStore: st,
		net:          network,
	}
}
