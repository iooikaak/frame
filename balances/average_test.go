package balance

import (
	"context"
	"fmt"
	"testing"
)

func TestNewAverage(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	average := NewAverage(mockNodeList)

	for i := 0; i < 50; i++ {
		t.Run(fmt.Sprintf("TestNewRandom：%d", i), func(t *testing.T) {
			average.Get(ctx)
		})
	}

	for i := 0; i < len(mockNodeList); i++ {
		t.Logf("mockNodeList i:%d RequestNum:%d", i, mockNodeList[i].GetRequestNum())
	}
}

func BenchmarkNewAverage(b *testing.B) {
	ctx := context.Background()
	b.ResetTimer()
	average := NewAverage(mockNodeList)

	for i := 0; i < b.N; i++ {
		average.Get(ctx)
	}
}
