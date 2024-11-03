package common

import "strings"

// 经纬度处理
func LocStr2float(location string) (float64, float64) {
	loc := strings.Split(location, ",")
	if len(loc) != 2 {
		return 0, 0
	}
	return Str2Float64(loc[0]), Str2Float64(loc[1])
}
