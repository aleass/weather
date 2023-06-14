package main

import (
	"fmt"
	"sort"
	"strings"
	"time"
	"weather/common"
)

const (
	caiYunUrl = "https://api.caiyunapp.com/v2.6/%s/%s/weather?alert=true&dailysteps=1&hourlysteps=24&unit=metric:v2"
	wechatUrl = "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key="
	msg       = " %s\n●温度:%.1fC° 体感:%s(%.1fC°)\n●紫外线:%s AQI:%s(%d) 湿度:%.1f%%\n●%s\n●未来24小时天气:%s"
)

type urlInfo struct {
	name      string        `desc:"地址"`
	caiYunUrl string        `desc:"caiyun url"`
	weChatUrl string        `desc:"wechat url"`
	_switch   chan struct{} `desc:"开关"`
	isRun     bool          `desc:"是否运行"`
	watchTime time.Duration `desc:"监控时间:分钟"`
	msg       strings.Builder
}

// 减少一瞬间请求
var delay = make(chan struct{}, 1)

// 天气监控
func (info *urlInfo) watchWeather() {
	defer func() {
		if err := recover(); err != nil {
			common.LogSend(fmt.Sprintf("panic err:%v", err), common.PanicType)
			info.isRun = false
		}
	}()
	//默认轮询监控,时间频率为60分钟
	var (
		res                        *Weather
		err                        error
		lastTime                   int64
		realtime                   Realtime
		_weatherMsg, _alertMsg     string
		rainMsg, windMsg, alertMsg string
	)

	for {
		now := time.Now()
		select {
		case <-info._switch:
			println("任务退出：", info.name)
			return
		default:
		}
		//0点到6点 不发送
		if now.Hour() < 6 && now.Hour() > 0 {
			goto end
		}

		//并发控制
		delay <- struct{}{}
		now = time.Now()
		res, err = getWeatherRawData(info.caiYunUrl)
		if err != nil {
			common.LogSend(info.caiYunUrl+":发生错误:"+err.Error(), common.ErrType)
			goto end
		}
		//防止并发请求
		time.Sleep(time.Second)
		<-delay

		if err != nil {
			common.LogSend(info.caiYunUrl+":发生错误:"+err.Error(), common.ErrType)
			goto end
		}
		realtime = res.Result.Realtime

		//雨水
		rainMsg = info.getRainData(res, &realtime)

		// 风向
		windMsg = info.getWindData(&realtime)

		//预警
		if len(res.Result.Alert.Content) > 0 {
			alertMsg = info.getAlterData(res)
		}

		//发送大于6小时才发生 天气发生变化 预警变更(取消或新增,修改)
		if now.Unix()-lastTime >= 6*3600 || _weatherMsg != SkyconStatus[realtime.Skycon] || alertMsg != _alertMsg {
			//发送
			common.Send(now.Format("15:04:05 ")+info.name+
				fmt.Sprintf(msg, rainMsg, realtime.Temperature, realtime.LifeIndex.Comfort.Desc, realtime.ApparentTemperature,
					realtime.LifeIndex.Ultraviolet.Desc, realtime.AirQuality.Description.Chn,
					realtime.AirQuality.Aqi.Chn, realtime.Humidity*100, windMsg, res.Result.Hourly.Description)+alertMsg, info.weChatUrl)

			//记录这次发送时间和信息
			lastTime = now.Unix()
			_weatherMsg = SkyconStatus[realtime.Skycon]
			_alertMsg = alertMsg
			alertMsg = ""
		}
	end:
		time.Sleep(time.Minute * info.watchTime)
	}
}

// 雨水
func (info *urlInfo) getRainData(res *Weather, _realtime *Realtime) string {
	rainMsg := "\n●"
	var weatherMsg string
	//	雨水
	switch _realtime.Skycon {
	case "LIGHT_RAIN", "MODERATE_RAIN", "HEAVY_RAIN", "STORM_RAIN":
		rainMsg = fmt.Sprintf("\n●降水强度:%.1f毫米/小时,最近的降水带距离%.1f公里和降水强度%.1f毫米/小时,",
			_realtime.Precipitation.Local.Intensity, _realtime.Precipitation.Nearest.Distance, _realtime.Precipitation.Nearest.Intensity)
	}
	if index := strings.Index(res.Result.Minutely.Description, "还在加班么？注意休息哦"); index != -1 {
		weatherMsg = SkyconStatus[_realtime.Skycon] + rainMsg + res.Result.Minutely.Description[:index]
	} else {
		weatherMsg = SkyconStatus[_realtime.Skycon] + rainMsg + res.Result.Minutely.Description
	}
	return weatherMsg
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
		if windStr = windDirection[index]; windStr == "" {
			if _realtime.Wind.Direction >= UnusualWind[index] {
				index++
			} else {
				index--
			}
			windStr = windDirection[index]
		}
	}

	_windLevel := (*windLevel[int(_realtime.Wind.Speed)])
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
