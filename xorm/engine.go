package xorm

import (
	"context"
	"reflect"
	"time"

	"github.com/xormplus/xorm/log"

	"github.com/xormplus/xorm/dialects"

	"github.com/xormplus/xorm/schemas"

	_ "github.com/go-sql-driver/mysql"
	"github.com/opentracing/opentracing-go"
	"github.com/xormplus/core"
	"github.com/xormplus/xorm"
)

// Engine .
type Engine struct {
	e *xorm.Engine
}

func NewEngine(driverName string, dataSourceName string) (eng *Engine, err error) {
	var e *xorm.Engine
	e, err = xorm.NewEngine(driverName, dataSourceName)
	if err != nil {
		return
	}

	eng = &Engine{e: e}
	return
}

func (engine *Engine) ClearCache(beans ...interface{}) error {
	return engine.e.ClearCache(beans...)
}

func (engine *Engine) CreateTables(beans ...interface{}) error {
	return engine.e.CreateTables(beans...)
}
func (engine *Engine) DBMetas() ([]*schemas.Table, error) {
	return engine.e.DBMetas()
}
func (engine *Engine) Dialect() dialects.Dialect {
	return engine.e.Dialect()
}
func (engine *Engine) DropTables(beans ...interface{}) error {
	return engine.e.DropTables(beans...)
}
func (engine *Engine) DumpAllToFile(fp string, tp ...schemas.DBType) error {
	return engine.e.DumpAllToFile(fp, tp...)
}
func (engine *Engine) GetCacher(s string) core.Cacher {
	return engine.e.GetCacher(s)
}
func (engine *Engine) GetColumnMapper() core.IMapper {
	return engine.e.GetColumnMapper()
}
func (engine *Engine) GetDefaultCacher() core.Cacher {
	return engine.e.GetDefaultCacher()
}
func (engine *Engine) GetTableMapper() core.IMapper {
	return engine.e.GetTableMapper()
}
func (engine *Engine) GetTZDatabase() *time.Location {
	return engine.e.GetTZDatabase()
}
func (engine *Engine) GetTZLocation() *time.Location {
	return engine.e.GetTZLocation()
}
func (engine *Engine) MapCacher(m interface{}, c core.Cacher) error {
	return engine.e.MapCacher(m, c)
}
func (engine *Engine) NewSession(ctx context.Context) *Sessionx {
	var span opentracing.Span
	if span = opentracing.SpanFromContext(ctx); span != nil {
		span = opentracing.StartSpan("XORM-SQL", opentracing.ChildOf(span.Context()))
	}
	sess := &Sessionx{Sess: engine.e.NewSession(), span: span}
	sess.Sess.Context(ctx)
	return sess
}

func (engine *Engine) Quote(q string) string {
	return engine.e.Quote(q)
}
func (engine *Engine) SetCacher(s string, c core.Cacher) {
	engine.e.SetCacher(s, c)
}
func (engine *Engine) SetConnMaxLifetime(t time.Duration) {
	engine.e.SetConnMaxLifetime(t)
}
func (engine *Engine) SetColumnMapper(c core.IMapper) {
	engine.e.SetColumnMapper(c)
}
func (engine *Engine) SetDefaultCacher(c core.Cacher) {
	engine.e.SetDefaultCacher(c)
}
func (engine *Engine) SetLogger(logger core.ILogger) {
	engine.e.SetLogger(logger)
}
func (engine *Engine) SetLogLevel(c log.LogLevel) {
	engine.e.SetLogLevel(c)
}
func (engine *Engine) SetMapper(c core.IMapper) {
	engine.e.SetMapper(c)
}
func (engine *Engine) SetMaxOpenConns(n int) {
	engine.e.SetMaxOpenConns(n)
}
func (engine *Engine) SetMaxIdleConns(n int) {
	engine.e.SetMaxIdleConns(n)
}
func (engine *Engine) SetSchema(s string) {
	engine.e.SetSchema(s)
}
func (engine *Engine) SetTableMapper(c core.IMapper) {
	engine.e.SetTableMapper(c)
}
func (engine *Engine) SetTZDatabase(tz *time.Location) {
	engine.e.SetTZDatabase(tz)
}
func (engine *Engine) SetTZLocation(tz *time.Location) {
	engine.e.SetTZLocation(tz)
}
func (engine *Engine) ShowSQL(show ...bool) {
	engine.e.ShowSQL(show...)
}
func (engine *Engine) Sync(beans ...interface{}) error {
	return engine.e.Sync(beans...)
}
func (engine *Engine) Sync2(beans ...interface{}) error {
	return engine.e.Sync2(beans...)
}

func (engine *Engine) TableInfo(bean interface{}) (*schemas.Table, error) {
	return engine.e.TableInfo(bean)
}
func (engine *Engine) TableName(t interface{}, b ...bool) string {
	return engine.e.TableName(t, b...)
}
func (engine *Engine) UnMapType(r reflect.Type) {
	engine.e.UnMapType(r)
}

func (engine *Engine) Ping() error {
	sess := engine.NewSession(context.Background())
	defer sess.Close()
	return sess.Ping()
}

func (engine *Engine) Close() error {
	return engine.e.Close()
}
