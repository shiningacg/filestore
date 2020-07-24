package os

import (
	store "github.com/shiningacg/filestore"
	"github.com/shiningacg/filestore/store/common"
)

func (s *Store) Space() *store.Space {
	space := &store.Space{}
	dbInfo := s.db.Info()
	if dbInfo != nil {
		space.Total = dbInfo.MaxSize
		space.Used = dbInfo.UsedSize
		space.Free = dbInfo.FreeSize
	}
	diskStats := common.DiskUsage(s.storeManager.GetBasePath())
	if diskStats != nil {
		space.Cap = diskStats.Total - diskStats.Used
	}
	return space
}

func (s *Store) Network() *store.Network {
	panic("implement me")
}

func (s *Store) Gateway() *store.Bandwidth {
	return s.gateway.BandWidth()
}
