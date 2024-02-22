package utils

// nolint
func CheckIsExistString(l []string, i string) bool {
	for _, v := range l {
		if v == i {
			return true
		}
	}
	return false
}

// nolint
func CheckIsExist(l []int64, i int64) bool {
	for _, v := range l {
		if v == i {
			return true
		}
	}
	return false
}
