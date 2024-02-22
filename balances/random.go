package balance

import (
	"context"
	"sync"
	"time"
)

//随机
type Random struct {
	sync.Mutex
	nodeList NodeList
}

func (r *Random) Get(ctx context.Context) *Node {
	defer report(ctx, "Random", time.Now())

	l := len(r.nodeList)

	u := time.Now().UnixNano()
	i := int(u) % l
	if i >= l {
		i = 0
	}
	node := r.nodeList[i]
	node.AddRequestNum(1)
	return node
}

func (r *Random) Add(node *Node) {
	r.Lock()
	defer r.Unlock()

	for _, v := range r.nodeList {
		if v.Id == node.Id {
			return
		}
	}
	r.nodeList = append(r.nodeList, node)
}

func (r *Random) Remove(node *Node) {
	r.Lock()
	defer r.Unlock()
	for k, v := range r.nodeList {
		if v.Id == node.Id {
			r.nodeList = append(r.nodeList[0:k], r.nodeList[k+1:]...)
		}
	}
}

func NewRandom(nodeList NodeList) *Random {
	return &Random{
		nodeList: nodeList,
	}
}
