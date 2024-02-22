package util

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"strconv"
	"time"

	"github.com/gogf/gf/container/garray"
)

// 去重
func Unique(ids []int) []int {
	a := garray.NewIntArray(false)
	for _, id := range ids {
		a.Append(id)
	}
	return a.Unique().Slice()
}

func Float64(i interface{}) float64 {
	switch ti := i.(type) {
	case int:
		return Float64(ti)
	case string:
		r, _ := strconv.ParseFloat(ti, 64)
		return r
	case json.Number:
		r, _ := ti.Float64()
		return r
	default:
		return i.(float64)
	}
}

func Int64(i interface{}) int64 {
	switch ti := i.(type) {
	case int:
		return int64(ti)
	case int64:
		return ti
	case string:
		r, _ := strconv.ParseInt(ti, 10, 64)
		return r
	case json.Number:
		r, _ := ti.Int64()
		return r
	default:
		return i.(int64)
	}
}

func IntToBytes(n int64) []byte {
	x := int32(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	_ = binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}

func RemoveRepByMap(slc []string) []string {
	result := make([]string, 0)
	tempMap := map[string]byte{}
	for _, e := range slc {
		l := len(tempMap)
		tempMap[e] = 0
		if len(tempMap) != l {
			result = append(result, e)
		}
	}
	return result
}

func RemoveRepByMapInt64(slc []int64) []int64 {
	result := make([]int64, 0)
	tempMap := map[int64]byte{}
	for _, e := range slc {
		l := len(tempMap)
		tempMap[e] = 0
		if len(tempMap) != l {
			result = append(result, e)
		}
	}
	return result
}

func RemoveRepByMapInt(slc []int) []int {
	result := make([]int, 0)
	tempMap := map[int]byte{}
	for _, e := range slc {
		l := len(tempMap)
		tempMap[e] = 0
		if len(tempMap) != l {
			result = append(result, e)
		}
	}
	return result
}

func CheckIsExist(l []int64, i int64) bool {
	for _, v := range l {
		if v == i {
			return true
		}
	}
	return false
}

func CheckIsExistInt(l []int, i int) bool {
	for _, v := range l {
		if v == i {
			return true
		}
	}
	return false
}

func CheckIsExistJsonTime(l []JsonTime, i JsonTime) bool {
	for _, v := range l {
		if v == i {
			return true
		}
	}
	return false
}

func CheckIsExistString(l []string, i string) bool {
	for _, v := range l {
		if v == i {
			return true
		}
	}
	return false
}

func Retry3(callback func() error) (err error) {
	for i := 0; i < 3; i++ {
		if err = callback(); err == nil {
			return
		}
		time.Sleep(100 * time.Millisecond)
	}
	return
}
