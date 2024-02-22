package balance

import (
	"context"
	"strconv"
	"time"
)

// 权重
type Weight struct {
	consistent *Consistent
}

func (w *Weight) Get(ctx context.Context) *Node {
	defer report(ctx, "Weight", time.Now())
	return w.consistent.Get(strconv.FormatInt(time.Now().UnixNano(), 10))
}

func (w *Weight) Add(node *Node) {
	w.consistent.Add(node)
}

func (w *Weight) Remove(node *Node) {
	w.consistent.Remove(node)
}

func NewWeight(nodeList NodeList) *Weight {
	return &Weight{
		consistent: NewConsistent(nodeList),
	}
}

type HashRing []uint32

func (c HashRing) Len() int {
	return len(c)
}

func (c HashRing) Less(i, j int) bool {
	return c[i] < c[j]
}

func (c HashRing) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
