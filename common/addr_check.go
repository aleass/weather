package common

import "strings"

// 检查参数
func CheckAddrOrLoc(data string) {
	if strings.Index(data, ".") != -1 {
		MyConfig.Home.Loc = data
		MyConfig.Home.Addr = ""
		return
	}
	MyConfig.Home.Addr = data
}
