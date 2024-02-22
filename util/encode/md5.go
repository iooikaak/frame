package encode

import (
	"crypto/md5"
	"fmt"
)

// MD5 md5
func MD5(data []byte) string {
	h := md5.New()
	h.Write(data)
	return fmt.Sprintf("%x", h.Sum(nil))
}
