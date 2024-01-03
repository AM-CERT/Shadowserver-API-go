//go:build !windows

package shadowserver

import (
	"github.com/AM-CERT/Shadowserver-API-go/model"
	"syscall"
)

// DiskUsage disk usage of path/disk
func DiskUsage(path string) (disk model.DiskStatus, err error) {
	fs := syscall.Statfs_t{}
	err = syscall.Statfs(path, &fs)
	if err != nil {
		return
	}
	disk.All = fs.Blocks * uint64(fs.Bsize)
	disk.Free = fs.Bfree * uint64(fs.Bsize)
	disk.Used = disk.All - disk.Free
	return
}
