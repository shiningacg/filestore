package remote

import (
	"context"
	"errors"
	fs "github.com/shiningacg/filestore"
	"github.com/shiningacg/filestore/gateway"
	"github.com/shiningacg/filestore/store/common"
	"github.com/shiningacg/filestore/store/remote/rpc"
	"google.golang.org/grpc"
	"log"
	"net"
)

func NewStoreGRPCServer(addr string, g gateway.Gateway, adder Adder, fs fs.FileStore, r *common.Reporter) *StoreServer {
	ss := &StoreServer{
		addr:      addr,
		g:         g,
		Adder:     adder,
		FileStore: fs,
		Reporter:  r,
	}
	sk, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("无法监听地址: %v %v", addr, err)
	}
	s := grpc.NewServer()
	rpc.RegisterRemoteStoreServer(s, ss)
	go func() {
		err := s.Serve(sk)
		if err != nil {
			panic(err)
		}
	}()
	go func() {
		ss.UpdateInfo(ss.report())
		err := ss.KeepAlive(context.TODO())
		if err != nil {
			panic(err)
		}
	}()
	return ss
}

type StoreServer struct {
	addr string
	g    gateway.Gateway
	Adder
	fs.FileStore
	*common.Reporter
}

// TODO: 自动识别ip
func (s StoreServer) report() *common.NodeInfo {
	return &common.NodeInfo{
		NodeId:      "center",
		NodeType:    "store",
		GRPCAddr:    s.addr,
		GatewayAddr: s.g.Host(),
	}
}

func (s StoreServer) Get(ctx context.Context, uuid *rpc.UUID) (*rpc.File, error) {
	f, err := s.FileStore.Get(uuid.UUID)
	if err != nil {
		return nil, err
	}
	return toPBFile(f), nil
}

// add方法只能通过id来添加，不能read里面的内容
func (s StoreServer) Add(ctx context.Context, file *rpc.File) (*rpc.Empty, error) {
	rf := s.Find(wrapPBFile(file))
	if rf == nil {
		return nil, errors.New("无法下载指定文件")
	}
	return &rpc.Empty{}, s.FileStore.Add(rf)
}

func (s StoreServer) Remove(ctx context.Context, uuid *rpc.UUID) (*rpc.Empty, error) {
	return &rpc.Empty{}, s.FileStore.Remove(uuid.UUID)
}

func (s StoreServer) Space(ctx context.Context, empty *rpc.Empty) (*rpc.SpaceInfo, error) {
	return toPBSpace(s.FileStore.Space()), nil
}

func (s StoreServer) Network(ctx context.Context, empty *rpc.Empty) (*rpc.NetworkInfo, error) {
	return toPBNetwork(s.FileStore.Network()), nil
}

func (s StoreServer) Bandwidth(ctx context.Context, empty *rpc.Empty) (*rpc.GatewayInfo, error) {
	return toPBBandwidth(s.FileStore.Gateway()), nil
}

func wrapPBFile(file *rpc.File) fs.BaseFile {
	var bf = &fs.BaseFileStruct{}
	bf.SetUUID(file.UUID)
	bf.SetName(file.Name)
	bf.SetSize(file.Size)
	return bf
}

func toPBSpace(space *fs.Space) *rpc.SpaceInfo {
	if space == nil {
		return &rpc.SpaceInfo{}
	}
	return &rpc.SpaceInfo{
		Cap:   space.Cap,
		Total: space.Total,
		Free:  space.Free,
		Used:  space.Used,
	}
}

func toPBNetwork(network *fs.Network) *rpc.NetworkInfo {
	if network == nil {
		return &rpc.NetworkInfo{}
	}
	return &rpc.NetworkInfo{
		Upload:   network.Upload,
		Download: network.Download,
	}
}

func toPBBandwidth(bandwidth *fs.Bandwidth) *rpc.GatewayInfo {
	if bandwidth == nil {
		return &rpc.GatewayInfo{}
	}
	return &rpc.GatewayInfo{
		Visit:         bandwidth.Visit,
		DayVisit:      bandwidth.DayVisit,
		HourVisit:     bandwidth.HourVisit,
		Bandwidth:     bandwidth.Bandwidth,
		DayBandwidth:  bandwidth.DayBandwidth,
		HourBandwidth: bandwidth.HourBandwidth,
	}
}
