package service

import (
	"services/common"
	"time"
)

// 运行
func Run() {
	defer common.RecoverWithStackTrace(RunWeather, 0)
	//10000ms
	//go sysos.GetPowermetrics("600000")
	go NewsRun(time.Minute * 30)
	go WebService()
	go RunWeather(time.Minute * 10)
	select {}
}
