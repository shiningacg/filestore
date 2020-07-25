package remote

import (
	"context"
	"errors"
	fs "github.com/shiningacg/filestore"
	"github.com/shiningacg/filestore/store/remote/rpc"
	"google.golang.org/grpc"
	"log"
	"net"
)

func NewStoreGRPCServer(addr string, adder Adder, fs fs.FileStore) *StoreServer {
	ss := &StoreServer{
		addr:      addr,
		Adder:     adder,
		FileStore: fs,
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
	return ss
}

type StoreServer struct {
	addr string
	Adder
	fs.FileStore
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
	bf.SetUrl(file.Url)
	return bf
}

func toPBSpace(space *fs.Space) *rpc.SpaceInfo {
	return &rpc.SpaceInfo{
		Cap:   space.Cap,
		Total: space.Total,
		Free:  space.Free,
		Used:  space.Used,
	}
}

func toPBNetwork(network *fs.Network) *rpc.NetworkInfo {
	return &rpc.NetworkInfo{
		Upload:   network.Upload,
		Download: network.Download,
	}
}

func toPBBandwidth(gateway *fs.Bandwidth) *rpc.GatewayInfo {
	return &rpc.GatewayInfo{
		Visit:         gateway.Visit,
		DayVisit:      gateway.DayVisit,
		HourVisit:     gateway.HourVisit,
		Bandwidth:     gateway.Bandwidth,
		DayBandwidth:  gateway.DayBandwidth,
		HourBandwidth: gateway.HourBandwidth,
	}
}
