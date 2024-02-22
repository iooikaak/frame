package util

import (
	"time"
)

const (
	timeFormart = "2006-01-02 15:04:05"
)

//
type XTime time.Time

func (t XTime) String() string {
	return time.Time(t).Format(timeFormart)
}

func (t XTime) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(timeFormart)+2)
	b = append(b, '"')
	b = time.Time(t).AppendFormat(b, timeFormart)
	b = append(b, '"')
	return b, nil
}

func (t *XTime) UnmarshalJSON(data []byte) (err error) {
	now, err := time.ParseInLocation(`"`+timeFormart+`"`, string(data), time.Local)
	*t = XTime(now)
	return
}
