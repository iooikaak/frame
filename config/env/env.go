// Package env get env & app config, all the public field must after init()
// finished and flag.Parse().
package env

import (
	"flag"
	"os"
	"strconv"
	"time"
)

// deploy env.
const (
	DeployEnvDev  = "dev"  // 开发环境
	DeployEnvQa   = "qa"   // qa环境
	DeployEnvSit  = "sit"  // 压测环境
	DeployEnvGray = "gray" // 灰度环境
	DeployEnvProd = "prod" // 正式环境
)

// env default value.
const (
	// env
	_region    = "region01"
	_zone      = "zone01"
	_deployEnv = "dev"
	_cloud     = "ali"
)

// env configuration.
var (
	// Cloud available cloud where app at.
	Cloud string
	// Region available region where app at.
	Region string
	// Zone available zone where app at.
	Zone string
	// Hostname machine hostname.
	Hostname string
	// DeployEnv deploy env where app at.
	DeployEnv string
	// AppID is global unique application id, register by service tree.
	// such as main.arch.disocvery.
	AppID string
	// Color is the identification of different experimental group in one caster cluster.
	Color string
	// DiscoveryNodes is seed nodes.
	DiscoveryNodes string
)

func init() {
	var err error
	Hostname = os.Getenv("HOSTNAME")
	if Hostname == "" {
		Hostname, err = os.Hostname()
		if err != nil {
			Hostname = strconv.Itoa(int(time.Now().UnixNano()))
		}
	}
	addFlag(flag.CommandLine)
}

func addFlag(fs *flag.FlagSet) {
	// env
	fs.StringVar(&Cloud, "cloud", defaultString("CLOUD", _cloud), "available cloud. or use CLOUD env variable, value: ali/tecent etc.")
	fs.StringVar(&Region, "region", defaultString("REGION", _region), "available region. or use REGION env variable, value: sh etc.")
	fs.StringVar(&Zone, "zone", defaultString("ZONE", _zone), "available zone. or use ZONE env variable, value: sh001/sh002 etc.")
	fs.StringVar(&AppID, "appid", os.Getenv("APP_ID"), "appid is global unique application id, register by service tree. or use APP_ID env variable.")
	fs.StringVar(&DeployEnv, "deploy.env", defaultString("DeployEnv", _deployEnv), "deploy env. or use DEPLOY_ENV env variable, value: dev/fat1/uat/pre/prod etc.")
	fs.StringVar(&Color, "deploy.color", os.Getenv("DeployColor"), "deploy.color is the identification of different experimental group.")
	fs.StringVar(&DiscoveryNodes, "discovery.nodes", os.Getenv("DISCOVERY_NODES"), "discovery.nodes is seed nodes. value: 127.0.0.1:7171,127.0.0.2:7171 etc.")
}

func defaultString(env, value string) string {
	v := os.Getenv(env)
	if v == "" {
		return value
	}
	return v
}
