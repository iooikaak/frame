package sig

import (
	"io/ioutil"
	"net/http"
	"sort"
	"strings"

	"github.com/iooikaak/frame/utils"
)

func GetSign(salt string, req *http.Request) (r string, err error) {
	var signString string
	if req.Method == http.MethodGet {
		err = req.ParseForm()
		if err != nil {
			return
		}
		params := make(map[string]string)
		keys := make([]string, 0)
		for k, v := range req.Form {
			params[k] = v[0]
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, v := range keys {
			signString = v + "=" + params[v] + "&"
		}
		signString = strings.Trim(signString, "&")
		return utils.Md5Encrypt(signString + salt), nil
	}
	if req.Method == http.MethodPost {
		var b []byte
		b, err = ioutil.ReadAll(req.Body)
		if err != nil {
			return
		}
		signString = string(b)
		return utils.Md5Encrypt(signString + salt), nil
	}
	return
}
