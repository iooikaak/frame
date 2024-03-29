package eureka

import (
	"strings"
	"sync/atomic"
)

// KeyNamed key naming to lower case.
func KeyNamed(key string) string {
	return strings.ToLower(key)
}

// Map is config map, key(filename) -> value(file).
type Map struct {
	values atomic.Value
}

// Store sets the value of the Value to values map.
func (m *Map) Store(values map[string]*Value) {
	dst := make(map[string]*Value, len(values))
	for k, v := range values {
		dst[KeyNamed(k)] = v
	}
	m.values.Store(dst)
}

// Load returns the value set by the most recent Store.
func (m *Map) Load() map[string]*Value {
	l := m.values.Load()
	if l == nil {
		return nil
	}
	src := l.(map[string]*Value)
	dst := make(map[string]*Value, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

// Exist check if values map exist a key.
func (m *Map) Exist(key string) bool {
	values := m.Load()
	if values == nil {
		return false
	}
	_, ok := values[KeyNamed(key)]
	return ok
}

// Get return get value by key.
func (m *Map) Get(key string) *Value {
	values := m.Load()
	if values == nil {
		return nil
	}
	if v, ok := values[KeyNamed(key)]; ok {
		return v
	}
	return &Value{}
}

// Keys return map keys.
func (m *Map) Keys() []string {
	values := m.Load()
	if values == nil {
		return nil
	}
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	return keys
}
