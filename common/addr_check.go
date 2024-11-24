package common

import "strings"

// 检查参数
func CheckAddrOrLoc(data string) (loc, addr string) {
	if strings.Index(data, ".") != -1 {
		loc = data
		addr = ""
		return
	}
	addr = data
	return
}
