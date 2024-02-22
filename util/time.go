package util

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"time"
)

const TimeLayout = "2006-01-02 15:04:05"

type JsonTime struct {
	time.Time
}

func (t JsonTime) MarshalJSON() ([]byte, error) {
	var stamp = fmt.Sprintf("\"%s\"", t.Time.Format(TimeLayout))
	return []byte(stamp), nil
}

func (t *JsonTime) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		return nil
	}
	var err error
	str, err := strconv.Unquote(string(data))
	if err != nil {
		return err
	}
	i, err := time.ParseInLocation(TimeLayout, str, time.Local)
	t.Time = i
	return err
}

func (t JsonTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

func (t JsonTime) Empty() bool {
	return t == JsonTime{} || t.Time.Unix() == 0
}

func (t *JsonTime) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = JsonTime{Time: value}
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}

type JsonTimeSlice []JsonTime

func (s JsonTimeSlice) Len() int { return len(s) }

func (s JsonTimeSlice) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s JsonTimeSlice) Less(i, j int) bool { return s[i].Before(s[j].Time) }
