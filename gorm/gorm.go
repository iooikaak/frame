package gorm

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/opentracing/opentracing-go"
)

const (
	parentSpanGormKey = "opentracingParentSpan"
	spanGormKey       = "opentracingSpan"
	spanDuration      = "opentracingSpanDuration"
)

var (
	defaultDatabase     = "mysql"
	ConnStrTmpl         = "%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=%s&timeout=%s&readTimeout=%s&writeTimeout=%s"
	DefaultMaxOpenConns = 200
	DefaultMaxIdleConns = 60
	DefaultMaxLeftTime  = 300 * time.Second
	DefaultTimeout      = 5 * time.Second
	DefaultWriteTimeout = 5 * time.Second
	DefaultReadTimeout  = 10 * time.Second
	Charset             = "utf8mb4"
	MPort               = 3306
	TimeZone            = "Local"
	gormEngine          *Engine
)

type Engine struct {
	gorm *gorm.DB
}

//New 实例化新的Gorm实例
func New(conf *Config) *Engine {
	err := authConfig(conf)
	if err != nil {
		panic(err)
	}

	if strings.TrimSpace(conf.Charset) == "" {
		conf.Charset = Charset
	}

	if strings.TrimSpace(conf.TimeZone) == "" {
		conf.TimeZone = TimeZone
	}

	connStr := fmt.Sprintf(
		ConnStrTmpl,
		conf.User,
		conf.Password,
		conf.Server,
		conf.Port,
		conf.Database,
		conf.Charset,
		conf.TimeZone,
		conf.Timeout,
		conf.ReadTimeout,
		conf.WriteTimeout)

	db, err := gorm.Open(conf.Type, connStr)
	if err != nil {
		panic(err)
	}
	gormEngine = &Engine{db}
	gormEngine.wrapLog()
	db.DB().SetConnMaxLifetime(conf.MaxLeftTime)
	db.DB().SetMaxIdleConns(conf.MaxIdleConns)
	db.DB().SetMaxOpenConns(conf.MaxOpenConns)

	addGormCallbacks(db)
	return gormEngine
}

func IsByteArrayOrSlice(v reflect.Value) bool {
	return gorm.IsByteArrayOrSlice(v)
}
func IsRecordNotFoundError(err error) bool {
	return gorm.IsRecordNotFoundError(err)
}

func AddNamingStrategy(ns *gorm.NamingStrategy) {
	gorm.AddNamingStrategy(ns)
}

func Expr(expression string, args ...interface{}) *gorm.SqlExpr {
	return gorm.Expr(expression, args...)
}

func (db *Engine) Context(ctx context.Context) *gorm.DB {
	if ctx == nil {
		return db.gorm
	}

	if parentSpan := opentracing.SpanFromContext(ctx); parentSpan != nil {
		return db.gorm.Set(parentSpanGormKey, parentSpan)
	}

	return db.gorm
}

func (db *Engine) SetLogMode(mode bool) {
	db.gorm.LogMode(mode)
}

func (db *Engine) Close() error {
	return db.gorm.Close()
}

func Context(ctx context.Context) *gorm.DB {
	if gormEngine == nil {
		panic(fmt.Errorf("must init gorm.New"))
	}

	if parentSpan := opentracing.SpanFromContext(ctx); parentSpan != nil {
		return gormEngine.gorm.Set(parentSpanGormKey, parentSpan)
	}

	return gormEngine.gorm
}

func addGormCallbacks(db *gorm.DB) {
	callbacks := newCallbacks()
	registerCallbacks(db, "create", callbacks)
	registerCallbacks(db, "query", callbacks)
	registerCallbacks(db, "update", callbacks)
	registerCallbacks(db, "delete", callbacks)
	registerCallbacks(db, "row_query", callbacks)
}

type callbacks struct{}

func newCallbacks() *callbacks {
	return &callbacks{}
}

func (c *callbacks) beforeCreate(scope *gorm.Scope)   { c.before(scope) }
func (c *callbacks) afterCreate(scope *gorm.Scope)    { c.after(scope) }
func (c *callbacks) beforeQuery(scope *gorm.Scope)    { c.before(scope) }
func (c *callbacks) afterQuery(scope *gorm.Scope)     { c.after(scope) }
func (c *callbacks) beforeUpdate(scope *gorm.Scope)   { c.before(scope) }
func (c *callbacks) afterUpdate(scope *gorm.Scope)    { c.after(scope) }
func (c *callbacks) beforeDelete(scope *gorm.Scope)   { c.before(scope) }
func (c *callbacks) afterDelete(scope *gorm.Scope)    { c.after(scope) }
func (c *callbacks) beforeRowQuery(scope *gorm.Scope) { c.before(scope) }
func (c *callbacks) afterRowQuery(scope *gorm.Scope)  { c.after(scope) }

func (c *callbacks) before(scope *gorm.Scope) {
	val, ok := scope.Get(parentSpanGormKey)
	if !ok {
		return
	}
	scope.Set(spanDuration, time.Now())
	scope.Set(spanGormKey, opentracing.StartSpan("GORM-SQL", opentracing.ChildOf(val.(opentracing.Span).Context())))
}

func (c *callbacks) after(scope *gorm.Scope) {
	val, ok := scope.Get(spanGormKey)
	if !ok {
		return
	}
	sp := val.(opentracing.Span)

	t, ok := scope.Get(spanDuration)
	if !ok {
		t = time.Now()
	}

	sp.SetTag("db.statement", scope.SQLVars)
	sp.SetTag("db.instance", scope.InstanceID())
	sp.SetTag("db.sql", scope.SQL)
	sp.SetTag("db.err", scope.HasError())
	sp.SetTag("db.took", time.Since(t.(time.Time)))

	sp.Finish()
}

func registerCallbacks(db *gorm.DB, name string, c *callbacks) {
	beforeName := fmt.Sprintf("tracing:%v_before", name)
	afterName := fmt.Sprintf("tracing:%v_after", name)
	gormCallbackName := fmt.Sprintf("gorm:%v", name)

	switch name {
	case "create":
		db.Callback().Create().Before(gormCallbackName).Register(beforeName, c.beforeCreate)
		db.Callback().Create().After(gormCallbackName).Register(afterName, c.afterCreate)
	case "query":
		db.Callback().Query().Before(gormCallbackName).Register(beforeName, c.beforeQuery)
		db.Callback().Query().After(gormCallbackName).Register(afterName, c.afterQuery)
	case "update":
		db.Callback().Update().Before(gormCallbackName).Register(beforeName, c.beforeUpdate)
		db.Callback().Update().After(gormCallbackName).Register(afterName, c.afterUpdate)
	case "delete":
		db.Callback().Delete().Before(gormCallbackName).Register(beforeName, c.beforeDelete)
		db.Callback().Delete().After(gormCallbackName).Register(afterName, c.afterDelete)
	case "row_query":
		db.Callback().RowQuery().Before(gormCallbackName).Register(beforeName, c.beforeRowQuery)
		db.Callback().RowQuery().After(gormCallbackName).Register(afterName, c.afterRowQuery)
	}
}
