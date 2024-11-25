package service

import (
	"services/api/telegram"
	"services/common"
	"time"
)

var NewAddr = make(chan bool, 1)
var NewTempAddr = make(chan bool, 1)

func RunWeather(sleepTimes time.Duration) {
	defer common.RecoverWithStackTrace(RunWeather, sleepTimes)
	var (
		weather = NewWeather(nil, nil)
	)

	go isNewAddress()
	go TemWeather()

	for {
		var now = time.Now()
		//获取地址
		if !weather.IsNewAddr {
			telegram.GetMessage()
		}

		weather.GetWeatcherInfo()
		weather.IsNewAddr = false

		//正常睡眠
		var curSleepTime = time.After(sleepTimes)
		var h = now.Hour()
		switch {
		case h < 6:
			//curSleepTime = time.After(sleepTimes * 2)
		}

		select {
		case <-curSleepTime: //正常睡眠
			break
		case <-NewAddr: //新地址马上触发
			weather.IsNewAddr = true
			break
		}
	}
}

// 定时检测新地址
func isNewAddress() {
	for {
		time.Sleep(5 * time.Minute)
		if ok, isTem := telegram.GetMessage(); ok {
			if isTem {
				NewTempAddr <- true
			} else {
				NewAddr <- true
			}
		}
	}
}
