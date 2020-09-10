package remote

import (
	"context"
	"errors"
	fs "github.com/shiningacg/filestore"
	"github.com/shiningacg/filestore/store/remote/rpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"time"
)

func NewRemoteStore(addr string) (*Store, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	if status := conn.GetState(); status != connectivity.Ready {
		ctx, cf := context.WithTimeout(context.Background(), time.Second*1)
		defer cf()
		// 等待连接建立
		for {
			status := conn.GetState()
			if status == connectivity.Ready {
				break
			}
			// 超时
			if !conn.WaitForStateChange(ctx, status) {
				break
			}
			status = conn.GetState()
			// 状态变化后检测状态
			if status == connectivity.TransientFailure {
				break
			}
			if status != connectivity.Ready {
				continue
			}
		}
		if conn.GetState() != connectivity.Ready {
			return nil, errors.New("无法建立连接")
		}
	}
	return &Store{
		RemoteStoreClient: rpc.NewRemoteStoreClient(conn),
		conn:              conn,
	}, nil
}

type Store struct {
	rpc.RemoteStoreClient
	conn *grpc.ClientConn
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

func (s *Store) Close() error {
	return s.conn.Close()
}

func toBaseFile(file *rpc.File) fs.BaseFile {
	var bf = &fs.BaseFileStruct{}
	bf.SetUUID(file.UUID)
	bf.SetSize(file.Size)
	bf.SetName(file.Name)
	return bf
}

func toPBFile(file fs.BaseFile) *rpc.File {
	return &rpc.File{
		UUID: file.UUID(),
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
