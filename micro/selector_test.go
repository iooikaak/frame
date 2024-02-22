package micro

import (
	"testing"

	"github.com/micro/go-micro/v2/client/selector"
)

func TestSelector(t *testing.T) {
	s := NewSelector(_consulAddr, selector.SetStrategy(Random))

	next, err := s.Select("mf.merchant")
	if err != nil {
		t.Fatal(err)
	}

	node, err := next()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(node)
}
