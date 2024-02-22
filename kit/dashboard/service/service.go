package service

// Service 服务
type Service struct {
	Info *Info
}

// New New
func New() *Service {
	s := &Service{}
	s.Info = NewInfo()
	return s
}
