package xorm

import (
	"context"
	"database/sql"
	"reflect"
	"time"

	"github.com/xormplus/xorm/dialects"
	"github.com/xormplus/xorm/log"

	"github.com/xormplus/xorm/schemas"

	"github.com/xormplus/core"
	"github.com/xormplus/xorm"
)

//xormplus 核心接口
type Interface interface {
	AllCols() *Sessionx
	Alias(alias string) *Sessionx
	Asc(colNames ...string) *Sessionx
	BufferSize(size int) *Sessionx
	Cols(columns ...string) *Sessionx
	Count(...interface{}) (int64, error)
	CreateIndexes(bean interface{}) error
	CreateUniques(bean interface{}) error
	Decr(column string, arg ...interface{}) *Sessionx
	Desc(...string) *Sessionx
	Delete(interface{}) (int64, error)
	Distinct(columns ...string) *Sessionx
	DropIndexes(bean interface{}) error
	Exec(sqlOrAgrs ...interface{}) (sql.Result, error)
	Exist(bean ...interface{}) (bool, error)
	Find(interface{}, ...interface{}) error
	FindAndCount(interface{}, ...interface{}) (int64, error)
	Get(interface{}) (bool, error)
	GroupBy(keys string) *Sessionx
	ID(interface{}) *Sessionx
	In(string, ...interface{}) *Sessionx
	Incr(column string, arg ...interface{}) *Sessionx
	Insert(...interface{}) (int64, error)
	InsertOne(interface{}) (int64, error)
	IsTableEmpty(bean interface{}) (bool, error)
	IsTableExist(beanOrTableName interface{}) (bool, error)
	Iterate(interface{}, xorm.IterFunc) error
	Limit(int, ...int) *Sessionx
	MustCols(columns ...string) *Sessionx
	NoAutoCondition(...bool) *Sessionx
	NotIn(string, ...interface{}) *Sessionx
	Join(joinOperator string, tablename interface{}, condition string, args ...interface{}) *Sessionx
	Omit(columns ...string) *Sessionx
	OrderBy(order string) *Sessionx
	Ping() error
	QueryBytes(sqlOrAgrs ...interface{}) (resultsSlice []map[string][]byte, err error)
	QueryInterface(sqlOrArgs ...interface{}) ([]map[string]interface{}, error)
	QueryString(sqlOrArgs ...interface{}) ([]map[string]string, error)
	QueryValue(sqlOrArgs ...interface{}) ([]map[string]xorm.Value, error)
	QueryResult(sqlOrArgs ...interface{}) (result *xorm.ResultValue)
	Rows(bean interface{}) (*xorm.Rows, error)
	SetExpr(string, interface{}) *Sessionx
	SQL(interface{}, ...interface{}) *Sessionx
	Sum(bean interface{}, colName string) (float64, error)
	SumInt(bean interface{}, colName string) (int64, error)
	Sums(bean interface{}, colNames ...string) ([]float64, error)
	SumsInt(bean interface{}, colNames ...string) ([]int64, error)
	Table(tableNameOrBean interface{}) *Sessionx
	Unscoped() *Sessionx
	Update(bean interface{}, condiBeans ...interface{}) (int64, error)
	UseBool(...string) *Sessionx
	Where(interface{}, ...interface{}) *Sessionx
}

//xormplus 核心Engine
type EngineInterface interface {
	//Interface
	ClearCache(...interface{}) error
	CreateTables(...interface{}) error
	DBMetas() ([]*schemas.Table, error)
	Dialect() dialects.Dialect
	DropTables(...interface{}) error
	DumpAllToFile(fp string, tp ...schemas.DBType) error
	GetCacher(string) core.Cacher
	GetColumnMapper() core.IMapper
	GetDefaultCacher() core.Cacher
	GetTableMapper() core.IMapper
	GetTZDatabase() *time.Location
	GetTZLocation() *time.Location
	MapCacher(interface{}, core.Cacher) error
	NewSession(context.Context) *Sessionx
	Quote(string) string
	SetCacher(string, core.Cacher)
	SetConnMaxLifetime(time.Duration)
	SetColumnMapper(core.IMapper)
	SetDefaultCacher(core.Cacher)
	SetLogger(logger core.ILogger)
	SetLogLevel(log.LogLevel)
	SetMapper(core.IMapper)
	SetMaxOpenConns(int)
	SetMaxIdleConns(int)
	SetSchema(string)
	SetTableMapper(core.IMapper)
	SetTZDatabase(tz *time.Location)
	SetTZLocation(tz *time.Location)
	ShowSQL(show ...bool)
	Sync(...interface{}) error
	Sync2(...interface{}) error
	TableInfo(bean interface{}) (*schemas.Table, error)
	TableName(interface{}, ...bool) string
	UnMapType(reflect.Type)
}

var (
	_ Interface       = &Sessionx{}
	_ EngineInterface = &Engine{}
)
