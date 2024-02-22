package balance

import (
	"context"
	"fmt"
	"testing"
)

func TestNewWeight(t *testing.T) {
	ctx := context.Background()

	t.Parallel()

	weight := NewWeight(mockNodeList)

	for i := 0; i < 50; i++ {
		t.Run(fmt.Sprintf("TestNewWeight：%d", i), func(t *testing.T) {
			weight.Get(ctx)
		})
	}

	for i := 0; i < len(mockNodeList); i++ {
		t.Logf("mockNodeList i:%d RequestNum:%d", i, mockNodeList[i].GetRequestNum())
	}

}

func BenchmarkNewWeight(b *testing.B) {
	ctx := context.Background()
	weight := NewWeight(mockNodeList)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		weight.Get(ctx)
	}
}
