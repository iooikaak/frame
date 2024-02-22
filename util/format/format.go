package format

import (
	"fmt"
	"strings"
)

// OrderNumberk 把订单数转为 k，当大于 1000
func OrderNumberk(num int64) string {
	if num < 1000 {
		return fmt.Sprintf("%d", num)
	}

	v := fmt.Sprintf("%.4f", float64(num)/1000.0)
	return v[:len(v)-3] + "k"
}

// 人民币（分）转狗粮
func FormatRMB2Gouliang(input int64) int64 {
	return input
}

func FormatAcceptOrderNumber3(num int64) string {
	if num < 10000 {
		return fmt.Sprintf("%d次", num)
	}
	return fmt.Sprintf("%.1f万次", float64(num)/10000)
}

// 大神评分 e.g. 5 > 5.0
func FormatScore(score int64) string {
	return fmt.Sprintf("%.1f", float64(score))
}

func FormatUserNameV1(userName string) string {
	s := []rune(userName)
	l := len(s)
	if l <= 1 {
		return "***" + userName
	}
	return fmt.Sprintf("%s***%s", string(s[0]), string(s[l-1]))
}

func FormatCommentTags(tags string) []string {
	return strings.Split(tags, ",")
}
