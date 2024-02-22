package balance

import (
	"context"
	"fmt"
	"testing"
)

func TestNewRandom(t *testing.T) {
	ctx := context.Background()

	t.Parallel()

	random := NewRandom(mockNodeList)

	for i := 0; i < 50; i++ {
		t.Run(fmt.Sprintf("TestNewRandom：%d", i), func(t *testing.T) {
			random.Get(ctx)
		})
	}

	for i := 0; i < len(mockNodeList); i++ {
		t.Logf("mockNodeList i:%d RequestNum:%d", i, mockNodeList[i].GetRequestNum())
	}

}

func BenchmarkNewRandom(b *testing.B) {
	ctx := context.Background()

	b.ResetTimer()
	random := NewRandom(mockNodeList)
	for i := 0; i < b.N; i++ {
		random.Get(ctx)
	}
}
