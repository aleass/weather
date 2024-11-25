package service

import (
	"fmt"
	"services/api/atmp"
	"services/api/gz_weather"
	"services/api/he_feng"
	"services/api/sysos"
	"services/api/telegram"
	"services/api/typhoon"
	"services/common"
	"time"
)

type Weather struct {
	Addr           *string
	Loc            *string
	LastUpdateUnix int64
	LastUpdateMsg  string
	LastUpdateHour int
	IsNewAddr      bool
}

func NewWeather(loc, addr *string) *Weather {
	var obj = &Weather{Loc: loc, Addr: addr}
	if loc == nil {
		obj.Loc = &common.MyConfig.Home.Loc
		obj.Addr = &common.MyConfig.Home.Addr
	}
	return obj
}

func (w *Weather) GetWeatcherInfo() {
	var (
		now     = time.Now()
		sendMsg string
		weather string
		addr    string
	)

	//typhoon per hour
	typhoonMsg := typhoon.TyphoonActive(*w.Loc)
	//获取天气 5分钟降雨
	var rainInfo = he_feng.FiveMinRain(*w.Loc)

	//获取位置
	//loc, addr, ok := he_feng.Lookup()
	loc, addr, ok := atmp.SearchAddrs(*w.Addr, *w.Loc)
	*w.Loc = loc

	//获取实时降雨量测试点
	var realData string
	if ok {
		realData = gz_weather.GZWeather(*w.Loc)
	}

	warningTitle, warningText := he_feng.CityWarning(*w.Loc)

	weather = he_feng.WeatherInfo(*w.Loc)

	//推送
	msg := addr + warningTitle + rainInfo + typhoonMsg + warningText

	switch {
	//初次发送 地址变更 离上次更新时间大于5小时
	case w.LastUpdateMsg == "" || w.IsNewAddr || (w.LastUpdateHour != now.Hour() && threeHour(now)) || now.Unix()-3600*5 > w.LastUpdateUnix:
		sendMsg = weather + warningTitle + rainInfo + realData + typhoonMsg + warningText

	//消息变更
	case w.LastUpdateMsg != msg:
		sendMsg = warningTitle + rainInfo + realData + weather + typhoonMsg + warningText
	}
	w.LastUpdateHour = now.Hour()
	w.LastUpdateMsg = msg

	if sendMsg != "" {
		w.LastUpdateUnix = now.Unix()
		sendMsg += addr + now.Format(" 15:04 ") + sysos.OSPower
		telegram.SendMessage(sendMsg, common.MyConfig.Telegram.WeatherToken)
		common.Logger.Info(sendMsg)
	}

	var used = fmt.Sprintf("%s %s  MapApi:%d,WeatherApi:%d", *w.Addr, now.Format("15:04"), common.AmtpApiCount, common.HeFengApiCount)
	common.Logger.Info(used)

}

// 正点天气
func threeHour(t time.Time) bool {
	switch t.Hour() {
	case 0, 6, 9, 12, 15, 18, 21:
		return true
	}
	return false
}
