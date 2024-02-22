package eureka

import (
	"sync"

	"github.com/iooikaak/frame/xlog"

	"github.com/iooikaak/frame/eureka"
	"github.com/micro/go-micro/v2/registry"
)

type eurekaRegistry struct {
	opts     registry.Options
	client   *eureka.Client
	config   *eureka.Config
	register bool
	sync.RWMutex
}

type Action struct {
	ActionType string
	AppName    string
}

func configure(c *eurekaRegistry, opts ...registry.Option) {
	for _, o := range opts {
		o(&c.opts)
	}
}

func (c *eurekaRegistry) Init(opts ...registry.Option) error {
	configure(c, opts...)
	return nil
}

func (c *eurekaRegistry) Register(s *registry.Service, opts ...registry.RegisterOption) error {
	c.RLock()
	register := c.register
	c.RUnlock()

	if !register {
		if c.client.Config != nil && c.client.Config.Instance != nil {
			if c.client.Config.Instance.Metadata == nil {
				c.client.Config.Instance.Metadata = make(map[string]string)
			}
			c.client.Config.Instance.Metadata["endpoints"] = encodeEndpoints(s.Endpoints)
			c.client.Config.Instance.Metadata["version"] = s.Version
			c.client.Config.Instance.Metadata["nodes"] = encodeNodes(s.Nodes)
		}
		err := c.client.Register()
		if err != nil {
			return err
		}
	} else {
		err := c.client.Heartbeat()
		if err != nil {
			return err
		}
		return nil
	}

	c.Lock()
	c.register = true
	c.Unlock()
	return nil
}

func (c *eurekaRegistry) Deregister(s *registry.Service, r ...registry.DeregisterOption) error {
	xlog.Debugf("Deregister start serviceName: %v", s.Name)
	err := c.client.UnRegister()
	if err != nil {
		return err
	}
	xlog.Debugf("Deregister successful serviceName: %v", s.Name)
	return nil
}

func (c *eurekaRegistry) GetService(name string, opts ...registry.GetOption) ([]*registry.Service, error) {
	l, err := c.client.GetService(name)
	if err != nil {
		return nil, err
	}
	serviceMap := make(map[string]*registry.Service)
	var service *registry.Service
	var ok bool
	for _, s := range l {
		metadata := s.Metadata
		version := metadata["version"]
		nodes := decodeNodes(metadata["nodes"])
		endpoints := decodeEndpoints(metadata["endpoints"])
		delete(metadata, "version")
		delete(metadata, "endpoints")
		delete(metadata, "nodes")
		if service, ok = serviceMap[version]; ok {
			service.Nodes = append(service.Nodes, nodes...)
			continue
		}

		serviceMap[version] = &registry.Service{
			Name:      name,
			Version:   version,
			Metadata:  metadata,
			Endpoints: endpoints,
			Nodes:     nodes,
		}
	}
	services := make([]*registry.Service, 0)
	for _, i := range serviceMap {
		services = append(services, i)
	}
	return services, nil
}

func (c *eurekaRegistry) ListServices(l ...registry.ListOption) ([]*registry.Service, error) {
	m, err := c.client.GetServices()
	if err != nil {
		return nil, err
	}
	serviceMap := make(map[string]*registry.Service)
	var service *registry.Service
	var ok bool
	for k, l := range m {
		appName := KeyNamed(k)
		for _, s := range l {
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
	xlog.Debugf("Start ListServices")
	services := make([]*registry.Service, 0)
	for _, i := range serviceMap {
		services = append(services, i)
	}
	return services, nil
}

func (c *eurekaRegistry) Watch(opts ...registry.WatchOption) (registry.Watcher, error) {
	return newEurekaWatcher(c, opts...)
}

func (c *eurekaRegistry) String() string {
	return "eureka"
}

func (c *eurekaRegistry) Options() registry.Options {
	return c.opts
}

func (c *eurekaRegistry) Client() *eureka.Client {
	if c.client != nil {
		return c.client
	}
	return nil
}

func NewRegistry(config *eureka.Config) registry.Registry {
	cr := &eurekaRegistry{
		config: config,
	}

	client := eureka.NewClient(config)
	cr.client = client
	configure(cr, registry.Addrs(config.DefaultZone...))
	cr.Lock()
	cr.client.Running = true
	cr.Unlock()
	if config.ServerOnly {
		go cr.Client().Refresh()
		go cr.Client().RollDiscoveries()
	}
	return cr
}
