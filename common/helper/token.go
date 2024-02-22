package helper

import (
	"crypto/md5"
	"fmt"
	"sort"
	"strings"
)

func GetToken(salt string, params map[string][]string) (r string, err error) {
	var (
		keys   []string
		values []string
	)
	for k := range params {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, v := range keys {
		value := params[v]
		item := v + "=" + value[0]
		values = append(values, item)
	}

	paramsStr := strings.Join(values, "&")
	paramsStr = paramsStr + salt
	fmt.Println(paramsStr)
	has := md5.Sum([]byte(paramsStr))
	r = fmt.Sprintf("%x", has) //将[]byte转成16进制
	return
}
