package common

import (
	"strconv"
	"strings"
)

func Str2Int64(val string) int64 {
	if val == "" {
		return 0
	}
	num, _ := strconv.Atoi(val)
	return int64(num)
}

func Str2Int(val string) int {
	if val == "" {
		return 0
	}
	num, _ := strconv.Atoi(val)
	return num
}

func Int642Str(val int64) string {
	if val == 0 {
		return ""
	}
	return strconv.Itoa(int(val))
}

func Str2Float64(param string) float64 {
	if param == "" {
		return 0
	}
	val, err := strconv.ParseFloat(param, 64)
	if err != nil {
		return 0
	}
	return val
}

// countChineseCharacters 计算字符串中的中文字符数量
func CountChineseCharacters(s string) int {
	count := 0
	for _, r := range s {
		if r >= 0x4E00 && r <= 0x9FA5 { // 判断是否为常用中文字符范围
			count++
		}
	}
	return count
}

// 检查经纬度合法
func CheckLoc(loc string) bool {
	data := strings.Split(loc, ",")
	if len(data) != 2 {
		return false
	}
	long := Str2Float64(data[0])
	lat := Str2Float64(data[1])
	if long > 180 || long < -180 {
		return false
	}
	if lat > 85 || lat < -85 {
		return false
	}

	return true
}
