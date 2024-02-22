package eureka

import (
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/iooikaak/frame/utils"

	"container/ring"

	"github.com/iooikaak/frame/net/ip"
	"github.com/iooikaak/frame/xlog"
)

const (
	InstanceStatusUp = "UP"
)

type Client struct {
	signalChan         chan os.Signal
	mutex              sync.RWMutex
	Running            bool
	Config             *Config
	DefaultDiscoveries *ring.Ring
	Discoveries        *ring.Ring
	InstancesMap       *Map
}

type Action struct {
	ActionType string
	AppName    string
}

func (c *Client) Start() {
	c.mutex.Lock()
	c.Running = true
	c.InstancesMap = &Map{}
	c.mutex.Unlock()
	if err := c.Register(); err != nil {
		xlog.Errorf("eureka Register failed err: %v", err)
		return
	}
	go c.Refresh()
	go c.heartbeat()
	go c.handleSignal()
	go c.RollDiscoveries()
}

func (c *Client) SetDiscoveries(discoveries *ring.Ring) error {
	c.mutex.Lock()
	c.Discoveries = discoveries
	c.mutex.Unlock()
	return nil
}

func (c *Client) GetService(serviceName string) (l []Instance, err error) {
	var (
		application *Application
	)
	c.mutex.RLock()
	r := c.getCenterUrl()
	var url string
	for i := 0; i < r.Len(); i++ {
		url = r.Value.(string)
		application, err = GetApplication(url, serviceName)
		if err == nil {
			xlog.Debugf("GetService serviceName: %v url: %v", serviceName, url)
			break
		}
		r = r.Next()
	}
	c.mutex.RUnlock()
	if err != nil {
		return nil, err
	}
	if application == nil {
		return nil, err
	}
	return application.Instances, nil
}

func (c *Client) GetServices() (m map[string][]Instance, err error) {
	var (
		applications *Applications
	)
	m = make(map[string][]Instance)
	c.mutex.RLock()
	r := c.getCenterUrl()
	var url string
	for i := 0; i < r.Len(); i++ {
		url = r.Value.(string)
		applications, err = GetApplications(url)
		if err == nil {
			xlog.Debugf("GetServices serviceName: %v url: %v", c.Config.App, url)
			break
		}
		r = r.Next()
	}
	c.mutex.RUnlock()
	if err != nil {
		return m, err
	}

	for _, v := range applications.Applications {
		m[v.Name] = v.Instances
	}
	return m, nil
}

func (c *Client) GetServiceUrl(serviceName string) (l []string, err error) {
	var value *Value
	if c.InstancesMap != nil {
		value = c.InstancesMap.Get(KeyNamed(serviceName))

	}
	if value != nil {
		instances := value.Val
		for _, instance := range instances {
			l = append(l, instance.HomePageURL)
		}
	} else {
		instances, _ := c.GetService(serviceName)
		for _, instance := range instances {
			l = append(l, instance.HomePageURL)
		}
	}

	return l, nil
}

func (c *Client) Register() (err error) {
	instance := c.Config.Instance
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	r := c.getCenterUrl()
	var url string
	for i := 0; i < r.Len(); i++ {
		if !c.Running {
			err = nil
			break
		}
		url = r.Value.(string)
		err = Register(url, c.Config.App, instance)
		if err == nil {
			xlog.Debugf("Register serviceName: %v url: %v", c.Config.App, url)
			break
		}
		r = r.Next()
	}
	if err != nil {
		return
	}
	return nil
}

func (c *Client) UnRegister() (err error) {
	instance := c.Config.Instance
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	r := c.getCenterUrl()
	var url string
	for i := 0; i < r.Len(); i++ {
		url = r.Value.(string)
		err = UnRegister(url, instance.App, instance.InstanceID)
		if err == nil {
			xlog.Debugf("UnRegister serviceName: %v url: %v", c.Config.App, url)
			break
		}
		r = r.Next()
	}
	return err
}

func (c *Client) Heartbeat() error {
	instance := c.Config.Instance
	var (
		err error
	)
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	r := c.getCenterUrl()
	var url string
	for i := 0; i < r.Len(); i++ {
		url = r.Value.(string)
		err = Heartbeat(url, instance.App, instance.InstanceID)
		if err == nil {
			xlog.Debugf("Heartbeat serviceName: %v url: %v", c.Config.App, url)
			break
		}
		r = r.Next()
	}
	return err
}

func (c *Client) doRefresh() error {
	var (
		applications *Applications
		err          error
	)
	c.mutex.RLock()
	r := c.getCenterUrl()
	var url string
	for i := 0; i < r.Len(); i++ {
		url = r.Value.(string)
		applications, err = GetApplications(url)
		if err == nil {
			break
		}
		r = r.Next()
	}
	c.mutex.RUnlock()
	if err != nil {
		return err
	}

	var checkFlag bool
	if applications == nil {
		return nil
	}
	if c.InstancesMap == nil {
		c.InstancesMap = &Map{}
	}
	m := make(map[string]*Value)
	discoveryUrls := make([]string, 0)
	oldValues := c.InstancesMap.Load()
	for _, app := range applications.Applications {
		appName := KeyNamed(app.Name)
		instances := make([]*Instance, 0)
		for k, instance := range app.Instances {
			if !strings.EqualFold(instance.Status, InstanceStatusUp) {
				continue
			}
			if strings.EqualFold(app.Name, c.Config.DataCenterName) {
				discoveryUrls = append(discoveryUrls, strings.Trim(instance.HomePageURL, "/"))
			}
			app.Instances[k].HomePageURL = strings.Trim(instance.HomePageURL, "/")
			instances = append(instances, &app.Instances[k])
		}

		value := NewValue(instances)
		m[appName] = value
		oldValue := oldValues[appName]
		if oldValue == nil || oldValue.Md5 != value.Md5 {
			checkFlag = true
		}
		delete(oldValues, appName)
	}
	if len(oldValues) > 0 {
		checkFlag = true
	}

	if checkFlag {
		c.InstancesMap.Store(m)
	}
	var flag bool
	if c.Discoveries != nil && c.Discoveries.Len() == len(discoveryUrls) {
		for i := 0; i < c.Discoveries.Len(); i++ {
			if !utils.CheckIsExistString(discoveryUrls, c.Discoveries.Value.(string)) {
				flag = true
				break
			}
			c.Discoveries = c.Discoveries.Next()
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
		c.mutex.Lock()
		c.Discoveries = r
		c.mutex.Unlock()
	}
	return nil
}

// nolint
func (c *Client) heartbeat() {
	ticker := time.NewTicker(time.Duration(c.Config.RenewalIntervalInSecs) * time.Second)
	for {
		select {
		case <-ticker.C:
			if c.Running {
				if err := c.Heartbeat(); err != nil {
					if !c.Running {
						break
					}
					if err == ErrNotFound {
						xlog.Warnf("eureka heartbeat Not Found need register serviceName: %s", c.Config.App)
						if err = c.Register(); err != nil {
							xlog.Errorf("eureka Register serviceName: %s err: %v", c.Config.App, err)
						}
						break
					}
					xlog.Errorf("eureka heartbeat serviceName: %s err: %v", c.Config.App, err)
					break
				}
				xlog.Debugf("eureka heartbeat successful serviceName: %s", c.Config.App)
			} else {
				break
			}
		}
	}
}

// nolint
func (c *Client) Refresh() chan string {
	ticker := time.NewTicker(time.Duration(c.Config.RegistryFetchIntervalSeconds) * time.Second)
	for {
		select {
		case <-ticker.C:
			if err := c.doRefresh(); err != nil {
				xlog.Errorf("eureka Refresh failed serviceName: %s err: %v", c.Config.App, err)
				break
			}
			xlog.Debugf("eureka Refresh successful serviceName: %s", c.Config.App)
		}
	}
}

// nolint
func (c *Client) RollDiscoveries() {
	ticker := time.NewTicker(time.Duration(c.Config.RollDiscoveriesIntervalSeconds) * time.Second)
	for {
		select {
		case <-ticker.C:
			if c.Discoveries != nil {
				c.mutex.Lock()
				c.Discoveries = c.Discoveries.Next()
				c.mutex.Unlock()
			}
		}
	}
}

func (c *Client) handleSignal() {
	if c.signalChan == nil {
		c.signalChan = make(chan os.Signal)
	}
	signal.Notify(c.signalChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	<-c.signalChan
	xlog.Warnf("eureka receive exit signal client Instance going to UnRegister serviceName: %s", c.Config.App)
	c.mutex.Lock()
	c.Running = false
	c.mutex.Unlock()
	err := c.UnRegister()
	if err != nil {
		xlog.Errorf("eureka UnRegister failed serviceName: %s err: %v", c.Config.App, err)
	} else {
		xlog.Infof("eureka UnRegister application Instance successful serviceName: %s", c.Config.App)
	}
}

func (c *Client) getCenterUrl() *ring.Ring {
	if c.Discoveries != nil && c.Discoveries.Len() > 0 {
		return c.Discoveries
	}
	return c.DefaultDiscoveries
}

func NewClient(config *Config) *Client {
	defaultConfig(config)
	config.Instance = NewInstance(ip.InternalIP(), config)
	r := ring.New(len(config.DefaultZone))
	for _, v := range config.DefaultZone {
		if strings.Contains(v, "http://") || strings.Contains(v, "https://") {
			r.Value = v
		} else {
			r.Value = "http://" + v
		}
		r = r.Next()
	}
	return &Client{
		Config:             config,
		DefaultDiscoveries: r,
	}
}

func defaultConfig(config *Config) {
	if len(config.DefaultZone) == 0 {
		config.DefaultZone = []string{"http://localhost:8761"}
	}
	if config.RenewalIntervalInSecs <= 0 {
		config.RenewalIntervalInSecs = 30
	}
	if config.RegistryFetchIntervalSeconds <= 0 {
		config.RegistryFetchIntervalSeconds = 15
	}
	if config.DurationInSecs <= 0 {
		config.DurationInSecs = 90
	}
	if config.RollDiscoveriesIntervalSeconds <= 0 {
		config.RollDiscoveriesIntervalSeconds = 60
	}

	if len(strings.TrimSpace(config.App)) == 0 {
		config.App = "server"
	} else {
		config.App = strings.ToLower(config.App)
	}
	if config.Port <= 0 {
		config.Port = 80
	}
	if len(strings.TrimSpace(config.DataCenterName)) == 0 {
		config.DataCenterName = "discovery"
	}
}
