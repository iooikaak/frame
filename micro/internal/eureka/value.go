package eureka

import (
	"github.com/micro/go-micro/v2/registry"
)

type Value struct {
	Val []*registry.Service
}

func NewValue(val []*registry.Service) *Value {
	value := &Value{
		Val: val,
	}
	return value
}
