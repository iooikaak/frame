package utils

import "testing"

func TestUUIDv4(t *testing.T) {
	for i := 0; i < 10; i++ {
		r := UUIDv4()
		t.Log(r, len(r))
	}
}
