package os

import (
	store "github.com/shiningacg/filestore"
	"github.com/shiningacg/filestore/store/common"
)

type Stats Store

func (s Stats) Space() *store.Space {
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

func (s Stats) Network() *store.Network {
	panic("implement me")
}

func (s Stats) Bandwidth() *store.Gateway {
	return &store.Gateway{
		Visit:         0,
		DayVisit:      0,
		HourVisit:     0,
		Bandwidth:     0,
		DayBandwidth:  0,
		HourBandwidth: 0,
	}
}
