package eureka

import (
	"container/ring"
	"strings"
	"time"

	"github.com/iooikaak/frame/utils"

	"github.com/iooikaak/frame/xlog"
	"github.com/micro/go-micro/v2/registry"
)

const (
	InstanceStatusUp = "UP"
)

type eurekaWatcher struct {
	r            *eurekaRegistry
	opts         []registry.WatchOption
	next         chan *registry.Result
	exit         chan bool
	InstancesMap *Map
}

func newEurekaWatcher(cr *eurekaRegistry, opts ...registry.WatchOption) (registry.Watcher, error) {
	ew := &eurekaWatcher{
		r:    cr,
		opts: opts,
	}
	go func() { ew.refresh() }()
	go func() { cr.client.RollDiscoveries() }()
	return ew, nil
}

func (ew *eurekaWatcher) Next() (*registry.Result, error) {
	select {
	case <-ew.exit:
		return nil, registry.ErrWatcherStopped
	case r, ok := <-ew.next:
		if !ok {
			return nil, registry.ErrWatcherStopped
		}
		xlog.Debugf("eureka watcher Next serviceName: %s action: %s", r.Service.Name, r.Action)
		return r, nil
	}
}

func (ew *eurekaWatcher) Stop() {
	select {
	case <-ew.exit:
		return
	default:
		close(ew.exit)
		// drain results
		for {
			select {
			case <-ew.next:
			default:
				return
			}
		}
	}
}

// nolint
func (ew *eurekaWatcher) refresh() {
	ticker := time.NewTicker(time.Duration(ew.r.config.RegistryFetchIntervalSeconds) * time.Second)
	if ew.next == nil {
		ew.next = make(chan *registry.Result, 10)
	}
	if ew.exit == nil {
		ew.exit = make(chan bool)
	}
	for {
		select {
		case <-ticker.C:
			if err := ew.doRefresh(); err != nil {
				xlog.Errorf("eureka watcher refresh failed err: %v", err)
			}
			xlog.Debugf("eureka watcher refresh successful")
		}
	}
}

func (ew *eurekaWatcher) doRefresh() error {
	var (
		err error
	)
	services, err := ew.GetServices()
	if err != nil {
		return err
	}
	m := make(map[string]*Value)
	var changeFlag, ok bool
	for _, s := range services {
		if _, ok = m[s.Name]; !ok {
			m[s.Name] = NewValue([]*registry.Service{s})
			continue
		}
		m[s.Name].Val = append(m[s.Name].Val, s)
	}

	if ew.InstancesMap == nil {
		ew.InstancesMap = &Map{}
	}

	values := ew.InstancesMap.Load()
	var value *Value
	for k, v := range m {
		if value, ok = values[k]; !ok {
			for _, s := range v.Val {
				ew.next <- &registry.Result{
					Action:  "create",
					Service: s,
				}
			}
			changeFlag = true
			continue
		}
		if encodeValue(value) != encodeValue(v) {
			ew.CompareValue(value, v)
			changeFlag = true
		}
		delete(values, k)
	}

	if len(values) > 0 {
		for _, v := range values {
			for _, s := range v.Val {
				ew.next <- &registry.Result{
					Action:  "update",
					Service: s,
				}
			}
		}
		changeFlag = true
	}

	if changeFlag {
		ew.InstancesMap.Store(m)
	}

	return nil
}

func (ew *eurekaWatcher) GetServices() ([]*registry.Service, error) {
	m, err := ew.r.client.GetServices()
	if err != nil {
		return nil, err
	}
	serviceMap := make(map[string]*registry.Service)
	var service *registry.Service
	var ok bool
	var discoveryUrls []string
	for k, l := range m {
		appName := KeyNamed(k)
		for _, s := range l {
			if s.Status != InstanceStatusUp {
				continue
			}
			s.HomePageURL = strings.Trim(s.HomePageURL, "/")
			if appName == ew.r.config.DataCenterName {
				discoveryUrls = append(discoveryUrls, s.HomePageURL)
			}

			metadata := s.Metadata
			version := metadata["version"]
			nodes := decodeNodes(metadata["nodes"])
			endpoints := decodeEndpoints(metadata["endpoints"])
			delete(metadata, "version")
			delete(metadata, "endpoints")
			delete(metadata, "nodes")
			if service, ok = serviceMap[appName+version]; ok {
				serviceMap[appName+version].Nodes = append(service.Nodes, nodes...)
				continue
			}

			serviceMap[appName+version] = &registry.Service{
				Name:      appName,
				Version:   version,
				Metadata:  metadata,
				Endpoints: endpoints,
				Nodes:     nodes,
			}
		}
	}
	services := make([]*registry.Service, 0)
	for _, i := range serviceMap {
		services = append(services, i)
	}
	err = ew.StoreDiscoveries(discoveryUrls)
	if err != nil {
		return services, err
	}
	return services, nil
}

func (ew *eurekaWatcher) StoreDiscoveries(discoveryUrls []string) (err error) {
	var flag bool
	oldRing := ew.r.client.Discoveries
	if oldRing != nil && oldRing.Len() == len(discoveryUrls) {
		for i := 0; i < oldRing.Len(); i++ {
			if !utils.CheckIsExistString(discoveryUrls, oldRing.Value.(string)) {
				flag = true
				break
			}
			oldRing = oldRing.Next()
		}
	} else {
		flag = true
	}
	if flag {
		r := ring.New(len(discoveryUrls))
		for _, url := range discoveryUrls {
			r.Value = url
			r = r.Next()
		}
		err = ew.r.client.SetDiscoveries(r)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ew *eurekaWatcher) CompareValue(value, v *Value) {
	var ok, add bool

	mValue := make(map[string]*registry.Service)
	for _, v := range value.Val {
		mValue[v.Version] = v
	}

	var mService *registry.Service
	for _, s := range v.Val {
		if mService, ok = mValue[s.Version]; ok {
			mNode := make(map[string]*registry.Node)
			for _, mN := range mService.Nodes {
				mNode[mN.Id] = mN
			}
			for _, node := range s.Nodes {
				if _, ok = mNode[node.Id]; !ok {
					add = true
					continue
				}
				delete(mNode, node.Id)
			}
			if len(mNode) > 0 {
				nodes := make([]*registry.Node, 0, len(mNode))
				for _, n := range mNode {
					nodes = append(nodes, n)
				}
				delService := &registry.Service{
					Name:      s.Name,
					Version:   s.Version,
					Metadata:  s.Metadata,
					Endpoints: s.Endpoints,
					Nodes:     nodes,
				}
				ew.next <- &registry.Result{
					Action:  "delete",
					Service: delService,
				}
			}
			if add {
				ew.next <- &registry.Result{
					Action:  "update",
					Service: s,
				}
			}
			continue
		}
		ew.next <- &registry.Result{
			Action:  "update",
			Service: s,
		}
	}

}
