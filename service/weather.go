package service

import (
	"fmt"
	"time"
	"weather/api/atmp"
	"weather/api/gz_weather"
	"weather/api/he_feng"
	"weather/api/telegram"
	"weather/api/typhoon"
	"weather/common"
)

func RunWeather(selectTime time.Duration) {
	defer func() {
		if err := recover(); err != nil {
			common.LogSend(fmt.Sprintf("panic err:%v", err), common.PanicType)
		}
		time.Sleep(selectTime)
		go RunWeather(selectTime)
	}()

	var (
		lastUpdateMsg  string
		lastUpdateHour int64
	)

	for {
		now := time.Now()

		//获取地址
		var isNewAddr = telegram.GetAddress()

		//typhoon per hour
		typhoonMsg := typhoon.TyphoonActive()

		//获取天气 5分钟降雨
		var rainInfo = he_feng.FiveMinRain()

		//获取位置
		//loc, addr, ok := he_feng.Lookup()
		loc, addr, ok := atmp.SearchByLonLac(common.MyConfig.Atmp.Loc)

		//获取实时 y
		var realData string
		if ok {
			realData = gz_weather.GZWeather(loc)
		}

		warningTitle, warningText := he_feng.CityWarning()

		//推送
		msg := warningTitle + rainInfo + typhoonMsg + warningText

		var (
			sendMsg string
			weather string
		)

		//【天气】25° 阴 西北风10m/s 小时降雨：0.0 能见度：17
		if lastUpdateMsg == "" || isNewAddr || now.Unix()-3600*5 > lastUpdateHour {
			weather = he_feng.WeatherInfo()
		}

		//消息变更
		if lastUpdateMsg != msg {
			sendMsg = warningTitle + rainInfo + realData + weather + typhoonMsg + warningText
		} else if weather != "" { //地址变更 离上次更新时间大于一小时
			sendMsg = weather + warningTitle + rainInfo + realData + typhoonMsg + warningText
		}

		if sendMsg != "" {
			sendMsg += "\n" + now.Format(addr+" 15:04")
			lastUpdateHour = now.Unix()
			telegram.SendMessage(sendMsg, common.MyConfig.Telegram.Token)
			lastUpdateMsg = msg
		}

		var used = fmt.Sprintf("\n%s %s  MapApi：%d 	WeatherApi：%d\n\n", addr, now.Format("15:04"), common.AmtpApiCount, common.HeFengApiCount)
		common.Logger.Info(used)

		var h = now.Hour()
		switch {
		case h < 6:
			time.Sleep(selectTime)
		}

		time.Sleep(selectTime)
	}
}
