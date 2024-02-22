package breaker

import (
	"github.com/afex/hystrix-go/hystrix"
	"github.com/iooikaak/frame/xlog"
)

type CommandConfig hystrix.CommandConfig

var (
	Go               = hystrix.Go
	GoC              = hystrix.GoC
	Do               = hystrix.Do
	DoC              = hystrix.DoC
	Configure        = hystrix.Configure
	ConfigureCommand = hystrix.ConfigureCommand
	Flush            = hystrix.Flush
	// 影响后续替换的直接干掉
	//GetCircuit         = hystrix.GetCircuit
	//GetCircuitSettings = hystrix.GetCircuitSettings
	//NewStreamHandler   = hystrix.NewStreamHandler
)

var (
	// DefaultTimeout is how long to wait for command to complete, in milliseconds
	DefaultTimeout = &hystrix.DefaultTimeout
	// DefaultMaxConcurrent is how many commands of the same type can run at the same time
	DefaultMaxConcurrent = &hystrix.DefaultMaxConcurrent
	// DefaultVolumeThreshold is the minimum number of requests needed before a circuit can be tripped due to health
	DefaultVolumeThreshold = &hystrix.DefaultVolumeThreshold
	// DefaultSleepWindow is how long, in milliseconds, to wait after a circuit opens before testing for recovery
	DefaultSleepWindow = &hystrix.DefaultSleepWindow
	// DefaultErrorPercentThreshold causes circuits to open once the rolling measure of errors exceeds this percent of requests
	DefaultErrorPercentThreshold = &hystrix.DefaultErrorPercentThreshold
)

var (
	// ErrMaxConcurrency occurs when too many of the same named command are executed at the same time.
	ErrMaxConcurrency = hystrix.ErrMaxConcurrency
	// ErrCircuitOpen returns when an execution attempt "short circuits". This happens due to the circuit being measured as unhealthy.
	ErrCircuitOpen = hystrix.ErrCircuitOpen
	// ErrTimeout occurs when the provided function takes too long to execute.
	ErrTimeout = hystrix.ErrTimeout
)

func init() {
	hystrix.SetLogger(xlog.Logger())
}
