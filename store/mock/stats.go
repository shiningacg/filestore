package mock

import store "github.com/shiningacg/filestore"

type Stats API

func (s Stats) Space() *store.Space {
	return &store.Space{
		Cap:   111,
		Total: 222,
		Free:  200,
		Used:  22,
	}
}

func (s Stats) Network() *store.Network {
	return &store.Network{
		Upload:   1000,
		Download: 2000,
	}
}

func (s Stats) Bandwidth() *store.Gateway {
	return &store.Gateway{
		Visit:         3,
		DayVisit:      2,
		HourVisit:     1,
		Bandwidth:     1000,
		DayBandwidth:  200,
		HourBandwidth: 100,
	}
}
