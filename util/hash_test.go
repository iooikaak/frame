package util

import (
	b64 "encoding/base64"
	"fmt"
	"testing"
)

func TestXOREncryptDecrypt(t *testing.T) {
	s := "31.172211371527776|121.410971408420139|021|上海市|310104|闵行区"
	k := "1234567890"
	s = b64.URLEncoding.EncodeToString([]byte(s))
	a := XOREncryptDecrypt(s, k)
	b := b64.StdEncoding.EncodeToString([]byte(a))
	fmt.Println(b)
	c, _ := b64.StdEncoding.DecodeString(b)
	d := XOREncryptDecrypt(string(c), k)
	bs, err := b64.URLEncoding.DecodeString(d)
	fmt.Println(string(bs), err)
}
