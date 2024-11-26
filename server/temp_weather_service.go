package service

import "services/common"

// 临时触发一次
func TemWeather() {
	for addr := range NewTempAddr {
		_addr, _loc := common.CheckAddrOrLoc(addr)
		var weather = NewWeather(&_addr, &_loc)
		weather.GetWeatcherInfo()
	}
}
