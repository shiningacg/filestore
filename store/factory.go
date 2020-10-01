package store

import (
	"context"
	"errors"
	fs "github.com/shiningacg/filestore"
	"github.com/shiningacg/filestore/gateway"
	"github.com/shiningacg/filestore/store/common"
	"github.com/shiningacg/filestore/store/core/ipfs"
	"github.com/shiningacg/filestore/store/core/os"
	"log"
	os2 "os"
)

type Core fs.FileFS

type store struct {
	cfg Config
	Core
	gtw fs.Gateway
	net common.Network
}

func (s *store) Space() *fs.Space {
	if core, ok := s.Core.(*os.Store); ok {
		return core.Space()
	}
	stats := common.DiskUsage(s.cfg.Path)
	// cap是当前磁盘的容量
	return &fs.Space{
		Cap:   stats.Total - stats.Used,
		Total: stats.Total,
		Free:  stats.Total - stats.Used,
		Used:  stats.Used,
	}
}

func (s *store) Gateway() *fs.Bandwidth {
	return s.gtw.BandWidth()
}

func (s *store) Network() *fs.Network {
	cur := s.net.Stats()
	return &fs.Network{
		Upload:   cur.Upload.FiveSec,
		Download: cur.Download.FiveSec,
	}
}

func (s *store) Run(ctx context.Context) error {
	err := s.gtw.Run(ctx)
	if err != nil {
		panic(err)
	}
	return nil
}

// 通过外层保证
func NewStore(config Config) (*store, error) {
	// 创建core对象
	core, err := newCore(&config, nil)
	if err != nil {
		return nil, err
	}
	// 创建checker
	gtw, err := gateway.NewMyginGateway(config.Gateway, gateway.GRPC, config.CheckerAddr, "")
	if err != nil {
		return nil, err
	}
	store := Combine(gtw, core)
	store.cfg = config
	return store, nil
}

func Combine(gtw fs.Gateway, core Core) *store {
	network := common.NewDefaultNetwork(context.Background())
	gtw.SetStore(core)
	return &store{
		Core: core,
		gtw:  gtw,
		net:  network,
	}
}

// TODO: 选用一个更好的log
func newCore(config *Config, logger *log.Logger) (Core, error) {
	if logger == nil {
		logger = log.New(os2.Stdout, "", log.Ldate)
	}
	switch config.Type {
	case IPFS:
		return ipfs.NewCore(logger)
	case OS:
		return os.NewCore(config.Path, logger)
	}
	return nil, errors.New("未知仓库")
}
