package balance

import (
	"hash/crc32"
	"sort"
	"strconv"
	"sync"
)

//默认分片
const DefaultReplicas = 160

type Consistent struct {
	sync.RWMutex
	numReps   int64
	ring      HashRing
	Nodes     map[uint32]Node
	Resources map[int64]bool
}

func NewConsistent(svc NodeList) *Consistent {
	c := &Consistent{
		numReps:   DefaultReplicas,
		ring:      HashRing{},
		Nodes:     make(map[uint32]Node),
		Resources: make(map[int64]bool),
	}
	c.BatchAdd(svc)
	return c
}

func (c *Consistent) Add(node *Node) bool {
	c.Lock()
	defer c.Unlock()

	if _, ok := c.Resources[node.Id]; ok {
		return false
	}

	count := c.numReps * node.Weight
	for i := int64(0); i < count; i++ {
		str := c.joinStr(i, node)
		c.Nodes[c.hashStr(str)] = *(node)
	}
	c.Resources[node.Id] = true
	c.sortHashRing()
	return true
}

func (c *Consistent) BatchAdd(svc NodeList) bool {
	c.Lock()
	defer c.Unlock()
	svcListLen := len(svc)
	for i := 0; i < svcListLen; i++ {
		if _, ok := c.Resources[svc[i].Id]; ok {
			return false
		}

		count := c.numReps * svc[i].Weight
		for j := int64(0); j < count; j++ {
			str := c.joinStr(j, svc[i])
			c.Nodes[c.hashStr(str)] = *(svc[i])
		}
		c.Resources[svc[i].Id] = true
	}

	c.sortHashRing()
	return true
}

func (c *Consistent) Get(key string) *Node {
	if len(c.ring) == 0 {
		return nil
	}
	c.RLock()
	defer c.RUnlock()

	hash := c.hashStr(key)
	i := c.search(hash)
	node := c.Nodes[c.ring[i]]
	return &node
}

func (c *Consistent) joinStr(i int64, nodeWeight *Node) string {
	return nodeWeight.IP + "*" +
		strconv.FormatInt(nodeWeight.Weight, 10) +
		"-" + strconv.FormatInt(i, 10) +
		"-" + strconv.FormatInt(nodeWeight.Id, 10)
}

func (c *Consistent) sortHashRing() {
	c.ring = HashRing{}
	for k := range c.Nodes {
		c.ring = append(c.ring, k)
	}
	sort.Sort(c.ring)
}

// MurMurHash算法 :https://github.com/spaolacci/murmur3
func (c *Consistent) hashStr(key string) uint32 {
	return crc32.ChecksumIEEE([]byte(key))
}

func (c *Consistent) search(hash uint32) int {
	// 二分查找
	i := sort.Search(len(c.ring), func(i int) bool { return c.ring[i] >= hash })
	if i < len(c.ring) {
		if i == len(c.ring)-1 {
			return 0
		} else {
			return i
		}
	} else {
		return len(c.ring) - 1
	}
}

func (c *Consistent) Remove(node *Node) {
	c.Lock()
	defer c.Unlock()

	if _, ok := c.Resources[node.Id]; !ok {
		return
	}

	delete(c.Resources, node.Id)

	count := c.numReps * node.Weight
	for i := int64(0); i < count; i++ {
		str := c.joinStr(i, node)
		delete(c.Nodes, c.hashStr(str))
	}
	c.sortHashRing()
}
