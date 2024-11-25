package service

import "services/common"

// 临时触发一次
func TemWeather() {
	for range NewTempAddr {
		var weather = NewWeather(&common.MyConfig.TemHome.Loc, &common.MyConfig.TemHome.Addr)
		weather.GetWeatcherInfo()
	}
}
