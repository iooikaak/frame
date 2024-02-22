package dts

import (
	"fmt"
	"strconv"
	"time"
	"unsafe"
)

const (
	Update     = "UPDATE"
	Insert     = "INSERT"
	Delete     = "DELETE"
	TimeLayout = "2006-01-02 15:04:05"
)

//parse db field
func FieldsName(d map[string]interface{}) (res []string) {
	if dd, ok := d["fields"].(map[string]interface{}); ok {
		if arr, ok := dd["array"].([]interface{}); ok {
			for _, v := range arr {
				if vv, ok := v.(map[string]interface{}); ok {
					switch vv["name"].(type) {
					case []byte:
						res = append(res, bytesToString(vv["name"].([]byte)))
					case string:
						res = append(res, vv["name"].(string))
					}
				}
			}
		}
	}
	return
}

func bytesToString(bytes []byte) string {
	return *(*string)(unsafe.Pointer(&bytes))
}

//map value to files key->value
func FieldsKeyValue(data interface{}, fieldsName []string) (map[string]string, bool) {
	res := make(map[string]string)
	if dd, ok := data.(map[string]interface{}); ok {
		if arr, ok := dd["array"].([]interface{}); ok {
			if len(arr) != len(fieldsName) {
				return res, false
			}
			for k := range arr {
				if vv, ok := arr[k].(map[string]interface{}); ok {
					for _, vvv := range vv {
						if vvvv, ok := vvv.(map[string]interface{}); ok {
							if i, ok := vvvv["value"]; ok {
								switch i.(type) {
								case []byte:
									res[fieldsName[k]] = bytesToString(vvvv["value"].([]byte))
								case string:
									res[fieldsName[k]] = vvvv["value"].(string)
								case float64:
									res[fieldsName[k]] = strconv.FormatFloat(vvvv["value"].(float64), 'f', 6, 64)
								default:
									res[fieldsName[k]] = fmt.Sprintf("%v", vvvv["value"])
								}
							}
							if _, ok := vvvv["year"]; ok {
								tArr := make(map[string]int)
								for k, ii := range vvvv {
									iii := ii.(map[string]interface{})
									iiii := iii["int"]
									switch iiii := iiii.(type) {
									case int:
										tArr[k] = iiii
									case int32:
										tArr[k] = int(iiii)
									case int64:
										tArr[k] = int(iiii)
									case int8:
										tArr[k] = int(iiii)
									}

								}
								it := time.Date(tArr["year"], time.Month(tArr["month"]), tArr["day"], tArr["hour"], tArr["minute"], tArr["second"], 0, time.Local)
								res[fieldsName[k]] = it.Format(TimeLayout)
							}
						}
					}
				}
			}
		}
	}
	//TODO len(res) == len(fieldsName)
	if len(res) == 0 {
		return res, false
	}
	return res, true
}
