package util

import (
	"math"
	"strings"
)

const (
	SYMBOLS = "0123456789abcdefghijklmnopqrsuvwxyzABCDEFGHIJKLMNOPQRSTUVXYZ-=_"
	BASE    = int64(len(SYMBOLS))
)

func Id2String(id int64) string {
	rest := id % BASE
	result := string(SYMBOLS[rest])
	if id-rest != 0 {
		newnumber := (id - rest) / BASE
		result = Id2String(newnumber) + result
	}
	return result
}

func String2Id(str string) int64 {
	const floatbase = float64(BASE)
	l := len(str)
	var sum int = 0
	for index := l - 1; index > -1; index -= 1 {
		current := string(str[index])
		pos := strings.Index(SYMBOLS, current)
		sum = sum + (pos * int(math.Pow(floatbase, float64((l-index-1)))))
	}
	return int64(sum)
}

// v2
const (
	SYMBOLS_V2 = "0123456789abcdefghijklmnopqrsuvwxyzABCDEFGHIJKLMNOPQRSTUVXYZ"
	BASE_V2    = int64(len(SYMBOLS_V2))
)

func Id2StringV2(id int64) string {
	rest := id % BASE_V2
	result := string(SYMBOLS_V2[rest])
	if id-rest != 0 {
		newnumber := (id - rest) / BASE_V2
		result = Id2StringV2(newnumber) + result
	}
	return result
}

func String2IdV2(str string) int64 {
	const floatbase = float64(BASE_V2)
	l := len(str)
	var sum int = 0
	for index := l - 1; index > -1; index -= 1 {
		current := string(str[index])
		pos := strings.Index(SYMBOLS_V2, current)
		sum = sum + (pos * int(math.Pow(floatbase, float64((l-index-1)))))
	}
	return int64(sum)
}
