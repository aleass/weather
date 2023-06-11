package main

import (
	"fmt"
	"strings"
	"time"
)

const (
	caiyunUrl = "https://api.caiyunapp.com/v2.6/%s/%s/weather?alert=true&dailysteps=1&hourlysteps=24&unit=metric:v2"
)

type url_info struct {
	addr      string //地址
	caiyunUrl string //caiyun url
	weChatUrl string //wechat url
	_switch   bool
}

var msg = ":%s\r\n当前温度:%.1fC°,体感温度:%.1fC,紫外线:%s,体感:%s,空气质量:%s(%d),未来24小时天气:%s"

// 减少一瞬间请求
var delay = make(chan struct{}, 1)

// 天气监控
func watch_weather(info *url_info) {
	//默认轮询监控,时间频率为60分钟

	var (
		rain_msg, weathereStatus string
		_url                     = info.caiyunUrl
		res                      *Weather
		err                      error
		lastTime                 int64
		now                      time.Time
		_realtime                realtime
		tempStatus               string
	)

	for {
		now = time.Now()
		if !info._switch {
			goto end
		}
		delay <- struct{}{}
		now = time.Now()
		//0点到6点 不发送
		if now.Hour() < 6 && now.Hour() > 0 {
			goto end
		}
		res, err = get_data(_url)
		if err != nil {
			Send(_url+":发生错误:"+err.Error(), info.weChatUrl)
			goto end
		}
		_realtime = res.Result.Realtime
		//发生了变化,减少时间,监控异常
		rain_msg = ","
		switch _realtime.Skycon {
		case "LIGHT_RAIN", "MODERATE_RAIN", "HEAVY_RAIN", "STORM_RAIN":
			rain_msg = fmt.Sprintf("(当前降水强度:%.1f毫米/小时,最近的降水带距离%.1f公里和降水强度%.1f毫米/小时),", _realtime.Precipitation.Local.Intensity,
				_realtime.Precipitation.Nearest.Distance, _realtime.Precipitation.Nearest.Intensity)
		}
		if index := strings.Index(res.Result.Minutely.Description, "还在加班么？注意休息哦"); index != -1 {
			weathereStatus = SkyconStatus[_realtime.Skycon] + rain_msg + res.Result.Minutely.Description[:index]
		} else {
			weathereStatus = SkyconStatus[_realtime.Skycon] + rain_msg + res.Result.Minutely.Description
		}
		//发送大于一小时才发生或天气发生变化
		if now.Unix()-lastTime >= 3600 || tempStatus != weathereStatus {
			//发送
			Send(now.Format("15:04:05 ")+info.addr+
				fmt.Sprintf(msg, weathereStatus, _realtime.Temperature, _realtime.ApparentTemperature,
					_realtime.LifeIndex.Ultraviolet.Desc, _realtime.LifeIndex.Comfort.Desc, _realtime.AirQuality.Description.Chn,
					_realtime.AirQuality.Aqi.Chn, res.Result.Hourly.Description), info.weChatUrl)
			//记录这次发送时间
			lastTime = now.Unix()
			tempStatus = weathereStatus
		}
		//更新状态
		//防止并发请求
		time.Sleep(time.Second)
		<-delay
	end:
		time.Sleep(time.Minute * 10)
	}
}
