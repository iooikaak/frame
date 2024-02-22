package xorm

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/xormplus/xorm"
)

//吐血的包装....
type Sessionx struct {
	Sess *xorm.Session
	span opentracing.Span
}

func NewSession(ctx context.Context, x *xorm.Engine) *Sessionx {
	var span opentracing.Span
	if span = opentracing.SpanFromContext(ctx); span != nil {
		span = opentracing.StartSpan("XORM-SQL", opentracing.ChildOf(span.Context()))
	}
	s := &Sessionx{Sess: x.NewSession(), span: span}
	s.Sess.Context(ctx)
	return s
}

func (s *Sessionx) Get(bean interface{}) (bool, error) {
	b4ExecTime := time.Now()
	b, err := s.Sess.Get(bean)
	s.finish(b4ExecTime)
	return b, err
}

func (s *Sessionx) SetSpan(span opentracing.Span) {
	s.span = span
}

func (s *Sessionx) SetSession(sess *xorm.Session) {
	s.Sess = sess
}

func (s *Sessionx) Find(rowsSlicePtr interface{}, condiBean ...interface{}) error {
	b4ExecTime := time.Now()
	err := s.Sess.Find(rowsSlicePtr, condiBean...)
	s.finish(b4ExecTime)
	return err
}

func (s *Sessionx) Count(v ...interface{}) (n int64, err error) {
	b4ExecTime := time.Now()
	n, err = s.Sess.Count(v...)
	s.finish(b4ExecTime)
	return
}
func (s *Sessionx) CreateIndexes(bean interface{}) error {
	b4ExecTime := time.Now()
	err := s.Sess.CreateIndexes(bean)
	s.finish(b4ExecTime)
	return err
}
func (s *Sessionx) CreateUniques(bean interface{}) error {
	b4ExecTime := time.Now()
	err := s.Sess.CreateUniques(bean)
	s.finish(b4ExecTime)
	return err
}

func (s *Sessionx) Delete(w interface{}) (int64, error) {
	b4ExecTime := time.Now()
	n, err := s.Sess.Delete(w)
	s.finish(b4ExecTime)
	return n, err
}

func (s *Sessionx) DropIndexes(bean interface{}) error {
	b4ExecTime := time.Now()
	err := s.Sess.DropIndexes(bean)
	s.finish(b4ExecTime)
	return err
}

func (s *Sessionx) Exec(sqlOrArgs ...interface{}) (r sql.Result, err error) {
	b4ExecTime := time.Now()
	r, err = s.Sess.Exec(sqlOrArgs...)
	s.finish(b4ExecTime)
	return
}
func (s *Sessionx) Exist(bean ...interface{}) (bool, error) {
	b4ExecTime := time.Now()
	b, err := s.Sess.Exist(bean...)
	s.finish(b4ExecTime)
	return b, err
}

func (s *Sessionx) FindAndCount(w interface{}, arg ...interface{}) (int64, error) {
	b4ExecTime := time.Now()
	n, err := s.Sess.FindAndCount(w, arg...)
	s.finish(b4ExecTime)
	return n, err
}

func (s *Sessionx) Insert(d ...interface{}) (int64, error) {
	b4ExecTime := time.Now()
	n, err := s.Sess.Insert(d...)
	s.finish(b4ExecTime)
	return n, err
}
func (s *Sessionx) InsertOne(w interface{}) (int64, error) {
	b4ExecTime := time.Now()
	n, err := s.Sess.InsertOne(w)
	s.finish(b4ExecTime)
	return n, err
}
func (s *Sessionx) IsTableEmpty(bean interface{}) (bool, error) {
	b4ExecTime := time.Now()
	b, err := s.Sess.IsTableEmpty(bean)
	s.finish(b4ExecTime)
	return b, err
}
func (s *Sessionx) IsTableExist(beanOrTableName interface{}) (bool, error) {
	b4ExecTime := time.Now()
	b, err := s.Sess.IsTableExist(beanOrTableName)
	s.finish(b4ExecTime)
	return b, err
}
func (s *Sessionx) Iterate(w interface{}, x xorm.IterFunc) error {
	b4ExecTime := time.Now()
	err := s.Sess.Iterate(w, x)
	s.finish(b4ExecTime)
	return err
}

func (s *Sessionx) Ping() error {
	return s.Sess.Ping()
}
func (s *Sessionx) Query(sqlOrArgs ...interface{}) (resultsSlice []map[string][]byte, err error) {
	b4ExecTime := time.Now()
	resultsSlice, err = s.Sess.QueryBytes(sqlOrArgs)
	s.finish(b4ExecTime)
	return
}

func (s *Sessionx) QueryBytes(sqlOrArgs ...interface{}) (resultsSlice []map[string][]byte, err error) {
	b4ExecTime := time.Now()
	resultsSlice, err = s.Sess.QueryBytes(sqlOrArgs)
	s.finish(b4ExecTime)
	return
}
func (s *Sessionx) QueryInterface(sqlOrArgs ...interface{}) (resultsSlice []map[string]interface{}, err error) {
	b4ExecTime := time.Now()
	resultsSlice, err = s.Sess.QueryInterface(sqlOrArgs...)
	s.finish(b4ExecTime)
	return
}
func (s *Sessionx) QueryString(sqlOrArgs ...interface{}) (resultsSlice []map[string]string, err error) {
	b4ExecTime := time.Now()
	resultsSlice, err = s.Sess.QueryString(sqlOrArgs...)
	s.finish(b4ExecTime)
	return
}
func (s *Sessionx) Rows(bean interface{}) (r *xorm.Rows, err error) {
	b4ExecTime := time.Now()
	r, err = s.Sess.Rows(bean)
	s.finish(b4ExecTime)
	return
}

func (s *Sessionx) Sum(bean interface{}, colName string) (float64, error) {
	b4ExecTime := time.Now()
	f, err := s.Sess.Sum(bean, colName)
	s.finish(b4ExecTime)
	return f, err
}
func (s *Sessionx) SumInt(bean interface{}, colName string) (int64, error) {
	b4ExecTime := time.Now()
	n, err := s.Sess.SumInt(bean, colName)
	s.finish(b4ExecTime)
	return n, err
}
func (s *Sessionx) Sums(bean interface{}, colNames ...string) ([]float64, error) {
	b4ExecTime := time.Now()
	arr, err := s.Sess.Sums(bean, colNames...)
	s.finish(b4ExecTime)
	return arr, err
}
func (s *Sessionx) SumsInt(bean interface{}, colNames ...string) ([]int64, error) {
	b4ExecTime := time.Now()
	arr, err := s.Sess.SumsInt(bean, colNames...)
	s.finish(b4ExecTime)
	return arr, err
}

func (s *Sessionx) Update(bean interface{}, condiBeans ...interface{}) (n int64, err error) {
	b4ExecTime := time.Now()
	n, err = s.Sess.Update(bean, condiBeans...)
	s.finish(b4ExecTime)
	return
}

// Deprecated: xorm does not support
func (s *Sessionx) QueryValue(sqlOrArgs ...interface{}) (res []map[string]xorm.Value, err error) {
	b4ExecTime := time.Now()
	res, err = s.Sess.QueryValue(sqlOrArgs...)
	s.finish(b4ExecTime)
	return
}

// Deprecated: xorm does not support
func (s *Sessionx) QueryResult(sqlOrArgs ...interface{}) (result *xorm.ResultValue) {
	b4ExecTime := time.Now()
	result = s.Sess.QueryResult(sqlOrArgs...)
	s.finish(b4ExecTime)
	return
}

func (s *Sessionx) SQL(query interface{}, args ...interface{}) *Sessionx {
	s.Sess.SQL(query, args...)
	return s
}

func (s *Sessionx) Sql(query string, args ...interface{}) *Sessionx {
	s.Sess.SQL(query, args...)
	return s
}

func (s *Sessionx) AllCols() *Sessionx {
	s.Sess.AllCols()
	return s
}

func (s *Sessionx) Alias(alias string) *Sessionx {
	s.Sess.Alias(alias)
	return s
}
func (s *Sessionx) Asc(colNames ...string) *Sessionx {
	s.Sess.Asc(colNames...)
	return s
}
func (s *Sessionx) BufferSize(size int) *Sessionx {
	s.Sess.BufferSize(size)
	return s
}
func (s *Sessionx) Cols(columns ...string) *Sessionx {
	s.Sess.Cols(columns...)
	return s
}

func (s *Sessionx) Select(str string) *Sessionx {
	s.Sess.Select(str)
	return s
}

func (s *Sessionx) Decr(column string, arg ...interface{}) *Sessionx {
	s.Sess.Decr(column, arg...)
	return s
}
func (s *Sessionx) Desc(c ...string) *Sessionx {
	s.Sess.Desc(c...)
	return s
}

func (s *Sessionx) Distinct(columns ...string) *Sessionx {
	s.Sess.Distinct(columns...)
	return s
}

func (s *Sessionx) GroupBy(keys string) *Sessionx {
	s.Sess.GroupBy(keys)
	return s
}
func (s *Sessionx) ID(i interface{}) *Sessionx {
	s.Sess.ID(i)
	return s
}
func (s *Sessionx) In(fields string, v ...interface{}) *Sessionx {
	s.Sess.In(fields, v...)
	return s
}
func (s *Sessionx) Incr(column string, arg ...interface{}) *Sessionx {
	s.Sess.Incr(column, arg...)
	return s
}

func (s *Sessionx) Limit(i int, e ...int) *Sessionx {
	s.Sess.Limit(i, e...)
	return s
}
func (s *Sessionx) MustCols(columns ...string) *Sessionx {
	s.Sess.MustCols(columns...)
	return s
}
func (s *Sessionx) NoAutoCondition(b ...bool) *Sessionx {
	s.Sess.NoAutoCondition(b...)
	return s
}
func (s *Sessionx) NotIn(str string, i ...interface{}) *Sessionx {
	s.Sess.NotIn(str, i...)
	return s
}
func (s *Sessionx) Join(joinOperator string, tablename interface{}, condition string, args ...interface{}) *Sessionx {
	s.Sess.Join(joinOperator, tablename, condition, args...)
	return s
}
func (s *Sessionx) Omit(columns ...string) *Sessionx {
	s.Sess.Omit(columns...)
	return s
}
func (s *Sessionx) OrderBy(order string) *Sessionx {
	s.Sess.OrderBy(order)
	return s
}

func (s *Sessionx) SetExpr(str string, i interface{}) *Sessionx {
	s.Sess.SetExpr(str, i)
	return s
}

func (s *Sessionx) Table(tableNameOrBean interface{}) *Sessionx {
	s.Sess.Table(tableNameOrBean)
	return s
}
func (s *Sessionx) Unscoped() *Sessionx {
	s.Sess.Unscoped()
	return s
}
func (s *Sessionx) UseBool(str ...string) *Sessionx {
	s.Sess.UseBool(str...)
	return s
}
func (s *Sessionx) Where(w interface{}, arg ...interface{}) *Sessionx {
	s.Sess.Where(w, arg...)
	return s
}

func (s *Sessionx) LastSQL() (string, []interface{}) {
	sqlx, arg := s.Sess.LastSQL()
	return sqlx, arg
}

func (s *Sessionx) Begin() error {
	return s.Sess.Begin()
}

func (s *Sessionx) Commit() error {
	return s.Sess.Commit()
}

func (s *Sessionx) Rollback() error {
	return s.Sess.Rollback()
}

func (s *Sessionx) Close() {
	if s.span != nil {
		s.span.Finish()
		s.span = nil
	}
	s.Sess.Close()
}

func (s *Sessionx) finish(t time.Time) {

	if s.span != nil {
		sqlx, arg := s.Sess.LastSQL()
		s.span.SetTag("xorm.sql", fmt.Sprintf("sql:%s  param:%v", sqlx, arg))
		s.span.SetTag("xorm.took", time.Since(t))
	}
}
