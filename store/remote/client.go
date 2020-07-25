package remote

import (
	"context"
	fs "github.com/shiningacg/filestore"
	"github.com/shiningacg/filestore/store/remote/rpc"
	"google.golang.org/grpc"
)

func NewRemoteStore(addr string) (*Store, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return &Store{rpc.NewRemoteStoreClient(conn)}, nil
}

type Store struct {
	rpc.RemoteStoreClient
}

func (s *Store) Get(uuid string) (fs.BaseFile, error) {
	bf, err := s.RemoteStoreClient.Get(context.TODO(), &rpc.UUID{UUID: uuid})
	if err != nil {
		return nil, err
	}
	return toBaseFile(bf), nil
}

func (s *Store) Add(file fs.BaseFile) error {
	_, err := s.RemoteStoreClient.Add(context.TODO(), toPBFile(file))
	return err
}

func (s *Store) Remove(uuid string) error {
	_, err := s.RemoteStoreClient.Remove(context.TODO(), &rpc.UUID{UUID: uuid})
	return err
}

func (s *Store) Space() *fs.Space {
	space, err := s.RemoteStoreClient.Space(context.TODO(), &rpc.Empty{})
	if err != nil {
		return nil
	}
	return toStoreSpace(space)
}

func (s *Store) Network() *fs.Network {
	nw, err := s.RemoteStoreClient.Network(context.TODO(), &rpc.Empty{})
	if err != nil {
		return nil
	}
	return toStoreNetwork(nw)
}

func (s *Store) Gateway() *fs.Bandwidth {
	bw, err := s.RemoteStoreClient.Bandwidth(context.TODO(), &rpc.Empty{})
	if err != nil {
		return nil
	}
	return toStoreBandwidth(bw)
}

func toBaseFile(file *rpc.File) fs.BaseFile {
	var bf = &fs.BaseFileStruct{}
	bf.SetUUID(file.UUID)
	bf.SetSize(file.Size)
	bf.SetName(file.Name)
	bf.SetUrl(file.Url)
	return bf
}

func toPBFile(file fs.BaseFile) *rpc.File {
	return &rpc.File{
		UUID: file.UUID(),
		Url:  file.Url(),
		Size: file.Size(),
		Name: file.Name(),
	}
}

func toStoreSpace(info *rpc.SpaceInfo) *fs.Space {
	return &fs.Space{
		Cap:   info.Cap,
		Total: info.Total,
		Free:  info.Free,
		Used:  info.Used,
	}
}

func toStoreNetwork(info *rpc.NetworkInfo) *fs.Network {
	return &fs.Network{
		Upload:   info.Upload,
		Download: info.Download,
	}
}

func toStoreBandwidth(info *rpc.GatewayInfo) *fs.Bandwidth {
	return &fs.Bandwidth{
		Visit:         info.Visit,
		DayVisit:      info.DayVisit,
		HourVisit:     info.HourVisit,
		Bandwidth:     info.Bandwidth,
		DayBandwidth:  info.DayBandwidth,
		HourBandwidth: info.HourBandwidth,
	}
}
