package balance

import (
	"context"
)

type Balance interface {
	Get(context.Context) *Node
	Add(*Node)
	Remove(*Node)
}

type NodeList []*Node

func (s NodeList) Len() int {
	return len(s)
}

func (s NodeList) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s NodeList) Less(i, j int) bool {
	return s[i].GetRequestNum() < s[j].GetRequestNum()
}

func NewBalance(nodeList NodeList, bt string) Balance {
	switch bt {
	case "random":
		return NewRandom(nodeList)
	case "average":
		return NewRandom(nodeList)
	case "weight":
		return NewWeight(nodeList)
	}
	return nil
}
