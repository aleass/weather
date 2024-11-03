package service

import (
	"services/api/sysos"
	"services/common"
	"time"
)

// 运行
func Run() {
	defer common.RecoverWithStackTrace(RunWeather, 0)
	//10000ms
	go sysos.GetPowermetrics("600000")
	go WebService()
	//go NewsRun(time.Minute * 30)
	go RunWeather(time.Minute * 30)
	select {}
}
