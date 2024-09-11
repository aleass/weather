package service

import (
	"fmt"
	"time"
	"weather/api/gz_weather"
	"weather/api/he_feng"
	"weather/api/typhoon"
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
	for {
		now := time.Now()
		//获取地址
		GetAddress()

		weather := he_feng.WeatherInfo()

		//typhoon per hour
		typhoonMsg := typhoon.TyphoonActive()

		//获取位置
		loc, addr, ok := he_feng.Lookup()
		//loc, addr, ok := atmp.SearchByLonLac(common.MyConfig.Atmp.Loc)
		//获取天气
		var rainInfo = he_feng.FiveMinRain()
		//获取实时
		if ok {
			rainInfo += "\n" + gz_weather.GZWeather(loc)
		}
		var used = fmt.Sprintf("\n%s %s  MapApi：%d 	WeatherApi：%d\n\n", addr, now.Format("15:04"), common.AmtpApiCount, common.HeFengApiCount)

		warningTitle, warningText := he_feng.CityWarning()

		//推送
		SendMessage(warningTitle + weather + rainInfo + typhoonMsg + warningText + used)
		common.Logger.Info(warningTitle + weather + rainInfo + typhoonMsg + warningText + used)
		time.Sleep(time.Minute * 30)

		//归0

		common.AmtpApiCount, common.HeFengApiCount = 0, 0
	}
}
