package common

import (
	"github.com/shirou/gopsutil/disk"
)

type DiskStats struct {
	Total uint64
	Used  uint64
}

func DiskUsage(path string) *DiskStats {
	s := &DiskStats{}
	if stat, err := disk.Usage(path); err == nil {
		s.Used = stat.Used
		s.Total = stat.Total
	}
	return s
}
