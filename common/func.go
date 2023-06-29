package common

import "strconv"

func Str2Int64(val string) int64 {
	if val == "" {
		return 0
	}
	num, _ := strconv.Atoi(val)
	return int64(num)
}
