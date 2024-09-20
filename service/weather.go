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

		//【天气】25° 阴 西北风10m/s 小时降雨：0.0 能见度：17
		var weather string

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
			realData = "\n" + gz_weather.GZWeather(loc)
		}
		var used = fmt.Sprintf("\n%s %s  MapApi：%d 	WeatherApi：%d\n\n", addr, now.Format("15:04"), common.AmtpApiCount, common.HeFengApiCount)

		warningTitle, warningText := he_feng.CityWarning()

		//推送
		msg := warningTitle + rainInfo + typhoonMsg + warningText

		//地址变更 离上次更新时间大于一小时
		if isNewAddr || now.Unix()-3600 > lastUpdateHour {
			lastUpdateHour = now.Unix()
			weather = he_feng.WeatherInfo()
		}

		//消息变更
		if lastUpdateMsg != msg || weather != "" {
			telegram.SendMessage(warningTitle+rainInfo+realData+weather+typhoonMsg+warningText, common.MyConfig.Telegram.Token)
			lastUpdateMsg = msg
		}

		common.Logger.Info(used)

		time.Sleep(selectTime)
	}
}
