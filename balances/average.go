package balance

import (
	"context"
	"sort"
	"sync"
	"time"
)

//平均
type Average struct {
	sync.Mutex
	nodeList NodeList
}

func (a *Average) Get(ctx context.Context) *Node {
	defer report(ctx, "Average", time.Now())

	node := a.nodeList[0]
	node.AddRequestNum(1)
	return node
}

func (a *Average) Add(node *Node) {
	a.Lock()
	defer a.Unlock()
	for _, v := range a.nodeList {
		if v.Id == node.Id {
			return
		}
	}
	a.nodeList = append(a.nodeList, node)
	sort.Sort(a.nodeList)
}

func (a *Average) Remove(node *Node) {
	a.Lock()
	defer a.Unlock()
	for k, v := range a.nodeList {
		if v.Id == node.Id {
			a.nodeList = append(a.nodeList[0:k], a.nodeList[k+1:]...)
		}
	}
	sort.Sort(a.nodeList)
}

func NewAverage(nodeList NodeList) *Average {
	m := make(map[int64]*Node)
	var l NodeList
	for _, v := range nodeList {
		if _, ok := m[v.Id]; ok {
			continue
		}
		m[v.Id] = v
		l = append(l, v)
	}
	sort.Sort(l)
	return &Average{
		nodeList: l,
	}
}
