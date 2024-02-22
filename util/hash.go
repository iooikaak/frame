package util

import (
	"crypto/md5"
	"crypto/sha1"
	"fmt"
)

func Sha1Sum(input []byte) string {
	h := sha1.New()
	h.Write(input)
	bs := h.Sum(nil)
	return fmt.Sprintf("%x", bs)
}

func MD5Sum(input []byte) string {
	return fmt.Sprintf("%x", md5.Sum(input))
}

func XOREncryptDecrypt(input, key string) (output string) {
	for i := range input {
		output += string(input[i] ^ key[i%len(key)])
	}
	return output
}
