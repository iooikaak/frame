package sony

import (
	"errors"

	"github.com/sony/sonyflake"
)

// Implement 实现
type Implement struct {
	SonyFlake *sonyflake.Sonyflake
}

// NewID NewID
func (s *Implement) NewID() (uint64, error) {
	return s.SonyFlake.NextID()
}

// NewID NewID
func (s *Implement) NewId() (int64, error) {
	id, err := s.SonyFlake.NextID()
	return int64(id), err
}

// New 实例
func New() (*Implement, error) {
	var st sonyflake.Settings
	//st.MachineID = awsutil.AmazonEC2MachineID
	sf := sonyflake.NewSonyflake(st)

	if sf == nil {
		return nil, errors.New("创建SonyFlake对象失败")
	}

	impl := &Implement{}
	impl.SonyFlake = sf

	return impl, nil
}
