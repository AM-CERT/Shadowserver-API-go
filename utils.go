package shadowserver

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"shadowserver/model"
	"syscall"
)

const (
	B  = 1
	KB = 1024 * B
	MB = 1024 * KB
	GB = 1024 * MB
)

// ComputeHmac compute HMAC of data
func ComputeHmac(secret string, data []byte) string {
	// create a new HMAC by defining the hash type and the key
	h := hmac.New(sha256.New, []byte(secret))

	// compute the HMAC
	h.Write(data)
	d := h.Sum(nil)

	return hex.EncodeToString(d)
}

// DiskUsage disk usage of path/disk
func DiskUsage(path string) (disk model.DiskStatus) {
	fs := syscall.Statfs_t{}
	err := syscall.Statfs(path, &fs)
	if err != nil {
		return
	}
	disk.All = fs.Blocks * uint64(fs.Bsize)
	disk.Free = fs.Bfree * uint64(fs.Bsize)
	disk.Used = disk.All - disk.Free
	return
}

// PrintJson print json string to stdout
func PrintJson(data []byte, pretty bool) {
	var err error
	var jsonString []byte

	if pretty {
		jsonString, err = json.MarshalIndent(json.RawMessage(data), "", " ")
	} else {
		jsonString, err = json.Marshal(json.RawMessage(data))
	}
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Fatal("failed to marshal json string")
	}

	fmt.Print(string(jsonString))
}

// FileExists check file exists
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
