package sql

import (
	"github.com/iooikaak/frame/net/netutil/breaker"
	"github.com/iooikaak/frame/time"
	"github.com/iooikaak/frame/xlog"

	// database driver
	_ "github.com/go-sql-driver/mysql"
)

// Config mysql config.
type Config struct {
	DSN          string          // write data source name.
	ReadDSN      []string        // read data source name.
	Active       int             // pool
	Idle         int             // pool
	IdleTimeout  time.Duration   // connect max life time.
	QueryTimeout time.Duration   // query sql timeout
	ExecTimeout  time.Duration   // execute sql timeout
	TranTimeout  time.Duration   // transaction sql timeout
	Breaker      *breaker.Config // breaker
}

// NewMySQL new db and retry connection when has error.
func NewMySQL(c *Config) (db *DB) {
	if c.QueryTimeout == 0 || c.ExecTimeout == 0 || c.TranTimeout == 0 {
		panic("mysql must be set query/execute/transaction timeout")
	}
	db, err := Open(c)
	if err != nil {
		xlog.Error("open mysql error(%v)", err)
		panic(err)
	}
	return
}
