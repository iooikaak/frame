package balance

import (
	"context"
	"testing"
)

var (
	mockNodeList = NodeList{
		&Node{
			requestNum: 0,
			IP:         "127.0.0.1",
			Port:       8080,
			Weight:     10,
			Id:         1,
		},
		&Node{
			requestNum: 0,
			IP:         "127.0.0.1",
			Port:       8081,
			Weight:     10,
			Id:         2,
		},
		&Node{
			requestNum: 0,
			IP:         "127.0.0.1",
			Port:       8082,
			Weight:     10,
			Id:         3,
		},
	}
)

func TestNewBalance(t *testing.T) {
	b := NewBalance(mockNodeList, "average")
	for i := 0; i < 100; i++ {
		n := b.Get(context.Background())
		t.Logf("--%+v--", n)
	}
}
