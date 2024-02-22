package micro

import (
	"time"

	"github.com/micro/go-micro/v2/client/selector"
	"github.com/micro/go-micro/v2/registry"

	"math/rand"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}
func NewSelector(consulIP []string, option ...selector.Option) selector.Selector {
	option = append(option, selector.Registry(NewRegister(consulIP)))
	return selector.NewSelector(option...)
}

func Random(services []*registry.Service) selector.Next {
	nodes := make([]*registry.Node, 0, len(services))

	for _, service := range services {
		if _, ok := service.Metadata["tag"]; ok {
			nodes = append(nodes, service.Nodes...)
		}
	}

	return func() (*registry.Node, error) {
		if len(nodes) == 0 {
			return nil, selector.ErrNoneAvailable
		}

		i := rand.Int() % len(nodes)
		return nodes[i], nil
	}

}
