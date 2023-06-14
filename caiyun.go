package main

import (
	"fmt"
	"strings"
	"time"
	"weather/common"
)

const (
	caiYunUrl = "https://api.caiyunapp.com/v2.6/%s/%s/weather?alert=true&dailysteps=1&hourlysteps=24&unit=metric:v2"
	wechatUrl = "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key="
)

type urlInfo struct {
	name      string        `desc:"地址"`
	caiYunUrl string        `desc:"caiyun url"`
	weChatUrl string        `desc:"wechat url"`
	_switch   chan struct{} `desc:"开关"`
	isRun     bool          `desc:"是否运行"`
	watchTime time.Duration `desc:"监控时间:分钟"` //
}

var msg = " %s\n\n温度(体感):%.1f(%.1f)C°\n紫外线:%s\n体感:%s\n空气质量:%s(%d)\n湿度:%.1f%%\n%s\n未来24小时天气:%s"

// 减少一瞬间请求
var delay = make(chan struct{}, 1)

// 天气监控
func watchWeather(info *urlInfo) {
	defer func() {
		if err := recover(); err != nil {
			common.LogSend(fmt.Sprintf("panic err:%v", err), common.PanicType)
			info.isRun = false
		}
	}()
	//默认轮询监控,时间频率为60分钟
	var (
		rainMsg, weatherMsg string
		_url                = info.caiYunUrl
		res                 *Weather
		err                 error
		lastTime            int64
		now                 time.Time
		_realtime           realtime
		tempStatus          string
		_windLevel          [2]string
		windStr, wind       string
		val                 float64
	)

	for {
		now = time.Now()
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
		res, err = getWeatherRawData(_url)
		if err != nil {
			common.LogSend(_url+":发生错误:"+err.Error(), common.ErrType)
			goto end
		}
		//防止并发请求
		time.Sleep(time.Second)
		<-delay

		if err != nil {
			common.LogSend(_url+":发生错误:"+err.Error(), common.ErrType)
			goto end
		}
		_realtime = res.Result.Realtime
		//发生了变化,减少时间,监控异常
		rainMsg = "\n"
		//	雨水
		switch _realtime.Skycon {
		case "LIGHT_RAIN", "MODERATE_RAIN", "HEAVY_RAIN", "STORM_RAIN":
			rainMsg = fmt.Sprintf("\n降水强度:%.1f毫米/小时,最近的降水带距离%.1f公里和降水强度%.1f毫米/小时,", _realtime.Precipitation.Local.Intensity,
				_realtime.Precipitation.Nearest.Distance, _realtime.Precipitation.Nearest.Intensity)
		}
		if index := strings.Index(res.Result.Minutely.Description, "还在加班么？注意休息哦"); index != -1 {
			weatherMsg = SkyconStatus[_realtime.Skycon] + rainMsg + res.Result.Minutely.Description[:index]
		} else {
			weatherMsg = SkyconStatus[_realtime.Skycon] + rainMsg + res.Result.Minutely.Description
		}

		// 风向
		val = (_realtime.Wind.Direction - 11.26) / 22.50
		if val < 0 {
			windStr = "北"
		} else {
			windStr = windDirection[int(val)]
		}
		_windLevel = (*windLevel[int(_realtime.Wind.Speed)])
		wind = fmt.Sprintf("%s%s风%s", _windLevel[0], windStr, _windLevel[1])

		//发送大于6小时才发生或天气发生变化
		if now.Unix()-lastTime >= 6*3600 || tempStatus != SkyconStatus[_realtime.Skycon] {
			//发送
			common.Send(now.Format("15:04:05 ")+info.name+
				fmt.Sprintf(msg, weatherMsg, _realtime.Temperature, _realtime.ApparentTemperature,
					_realtime.LifeIndex.Ultraviolet.Desc, _realtime.LifeIndex.Comfort.Desc, _realtime.AirQuality.Description.Chn,
					_realtime.AirQuality.Aqi.Chn, _realtime.Humidity*100, wind, res.Result.Hourly.Description), info.weChatUrl)

			//记录这次发送时间和信息
			lastTime = now.Unix()
			tempStatus = SkyconStatus[_realtime.Skycon]
		}
	end:
		time.Sleep(time.Minute * info.watchTime)
	}
}
