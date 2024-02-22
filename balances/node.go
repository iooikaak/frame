package balance

import (
	"fmt"
	"sync/atomic"
)

// Service 服务对象信息
type Node struct {
	Id     int64
	IP     string
	Port   int
	Weight int64
	Tags   []string

	addr       string
	requestNum int64
}

// Addr 获取地址
func (n *Node) Addr() string {
	if n.addr == "" {
		n.addr = fmt.Sprintf("%s:%d", n.IP, n.Port)
	}

	return n.addr
}

//获取请求数量
func (n *Node) GetRequestNum() int64 {
	return n.requestNum
}

//增加请求数量
func (n *Node) AddRequestNum(num int64) {
	atomic.AddInt64(&n.requestNum, num)
}
