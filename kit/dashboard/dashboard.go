package dashboard

import "github.com/iooikaak/frame/kit/dashboard/service"

// DashBoard DashBoard
type DashBoard struct {
	Service *service.Service
}

// New New
func New() *DashBoard {
	dashBoard := &DashBoard{}
	dashBoard.Service = service.New()
	return dashBoard
}
