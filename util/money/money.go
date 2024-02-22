package money

// FenToYuan 分转换为元
import (
	"fmt"
	"strconv"
	"strings"
)

func FenToYuan(number int) string {
	yuan := number / 100
	remind := number % 100
	jiao := remind / 10
	fen := remind % 10
	reslutStr := fmt.Sprint(yuan) + "." + fmt.Sprint(jiao) + fmt.Sprint(fen)
	return reslutStr
}

// FenToYuan2 分转换为元
func FenToYuan2(number string) string {
	l := len(number)
	if l == 1 {
		return "0.0" + number
	} else if l == 2 {
		return "0." + number
	} else {
		return number[:l-2] + "." + number[l-2:]
	}
}

// YuanToFen 字符串元转为 分
func YuanToFen(yuan string) (int64, error) {
	y := strings.Split(yuan, ".")
	if len(y) != 2 {
		return 0, fmt.Errorf("bad yuan:%s", yuan)
	}
	if len(y[1]) != 2 {
		return 0, fmt.Errorf("bad yuan:%s", yuan)
	}
	return strconv.ParseInt(y[0]+y[1], 10, 64)
}
