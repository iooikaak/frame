package eureka

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	xHttp "github.com/iooikaak/frame/net/http/blademaster"
)

var (
	ErrNotFound = errors.New("not found")
	c           = xHttp.NewClient(&xHttp.ClientConfig{
		Dial:      5 * time.Second,
		Timeout:   5 * time.Second,
		KeepAlive: 100,
	})
	registerUrl        = "/eureka/apps/"
	unRegisterUrl      = "/eureka/apps/"
	getApplicationsUrl = "/eureka/apps/"
	getApplicationUrl  = "/eureka/apps/"
	heartbeatUrl       = "/eureka/apps/"
)

func Register(zone, app string, instance *Instance) error {
	type InstanceInfo struct {
		Instance *Instance `json:"instance"`
	}
	var info = &InstanceInfo{
		Instance: instance,
	}

	u := zone + registerUrl + app
	ctx := context.Background()

	var resp interface{}
	err := c.PostJson(ctx, u, "", info, &resp)
	if err != nil {
		return fmt.Errorf("register application Instance failed, error: %s", err)
	}
	return nil
}

func UnRegister(zone, app, instanceID string) error {
	u := zone + unRegisterUrl + app + "/" + instanceID
	ctx := context.Background()

	var resp interface{}
	err := c.Delete(ctx, u, "", nil, &resp)
	if err != nil {
		return fmt.Errorf("unRegister application Instance failed, u: %v error: %s", u, err)
	}
	return nil
}

func GetApplications(zone string) (*Applications, error) {
	type Result struct {
		Applications *Applications `json:"applications"`
	}
	apps := new(Applications)
	res := &Result{
		Applications: apps,
	}
	u := zone + getApplicationsUrl
	ctx := context.Background()
	params := url.Values{}
	req, err := c.NewRequest(http.MethodGet, u, "", params)
	if err != nil {
		return apps, err
	}
	req.Header.Set("Accept", "application/json")
	err = c.Do(ctx, req, res)
	if err != nil {
		return nil, fmt.Errorf("GetApplications failed, error: %s", err)
	}
	return apps, nil
}

func GetApplication(zone, appId string) (*Application, error) {
	type Result struct {
		Application *Application `json:"application"`
	}
	res := &Result{
		Application: &Application{},
	}
	u := zone + getApplicationUrl + appId
	ctx := context.Background()
	params := url.Values{}
	req, err := c.NewRequest(http.MethodGet, u, "", params)
	if err != nil {
		return res.Application, err
	}
	req.Header.Set("Accept", "application/json")
	err = c.Do(ctx, req, res)
	if err != nil {
		return nil, fmt.Errorf("GetApplication failed, error: %s", err)
	}
	return res.Application, nil
}

func Heartbeat(zone, app, instanceID string) error {
	u := zone + heartbeatUrl + app + "/" + instanceID
	ctx := context.Background()
	params := url.Values{
		"status": {"UP"},
	}
	var res interface{}
	resp, err := c.RawPut(ctx, u, "", params, res)
	if resp != nil && resp.StatusCode == http.StatusNotFound {
		err = ErrNotFound
		return err
	}
	if err != nil {
		return fmt.Errorf("heartbeat failed, error: %s", err)
	}
	return nil
}
