package tp

import (
	"strings"
	"time"
)

type Unixtime int64

//NewUnixtime 创建一个当前时间的Unixtime
func NewUnixtime() Unixtime {
	return Unixtime(time.Now().Unix())
}

func (t *Unixtime) UnmarshalJSON(data []byte) error {
	str := strings.Replace(string(data), `"`, "", -1)
	return t.SetString(str)
}

func (t Unixtime) MarshalJSON() ([]byte, error) {
	return []byte(`"` + t.String() + `"`), nil
}

//String 获取日期时间，格式：yyyy-MM-dd HH:mm:ss
func (t Unixtime) String() string {
	return formatTimeToString(t.Time())
}

//SetString 设置日期时间，字符串格式：yyyy-MM-dd HH:mm:ss
func (t *Unixtime) SetString(str string) error {
	tt, err := parseStringToTime(str)
	if err != nil {
		return err
	}

	t.SetTime(tt)
	return nil
}

//Time 获取 Unixtime的 time.Time 类型形式
func (t Unixtime) Time() time.Time {
	return time.Unix(t.Unix(), 0)
}

//SetTime 用 time.Time 类型来设置 Unixtime 的值
func (t *Unixtime) SetTime(tt time.Time) {
	*t = Unixtime(tt.Unix())
}

//Unix 获取 Datetime 的 Unix 时间戮，单位：秒
func (t *Unixtime) Unix() int64 {
	return int64(*t)
}

//SetUnix 用时间戮来设置 Datetime 的值
func (t *Unixtime) SetUnix(i64 int64) {
	*t = Unixtime(i64)
}
