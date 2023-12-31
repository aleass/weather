package service

import (
	"fmt"
	"sort"
	"strings"
	"time"
	"weather/common"
)

const (
	weatherMsg = "%s %s[%s,%s] %.1fC°" + //title 时间 地点[体感描述，天气描述] 当前温度
		"%s" + //下雨描述
		"\n%s" + //温度
		"\n●紫外线:%s AQI:%s(%d) 湿度:%.1f%%" +
		"\n●%s°" + //风向
		"\n●未来24小时天气:%s"
	timeout = 3600 * 24 //1 day
)

// 减少一瞬间请求
var delay = make(chan struct{}, 1)

// 天气监控
func (info *urlInfo) WatchWeather() {
	defer func() {
		if err := recover(); err != nil {
			common.LogSend(fmt.Sprintf("panic err:%v", err), common.PanicType)
			info.IsRun = false
		}
	}()
	info.IsRun = true
	//默认轮询监控,时间频率为60分钟
	var (
		res                        *Weather
		err                        error
		lastDate                   int
		isTimeTo                   = true
		realtime                   Realtime
		_weatherMsg, _alertMsg     string
		rainMsg, windMsg, alertMsg string
	)

	for {
		now := time.Now()
		select {
		case <-info.Switch:
			println("任务退出：", info.Name)
			info.IsRun = false
			return
		default:
		}

		switch {
		case info.IsUrlConfig && info.RunTime+timeout < now.Unix(): //url配置的限时
			delete(taskMap, info.Name)
			_msg := fmt.Sprintf("%s %s 操作:%s-%s ip:%s 坐标:%s ", time.Now().Format("2006-01-02 15:04:05"),
				info.Name, "del", info.Main, info.Ip, info.Adcodes)
			common.LogSend(_msg, common.InfoErrorType)
			return
		case info.AllowWeek != nil && !info.AllowWeek[now.Weekday()]: //week allow
			goto end
		case !info.AllowNight && now.Hour() < 6 && now.Hour() > 0: //除了手动设置，0点到6点 不发送
			goto end
		}

		//并发控制
		delay <- struct{}{}
		res, err = GetWeatherRawData(info.CaiYunUrl)
		if err != nil {
			common.LogSend(info.CaiYunUrl+":发生错误:"+err.Error(), common.ErrType)
			goto end
		}
		//防止并发请求
		time.Sleep(time.Second)
		<-delay

		if err != nil {
			common.LogSend(info.CaiYunUrl+":发生错误:"+err.Error(), common.ErrType)
			goto end
		}
		realtime = res.Result.Realtime

		//地址
		info.address = " "
		if info.IsUrlConfig && len(res.Result.Alert.Adcodes) > 0 {
			info.address += info.Main + ":"
			for _, adcode := range res.Result.Alert.Adcodes {
				info.address += adcode.Name
			}
		}

		//雨水
		rainMsg = info.getRainData(&realtime)

		// 风向
		windMsg = info.getWindData(&realtime)

		//预警
		if len(res.Result.Alert.Content) > 0 {
			alertMsg = info.getAlterData(res)
		}

		//触发的时间段
		switch now.Hour() {
		case 0, 6, 12, 18:
			isTimeTo = true
		}

		//发送大于6小时才发生 天气发生变化 预警变更(取消或新增,修改)
		if info.isEdit || isTimeTo && now.Hour() != lastDate || _weatherMsg != SkyconStatus[realtime.Skycon] || alertMsg != _alertMsg {
			temp := ""
			if temps := res.Result.Daily.Temperature; len(temps) > 0 {
				temp = fmt.Sprintf("●气温:%.1fC ~ %.1fC 体感:%.1fC", temps[0].Min, temps[0].Max, realtime.ApparentTemperature)
			} else {
				temp = fmt.Sprintf("●气温: 体感:%.1fC", realtime.ApparentTemperature)
			}

			//发送
			common.Send(fmt.Sprintf(weatherMsg, now.Format("15:04:05 "),
				info.Name+info.address, realtime.LifeIndex.Comfort.Desc, SkyconStatus[realtime.Skycon], realtime.Temperature,
				rainMsg,
				temp,
				realtime.LifeIndex.Ultraviolet.Desc, realtime.AirQuality.Description.Chn, realtime.AirQuality.Aqi.Chn, realtime.Humidity*100,
				windMsg,
				res.Result.Hourly.Description)+alertMsg, info.WeChatUrl)
			//记录这次发送时间和信息
			lastDate = now.Hour()
			_weatherMsg = SkyconStatus[realtime.Skycon]
			_alertMsg = alertMsg
			alertMsg = ""
		}
	end:
		isTimeTo = false
		time.Sleep(time.Minute * 10)
	}
}

// 雨水
func (info *urlInfo) getRainData(_realtime *Realtime) string {
	//	雨水
	switch _realtime.Skycon {
	case "LIGHT_RAIN", "MODERATE_RAIN", "HEAVY_RAIN", "STORM_RAIN":
		return fmt.Sprintf("\n●降水强度:%.1f毫米/小时,最近的降水带距离%.1f公里和降水强度%.1f毫米/小时,",
			_realtime.Precipitation.Local.Intensity, _realtime.Precipitation.Nearest.Distance, _realtime.Precipitation.Nearest.Intensity)
	}

	return ""
}

// 风向
func (info *urlInfo) getWindData(_realtime *Realtime) string {
	var windStr string
	// 风向
	val := (_realtime.Wind.Direction - 11.26) / 22.50
	index := int(val)
	if val < 0 || index == 0 {
		windStr = "北"
	} else {
		if windStr = WindDirection[index]; windStr == "" {
			if _realtime.Wind.Direction >= UnusualWind[index] {
				index++
			} else {
				index--
			}
			windStr = WindDirection[index]
		}
	}

	_windLevel := (*WindLevel[int(_realtime.Wind.Speed)])
	windStr = fmt.Sprintf("%s风 %s", windStr, _windLevel)
	return windStr
}

// 预警
func (info *urlInfo) getAlterData(res *Weather) string {
	info.msg.Reset()
	info.msg.WriteString("\n\n------------预警------------")

	sort.Slice(res.Result.Alert.Content, func(i, j int) bool {
		return res.Result.Alert.Content[i].Pubtimestamp < res.Result.Alert.Content[j].Pubtimestamp
	})

	for i, content := range res.Result.Alert.Content {
		if content.RequestStatus != "ok" {
			continue
		}
		if index := strings.Index(content.Title, "布"); index != -1 {
			content.Title = content.Title[index+len("布"):]
		}
		if i > 0 {
			info.msg.WriteString("\n------------")
		}
		info.msg.WriteString(fmt.Sprintf("\n●标题:%s\n●内容:%s\n●状态:%s\n●发布时间:%s\n●来源:%s",
			content.Title, content.Description, content.Status, time.Unix(content.Pubtimestamp, 0).
				Format("2006-01-02 15:04:05"), content.Source))
	}
	return info.msg.String()
}
