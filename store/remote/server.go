package remote

import (
	"context"
	store "github.com/shiningacg/filestore"
	"github.com/shiningacg/filestore/gateway"
	"github.com/shiningacg/filestore/store/remote/rpc"
)

type StoreServer struct {
	*gateway.Gateway
	store.Store
}

func (s StoreServer) Get(ctx context.Context, uuid *rpc.UUID) (*rpc.File, error) {
	f, err := s.API().Get(uuid.UUID)
	if err != nil {
		return nil, err
	}
	return toPBFile(f), nil
}

// add方法只能通过id来添加，不能read里面的内容
func (s StoreServer) Add(ctx context.Context, file *rpc.File) (*rpc.Empty, error) {
	return nil, s.API().Add(wrapPBFile(file))
}

func (s StoreServer) Remove(ctx context.Context, uuid *rpc.UUID) (*rpc.Empty, error) {
	return nil, s.API().Remove(uuid.UUID)
}

func (s StoreServer) Space(ctx context.Context, empty *rpc.Empty) (*rpc.SpaceInfo, error) {
	return toPBSpace(s.Stats().Space()), nil
}

func (s StoreServer) Network(ctx context.Context, empty *rpc.Empty) (*rpc.NetworkInfo, error) {
	return toPBNetwork(s.Stats().Network()), nil
}

func (s StoreServer) Bandwidth(ctx context.Context, empty *rpc.Empty) (*rpc.GatewayInfo, error) {
	return toPBBandwidth(s.Stats().Bandwidth()), nil
}

func wrapPBFile(file *rpc.File) store.File {
	return PBFile{File: file}
}

func toPBFile(file store.File) *rpc.File {
	return &rpc.File{
		UUID: file.ID(),
		Url:  file.Url(),
		Size: file.Size(),
		Name: file.FileName(),
	}
}

func toPBSpace(space *store.Space) *rpc.SpaceInfo {
	return &rpc.SpaceInfo{
		Cap:   space.Cap,
		Total: space.Total,
		Free:  space.Free,
		Used:  space.Used,
	}
}

func toStoreSpace(info *rpc.SpaceInfo) *store.Space {
	return &store.Space{
		Cap:   info.Cap,
		Total: info.Total,
		Free:  info.Free,
		Used:  info.Used,
	}
}

func toPBNetwork(network *store.Network) *rpc.NetworkInfo {
	return &rpc.NetworkInfo{
		Upload:   network.Upload,
		Download: network.Download,
	}
}

func toStoreNetwork(info *rpc.NetworkInfo) *store.Network {
	return &store.Network{
		Upload:   info.Upload,
		Download: info.Download,
	}
}

func toPBBandwidth(gateway *store.Bandwidth) *rpc.GatewayInfo {
	return &rpc.GatewayInfo{
		Visit:         gateway.Visit,
		DayVisit:      gateway.DayVisit,
		HourVisit:     gateway.HourVisit,
		Bandwidth:     gateway.Bandwidth,
		DayBandwidth:  gateway.DayBandwidth,
		HourBandwidth: gateway.HourBandwidth,
	}
}

func toStoreBandwidth(info *rpc.GatewayInfo) *store.Bandwidth {
	return &store.Bandwidth{
		Visit:         info.Visit,
		DayVisit:      info.DayVisit,
		HourVisit:     info.HourVisit,
		Bandwidth:     info.Bandwidth,
		DayBandwidth:  info.DayVisit,
		HourBandwidth: info.HourBandwidth,
	}
}

type PBFile struct {
	*rpc.File
}

func (f PBFile) Read(p []byte) (n int, err error) {
	return 0, nil
}

func (f PBFile) Seek(offset int64, whence int) (int64, error) {
	return 0, nil
}

func (f PBFile) Close() error {
	return nil
}

func (f PBFile) FileName() string {
	return f.Name
}

func (f PBFile) ID() string {
	return f.UUID
}

func (f PBFile) Url() string {
	return f.File.Url
}

func (f PBFile) Size() uint64 {
	return f.File.Size
}
