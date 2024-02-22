package snowflake

import (
	"github.com/iooikaak/frame/snowflake/sony"
)

func init() {
	_ = _snowFlake.loadAdapter()
}

// ISnowFlake 接口
type ISnowFlake interface {
	NewID() (uint64, error)
	NewId() (int64, error)
}

// NewID NewID
func NewID() (uint64, error) {
	return _snowFlake.adapter.NewID()
}

// NewID NewID
func NewId() (int64, error) {
	id, err := _snowFlake.adapter.NewID()
	return int64(id), err
}

// SetAdapter 设置Adapter
func SetAdapter(adaptertype AdapterType) error {
	return _snowFlake.setAdapter(adaptertype)
}

var _snowFlake snowFlake

type snowFlake struct {
	adapterType AdapterType
	adapter     ISnowFlake
}

func (s *snowFlake) setAdapter(adapterType AdapterType) error {
	s.adapterType = adapterType
	return s.loadAdapter()
}

func (s *snowFlake) loadAdapter() error {
	if s.adapterType == "" {
		s.adapterType = SONY
	}

	var err error

	switch s.adapterType {
	case SONY:
		s.adapter, err = sony.New()
		if err != nil {
			return err
		}
	}

	return nil
}
