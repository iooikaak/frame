package util

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

func ReplaceAllSpace(s string) string {
	return strings.Replace(strings.Replace(strings.TrimSpace(s), " ", "", -1), "　", "", -1)
}

func StringArrayToSets(src []string) (dest []string) {
	var tmpMap = make(map[string]string)
	var ok bool
	dest = make([]string, 0, len(src))
	for _, i := range src {
		if _, ok = tmpMap[i]; i == "" || ok {
			continue
		}
		tmpMap[i] = i
		dest = append(dest, i)
	}
	return
}

func IntArrayToSets(src []int) (dest []int) {
	var tmpMap = make(map[int]int)
	var ok bool
	dest = make([]int, 0, len(src))
	for _, i := range src {
		if _, ok = tmpMap[i]; ok {
			continue
		}
		tmpMap[i] = i
		dest = append(dest, i)
	}
	return
}

func Int64ArrayToSets(src []int64) (dest []int64) {
	var tmpMap = make(map[int64]int64)
	var ok bool
	dest = make([]int64, 0, len(src))
	for _, i := range src {
		if _, ok = tmpMap[i]; ok {
			continue
		}
		tmpMap[i] = i
		dest = append(dest, i)
	}
	return
}

func Int64ArrayToInterfaceSets(src []int64) (dest []interface{}) {
	var tmpMap = make(map[int64]int64)
	var ok bool
	dest = make([]interface{}, 0, len(src))
	for _, i := range src {
		if _, ok = tmpMap[i]; ok {
			continue
		}
		tmpMap[i] = i
		dest = append(dest, i)
	}
	return
}

func UIntArrayToSets(src []uint) (dest []uint) {
	var tmpMap = make(map[uint]uint)
	var ok bool
	dest = make([]uint, 0, len(src))
	for _, i := range src {
		if _, ok = tmpMap[i]; ok {
			continue
		}
		tmpMap[i] = i
		dest = append(dest, i)
	}
	return
}

func UInt64ArrayToSets(src []uint64) (dest []uint64) {
	var tmpMap = make(map[uint64]uint64)
	var ok bool
	dest = make([]uint64, 0, len(src))
	for _, i := range src {
		if _, ok = tmpMap[i]; ok {
			continue
		}
		tmpMap[i] = i
		dest = append(dest, i)
	}
	return
}

func StringsToInts(src []string) (dest []int) {
	var err error
	var tmpInt int
	dest = make([]int, 0, len(src))
	for _, s := range src {
		if tmpInt, err = strconv.Atoi(s); err == nil {
			dest = append(dest, tmpInt)
		}
	}
	return
}

func StringsToInt64s(src []string) (dest []int64) {
	var err error
	var tmpInt int64
	dest = make([]int64, 0, len(src))
	for _, s := range src {
		if tmpInt, err = strconv.ParseInt(s, 10, 64); err == nil {
			dest = append(dest, tmpInt)
		}
	}
	return
}

func StringsToInt64Set(src []string) (dest []int64) {
	var err error
	var tmpInt int64
	var tmpMap = make(map[string]string)
	var ok bool
	dest = make([]int64, 0, len(src))
	for _, s := range src {
		if tmpInt, err = strconv.ParseInt(s, 10, 64); err == nil {
			if _, ok = tmpMap[s]; ok {
				continue
			}
			tmpMap[s] = s
			dest = append(dest, tmpInt)
		}
	}
	return
}

func StringsToUInt64s(src []string) (dest []uint64) {
	var err error
	var tmpInt uint64
	dest = make([]uint64, 0, len(src))
	for _, s := range src {
		if tmpInt, err = strconv.ParseUint(s, 10, 64); err == nil {
			dest = append(dest, tmpInt)
		}
	}
	return
}

func StringsToUInt64Set(src []string) (dest []uint64) {
	var err error
	var tmpInt uint64
	var tmpMap = make(map[string]string)
	var ok bool
	dest = make([]uint64, 0, len(src))
	for _, s := range src {
		if tmpInt, err = strconv.ParseUint(s, 10, 64); err == nil {
			if _, ok = tmpMap[s]; ok {
				continue
			}
			tmpMap[s] = s
			dest = append(dest, tmpInt)
		}
	}
	return
}

func Int64sToStrings(src []int64) (dest []string) {
	dest = make([]string, 0, len(src))
	for _, s := range src {
		dest = append(dest, strconv.FormatInt(s, 10))
	}
	return
}

func UInt64Join(a []uint64, sep string) string {
	switch len(a) {
	case 0:
		return ""
	case 1:
		return strconv.FormatUint(a[0], 10)
	case 2:
		return strconv.FormatUint(a[0], 10) + sep + strconv.FormatUint(a[1], 10)
	case 3:
		return strconv.FormatUint(a[0], 10) + sep + strconv.FormatUint(a[1], 10) + sep + strconv.FormatUint(a[2], 10)
	}
	n := len(sep) * (len(a) - 1)
	for i := 0; i < len(a); i++ {
		n += len(strconv.FormatUint(a[i], 10))
	}
	b := make([]byte, n)
	bp := copy(b, strconv.FormatUint(a[0], 10))
	for _, s := range a[1:] {
		bp += copy(b[bp:], sep)
		bp += copy(b[bp:], strconv.FormatUint(s, 10))
	}
	return string(b)
}

// 大于10000的数字转换成"%.1fk"格式，e.g.  10245 = 10.2k
func PrettyNumberV1(num int64) string {
	if num < 10000 {
		return fmt.Sprintf("%d", num)
	}
	return fmt.Sprintf("%.1fk", float64(num)/1000)
}

// 日期格式化V1，yyyyMMdd
func PrettyDateV1(t time.Time) string {
	return t.Format("20060102")
}

// 时间戳转字符串
func PrettyDateV2(ts int64) string {
	return time.Unix(ts, 0).Format("2006-01-02 15:04:05")
}

// 日期格式化V1，yyyy-MM-dd HH:mm:ss
func PrettyDateV3(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

// 日期格式化V4，yyyyMMddHH
func PrettyDateV4(t time.Time) string {
	return t.Format("2006010215")
}

func PrettyDateV5(ts int64) string {
	return time.Unix(ts, 0).Format("2006-01-02T15:04:05+08:00")
}

// 秒数转 xx:xx
func FormatSecondsV1(seconds int64) string {
	rem := seconds % 60
	quot := seconds / 60
	if quot < 100 {
		return fmt.Sprintf("%02d:%02d", quot, rem)
	}
	return fmt.Sprintf("%d:%02d", quot, rem)
}

// 浮点型精度处理
func Round(f float64, n int, roundDown bool) float64 {
	s := math.Pow10(n)
	if roundDown {
		return math.Floor(f*s) / s
	} else {
		return math.Ceil(f*s) / s
	}
}

// 整数四舍五入，12>10 15>20
func IntegerRound(i int64) int64 {
	b := i % 10
	if b < 5 {
		return i
	}
	return i + 10 - b
}
