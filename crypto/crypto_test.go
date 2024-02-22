package crypto

import "testing"

func TestMd5(t *testing.T) {
	md5, _ := MD5("testtttttttt4444444")
	t.Log(md5)
}
