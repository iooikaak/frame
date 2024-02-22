package micro

import (
	"github.com/iooikaak/frame/micro/internal/consul"
	"github.com/micro/go-micro/v2/registry"
)

//consul from http
func NewRegister(addr []string) registry.Registry {
	return consul.NewRegistryWithHttp(registry.Addrs(addr...))
}
