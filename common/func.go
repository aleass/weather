package common

import "strconv"

func Str2Int64(val string) int64 {
	if val == "" {
		return 0
	}
	num, _ := strconv.Atoi(val)
	return int64(num)
}

func Int642Str(val int64) string {
	if val == 0 {
		return ""
	}
	return strconv.Itoa(int(val))
}

func DefaultVal(val string) string {
	if val == "" {
		return "0"
	}
	return val
}
