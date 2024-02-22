package time

import (
	"database/sql/driver"
	"time"
)

const localDateTimeFormat string = "2006-01-02 15:04:05"

//TODO 单纯用于xorm datetime日期转换成"2006-01-02 15:04:05"仅此而为，但换来的代价是废除啦原有time本身应有的很多特性
//TODO 该类型不能用于指针类型,如果datetime存在null值，可能存在类型转换失败,使用前请明确知道使用范围
type LocalTime time.Time

func (l LocalTime) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(localDateTimeFormat)+2)
	b = append(b, '"')
	b = time.Time(l).AppendFormat(b, localDateTimeFormat)
	b = append(b, '"')
	return b, nil
}

func (l *LocalTime) UnmarshalJSON(b []byte) error {
	now, err := time.ParseInLocation(localDateTimeFormat, string(b), time.Local)
	*l = LocalTime(now)
	return err
}

func (l LocalTime) String() string {
	return time.Time(l).Format(localDateTimeFormat)
}

func (l LocalTime) Now() LocalTime {
	return LocalTime(time.Now())
}

func (l LocalTime) ParseTime(t time.Time) LocalTime {
	return LocalTime(t)
}

func (l LocalTime) format() string {
	return time.Time(l).Format(localDateTimeFormat)
}

func (l LocalTime) MarshalText() ([]byte, error) {
	return []byte(l.format()), nil
}

func (l *LocalTime) FromDB(b []byte) error {
	if nil == b || len(b) == 0 {
		l = nil
		return nil
	}
	var now time.Time
	var err error
	now, err = time.ParseInLocation(localDateTimeFormat, string(b), time.Local)
	if nil == err {
		*l = LocalTime(now)
		return nil
	}
	now, err = time.ParseInLocation("2006-01-02T15:04:05Z", string(b), time.Local)
	if nil == err {
		*l = LocalTime(now)
		return nil
	}
	panic("自己定義個layout日期格式處理一下數據庫裏面的日期型數據解析!")
}

func (l *LocalTime) ToDB() ([]byte, error) {
	if nil == l {
		return nil, nil
	}
	return []byte(time.Time(*l).Format(localDateTimeFormat)), nil
}

func (l *LocalTime) Value() (driver.Value, error) {
	if nil == l {
		return nil, nil
	}
	return time.Time(*l).Format(localDateTimeFormat), nil
}
