package consul

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	consul "github.com/hashicorp/consul/api"
	"github.com/micro/go-micro/v2/registry"
	mnet "github.com/micro/go-micro/v2/util/net"
	hash "github.com/mitchellh/hashstructure"
)

const (
	intervalDefault = time.Second
	timeOutDefault  = time.Second * 3
	versionDefault  = "1.0.0"
	defaultTag      = "tag" // 兼容老版本没有tag
)

type consulRegistryHttp struct {
	Address []string
	opts    registry.Options

	client *consul.Client
	config *consul.Config

	// connect enabled
	connect bool

	queryOptions *consul.QueryOptions

	sync.Mutex
	register map[string]uint64
}

func config(c *consulRegistryHttp, opts ...registry.Option) {
	// set opts
	for _, o := range opts {
		o(&c.opts)
	}

	// use default config DefaultConfig
	config := consul.DefaultNonPooledConfig()

	if c.opts.Context != nil {
		// Use the consul config passed in the options, if available
		if co, ok := c.opts.Context.Value("consul_config").(*consul.Config); ok {
			config = co
		}
		if cn, ok := c.opts.Context.Value("consul_connect").(bool); ok {
			c.connect = cn
		}

		// Use the consul query options passed in the options, if available
		if qo, ok := c.opts.Context.Value("consul_query_options").(*consul.QueryOptions); ok && qo != nil {
			c.queryOptions = qo
		}
		if as, ok := c.opts.Context.Value("consul_allow_stale").(bool); ok {
			c.queryOptions.AllowStale = as
		}
	}

	// check if there are any addrs
	var addrs []string

	// iterate the options addresses
	for _, address := range c.opts.Addrs {
		// check we have a port
		addr, port, err := net.SplitHostPort(address)
		if ae, ok := err.(*net.AddrError); ok && ae.Err == "missing port in address" {
			port = "8500"
			addr = address
			addrs = append(addrs, net.JoinHostPort(addr, port))
		} else if err == nil {
			addrs = append(addrs, net.JoinHostPort(addr, port))
		}
	}

	// set the addrs
	if len(addrs) > 0 {
		c.Address = addrs
		config.Address = c.Address[0]
	}

	if config.HttpClient == nil {
		config.HttpClient = new(http.Client)
	}

	// requires secure connection?
	if c.opts.Secure || c.opts.TLSConfig != nil {
		config.Scheme = "https"
		// We're going to support InsecureSkipVerify
		config.HttpClient.Transport = newTransport(c.opts.TLSConfig)
	}

	// set timeout
	if c.opts.Timeout > 0 {
		config.HttpClient.Timeout = c.opts.Timeout
	}

	// set the config
	c.config = config

	// remove client
	c.client = nil

	// setup the client
	c.Client()
}

func (c *consulRegistryHttp) Init(opts ...registry.Option) error {
	config(c, opts...)
	return nil
}

func (c *consulRegistryHttp) Deregister(s *registry.Service, r ...registry.DeregisterOption) error {
	if len(s.Nodes) == 0 {
		return errors.New("Require at least one node")
	}

	// delete our hash and time check of the service
	c.Lock()
	delete(c.register, s.Name)
	c.Unlock()

	node := s.Nodes[0]
	return c.Client().Agent().ServiceDeregister(node.Id)
}

func (c *consulRegistryHttp) Register(s *registry.Service, opts ...registry.RegisterOption) error {
	if len(s.Nodes) == 0 {
		return errors.New("Require at least one node")
	}

	var (
		regInterval = intervalDefault
		timeOut     = timeOutDefault
		options     registry.RegisterOptions
		httpUrl     string
	)
	for _, o := range opts {
		o(&options)
	}

	if c.opts.Context != nil {
		if CheckInterval, ok := c.opts.Context.Value("consul_http_check").(time.Duration); ok {
			regInterval = CheckInterval
		}
		if timeOutD, ok := c.opts.Context.Value("consul_http_check_time_out").(time.Duration); ok {
			timeOut = timeOutD
		}
	}

	// create hash of service; uint64
	h, err := hash.Hash(s, nil)
	if err != nil {
		return err
	}

	// use first node
	node := s.Nodes[0]

	// get existing hash and last checked time
	c.Lock()
	v, ok := c.register[s.Name]
	c.Unlock()

	// if it's already registered and matches then just pass the check
	if ok && v == h {
		services, _, err := c.Client().Health().Checks(s.Name, c.queryOptions)
		if err == nil {
			for _, v := range services {
				if v.ServiceID == node.Id {
					return nil
				}
			}
		}
	}

	httpUrl = "http://" + node.Address + "/health"
	if strings.Contains(node.Address, "http") {
		u, err := url.Parse(node.Address)
		if err != nil {
			return errors.New("url parse error " + httpUrl)
		}
		httpUrl = node.Address
		node.Address = u.Host
	}

	// encode the tags
	tags := encodeMetadata(node.Metadata)
	tags = append(tags, encodeEndpoints(s.Endpoints)...)
	tags = append(tags, encodeVersion(s.Version)...)

	var check = &consul.AgentServiceCheck{
		HTTP:     httpUrl,
		Interval: fmt.Sprintf("%v", regInterval),
		Timeout:  fmt.Sprintf("%v", timeOut),
	}

	host, pt, _ := net.SplitHostPort(node.Address)
	if host == "" {
		host = node.Address
	}
	port, _ := strconv.Atoi(pt)

	// register the service
	asr := &consul.AgentServiceRegistration{
		ID:      node.Id,
		Name:    s.Name,
		Tags:    tags,
		Port:    port,
		Address: host,
		Check:   check,
	}

	// Specify consul connect
	if c.connect {
		asr.Connect = &consul.AgentServiceConnect{
			Native: true,
		}
	}

	if err := c.Client().Agent().ServiceRegister(asr); err != nil {
		return err
	}

	// save our hash and time check of the service
	c.Lock()
	c.register[s.Name] = h
	c.Unlock()

	return nil
}

func (c *consulRegistryHttp) fetchOldTag(tags []string) string {
	var tmpData = make([]string, 0)
	for _, tag := range tags {
		if len(tag) == 0 {
			continue
		}

		switch tag[0] {
		case 't', 'v', 'e':
			if tag[1] == '-' {
				continue
			}
		}

		tmpData = append(tmpData, tag)

	}

	return strings.Join(tmpData, ",")
}

func (c *consulRegistryHttp) GetService(name string, opts ...registry.GetOption) ([]*registry.Service, error) {
	var rsp []*consul.ServiceEntry
	var err error

	// if we're connect enabled only get connect services
	if c.connect {
		rsp, _, err = c.Client().Health().Connect(name, "", false, c.queryOptions)
	} else {
		rsp, _, err = c.Client().Health().Service(name, "", false, c.queryOptions)
	}
	if err != nil {
		return nil, err
	}

	serviceMap := map[string]*registry.Service{}

	for _, s := range rsp {
		if s.Service.Service != name {
			continue
		}

		// version is now a tag
		version, _ := decodeVersion(s.Service.Tags)
		// service ID is now the node id
		id := s.Service.ID
		// key is always the version
		key := version
		// address is service address
		address := s.Service.Address

		// use node address
		if len(address) == 0 {
			address = s.Node.Address
		}

		// TODO Compatible with non-micro registered atg
		if len(key) == 0 {
			key = versionDefault
		}
		svc, ok := serviceMap[key]
		if !ok {
			svc = &registry.Service{
				Endpoints: decodeEndpoints(s.Service.Tags),
				Name:      s.Service.Service,
				Version:   version,
			}
			serviceMap[key] = svc
		}

		var del bool

		for _, check := range s.Checks {
			// delete the node if the status is critical
			if check.Status == "critical" {
				del = true
				break
			}
		}

		// if delete then skip the node
		if del {
			continue
		}

		svc.Nodes = append(svc.Nodes, &registry.Node{
			Id:       id,
			Address:  mnet.HostPort(address, s.Service.Port),
			Metadata: decodeMetadata(s.Service.Tags),
		})
		// TODO Compatible with non-micro registered atg
		if t := c.fetchOldTag(s.Service.Tags); len(t) > 0 {
			svc.Nodes[len(svc.Nodes)-1].Metadata[defaultTag] = t
		}
	}

	var services []*registry.Service
	for _, service := range serviceMap {
		services = append(services, service)
	}
	return services, nil
}

func (c *consulRegistryHttp) ListServices(l ...registry.ListOption) ([]*registry.Service, error) {
	rsp, _, err := c.Client().Catalog().Services(c.queryOptions)
	if err != nil {
		return nil, err
	}

	var services []*registry.Service

	for service := range rsp {
		services = append(services, &registry.Service{Name: service})
	}

	return services, nil
}

func (c *consulRegistryHttp) Watch(opts ...registry.WatchOption) (registry.Watcher, error) {
	return newConsulWatcherHttp(c, opts...)
}

func (c *consulRegistryHttp) String() string {
	return "consul-http"
}

func (c *consulRegistryHttp) Options() registry.Options {
	return c.opts
}

func (c *consulRegistryHttp) Client() *consul.Client {
	if c.client != nil {
		return c.client
	}

	for _, addr := range c.Address {
		// set the address
		c.config.Address = addr

		// create a new client
		tmpClient, _ := consul.NewClient(c.config)

		// test the client
		_, err := tmpClient.Agent().Host()
		if err != nil {
			continue
		}

		// set the client
		c.client = tmpClient
		return c.client
	}

	// set the default
	c.client, _ = consul.NewClient(c.config)

	// return the client
	return c.client
}

func NewRegistryWithHttp(opts ...registry.Option) registry.Registry {
	cr := &consulRegistryHttp{
		opts:     registry.Options{},
		register: make(map[string]uint64),
		queryOptions: &consul.QueryOptions{
			AllowStale: true,
		},
	}
	config(cr, opts...)
	return cr
}
