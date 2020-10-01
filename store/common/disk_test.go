package common

import (
	"fmt"
	"testing"
)

func TestDiskUsage(t *testing.T) {
	d := DiskUsage("/Users/shlande/Desktop/UI/USBPorts.kext")
	fmt.Println(d)
	d = DiskUsage("")
	fmt.Println(d)
}
