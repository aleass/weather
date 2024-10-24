package service

import (
	"fmt"
	"time"
	"weather/api/sysos"
	"weather/common"
)

// 运行
func Run() {
	defer func() {
		if err := recover(); err != nil {
			common.LogSend(fmt.Sprintf("panic err:%v", err), common.PanicType)
		}
		time.Sleep(time.Minute * 10)
		Run()
	}()

	//10000ms
	go sysos.GetPowermetrics("600000")

	go NewsRun(time.Minute * 30)
	go RunWeather(time.Minute * 30)
	select {}
}
