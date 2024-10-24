package service

import (
	"fmt"
	"time"
	"weather/api/atmp"
	"weather/api/gz_weather"
	"weather/api/he_feng"
	"weather/api/sysos"
	"weather/api/telegram"
	"weather/api/typhoon"
	"weather/common"
)

func RunWeather(sleepTimes time.Duration) {
	defer func() {
		if err := recover(); err != nil {
			common.LogSend(fmt.Sprintf("panic err:%v", err), common.PanicType)
		}
		time.Sleep(sleepTimes)
		go RunWeather(sleepTimes)
	}()

	var (
		lastUpdateMsg  string
		lastUpdateHour int64
		newAddr        = make(chan bool, 1)
		isNewAddr      bool
	)
	go isNewAddress(newAddr)

	for {
		var (
			sendMsg string
			weather string
			now     = time.Now()
		)

		//获取地址
		if !isNewAddr {
			isNewAddr = telegram.GetAddress()
		}

		//typhoon per hour
		typhoonMsg := typhoon.TyphoonActive()

		//获取天气 5分钟降雨
		var rainInfo = he_feng.FiveMinRain()

		//获取位置
		//loc, addr, ok := he_feng.Lookup()
		loc, addr, ok := atmp.SearchByLonLac(common.MyConfig.Atmp.Loc)

		//获取实时降雨量测试点
		var realData string
		if ok {
			realData = gz_weather.GZWeather(loc)
		}

		warningTitle, warningText := he_feng.CityWarning()

		weather = he_feng.WeatherInfo()

		//推送
		msg := warningTitle + rainInfo + typhoonMsg + warningText

		switch {
		//初次发送 地址变更 离上次更新时间大于一小时
		case lastUpdateMsg == "" || isNewAddr || now.Unix()-3600*5 > lastUpdateHour:
			sendMsg = weather + warningTitle + rainInfo + realData + typhoonMsg + warningText

		//消息变更
		case lastUpdateMsg != msg:
			sendMsg = warningTitle + rainInfo + realData + weather + typhoonMsg + warningText
		}

		if sendMsg != "" {
			lastUpdateHour = now.Unix()
			lastUpdateMsg = msg

			sendMsg += addr + now.Format(" 15:04 ") + sysos.OSPower
			telegram.SendMessage(sendMsg, common.MyConfig.Telegram.Token)
			common.Logger.Info(sendMsg)
		}

		var used = fmt.Sprintf("%s %s  MapApi:%d,WeatherApi:%d", addr, now.Format("15:04"), common.AmtpApiCount, common.HeFengApiCount)
		common.Logger.Info(used)
		isNewAddr = false

		//正常睡眠
		var curSleepTime = time.After(sleepTimes)
		var h = now.Hour()
		switch {
		case h < 6:
			curSleepTime = time.After(sleepTimes * 2)
		}

		select {
		case <-curSleepTime: //正常睡眠
			break
		case <-newAddr: //新地址马上触发
			isNewAddr = true
			break
		}
	}
}

// 定时检测新地址
func isNewAddress(newAddr chan bool) {
	for {
		time.Sleep(time.Minute)
		if telegram.GetAddress() {
			newAddr <- true
		}
	}
}
