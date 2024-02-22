package util

import (
	"crypto/md5"
	"hash/crc32"
	"io"
	"os"
	"strconv"
)

func GetWorkerID() int64 {

	hostname, err := os.Hostname()
	if err != nil {
		return -1
	}

	h := md5.New()
	io.WriteString(h, hostname)
	io.WriteString(h, strconv.Itoa(os.Getpid()))
	return int64(crc32.ChecksumIEEE(h.Sum(nil)) % 1024)

}
