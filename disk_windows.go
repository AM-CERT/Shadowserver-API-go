//go:build windows
// +build windows

package shadowserver

import (
	"github.com/AM-CERT/Shadowserver-API-go/model"
	"golang.org/x/sys/windows"
)

// DiskUsage disk usage of path/disk
func DiskUsage(path string) (disk model.DiskStatus, err error) {
	var freeBytesAvailable uint64
	var totalNumberOfBytes uint64
	var totalNumberOfFreeBytes uint64
	err = windows.GetDiskFreeSpaceEx(windows.StringToUTF16Ptr(path), &freeBytesAvailable, &totalNumberOfBytes, &totalNumberOfFreeBytes)
	if err != nil {
		return
	}
	disk.All = totalNumberOfBytes
	disk.Free = totalNumberOfFreeBytes
	disk.Used = disk.All - disk.Free
	return
}
