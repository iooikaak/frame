package util

import (
	"math"
	"strconv"
)

var eps = 0.00000001 //设置容忍度

func FloatEquals(a, b float64) bool {
	return math.Abs(a-b) < eps
}

func InIn64tSliceMapKeyFunc(haystack []int64) func(int64) bool {
	set := make(map[int64]struct{})

	for _, e := range haystack {
		set[e] = struct{}{}
	}

	return func(needle int64) bool {
		_, ok := set[needle]
		return ok
	}
}

func Int64SliceCheckRepeat(a []int64, b []int64) bool {
	la := len(a)
	lb := len(b)
	if la == 0 || lb == 0 {
		return false
	}
	if la > lb {
		fn := InIn64tSliceMapKeyFunc(a)
		for _, i := range b {
			if fn(i) {
				return true
			}
		}
	} else {
		fn := InIn64tSliceMapKeyFunc(b)
		for _, i := range a {
			if fn(i) {
				return true
			}
		}
	}

	return false
}

func GetFloatStringFromIntCent(i int) string {
	j := float64(i) / 100
	return strconv.FormatFloat(j, 'f', 2, 64)
}

func GetIntCentFromFloat64(i float64) int {
	return int(math.Round(i * 100))
}
