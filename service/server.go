package service

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"os"
	"strings"
	"time"
	"weather/common"
)

const (
	wechatUrl = "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key="
	caiYunUrl = "https://api.caiyunapp.com/v2.6/%s/%s/weather?alert=true&dailysteps=1&hourlysteps=24&unit=metric:v2"
)

type UrlInfo struct {
	name        string        `desc:"地址"`
	caiYunUrl   string        `desc:"caiyun url"`
	weChatUrl   string        `desc:"wechat url"`
	_switch     chan struct{} `desc:"开关"`
	isRun       bool          `desc:"是否运行"`
	watchTime   time.Duration `desc:"监控时间:分钟"`
	msg         strings.Builder
	isUrlConfig bool `desc:"是否url配置"`
}

var myConfig = common.Config{}
var notes = make(map[string]string, len(myConfig.Wechat))

// 运行
func Run() {
	var (
		taskMap = map[string]*UrlInfo{} //任务控制

		vip  = viper.New()
		path = "pkg/config.yaml"
	)
	// 使用 os.Stat 函数获取文件的信息
	_, err := os.Stat(path)
	// 检查文件是否存在
	if os.IsNotExist(err) {
		path = "config.yaml"
	}
	vip.SetConfigFile(path)
	vip.SetConfigType("yaml")

	//创建一个监控对象
	watch, err := fsnotify.NewWatcher()
	if err != nil {
		common.LogSend("error:"+err.Error(), common.ErrType)
	}
	//添加要监控的对象，文件或文件夹

	if err = watch.Add(path); err != nil {
		common.LogSend("error:"+err.Error(), common.ErrType)
	}
	defer watch.Close()
	for {
		if err = vip.ReadInConfig(); err != nil {
			panic(fmt.Errorf("无法读取配置文件: %w", err))
		}

		if err = vip.Unmarshal(&myConfig); err != nil {
			panic(fmt.Errorf("无法解析配置文件: %w", err))
		}

		if len(myConfig.CaiYun.Token) == 0 || len(myConfig.CaiYun.Addres) == 0 || len(myConfig.Wechat) == 0 {
			panic("token，任务，发送url为空")
		}

		for _, v := range myConfig.Wechat {
			notes[v.Notes] = v.Token
		}

		for _, v := range myConfig.CaiYun.Addres {
			info, ok := taskMap[v.Name]
			if !ok {
				//生成一个任务
				task := &UrlInfo{
					name:      v.Name,
					caiYunUrl: fmt.Sprintf(caiYunUrl, myConfig.CaiYun.Token, v.Coordinate),
					weChatUrl: wechatUrl + notes[v.WechatNotes],
					_switch:   make(chan struct{}),
					watchTime: 5, //默认10分钟
				}
				taskMap[v.Name] = task
				info = task
			}
			if !v.Switch && info.isRun {
				info._switch <- struct{}{} //关闭一个任务
				info.isRun = false
			} else if v.Switch && !info.isRun {
				info.isRun = true
				go info.WatchWeather() //生成一个监控任务
			}
		}
		common.ErrorUrl = wechatUrl + notes["error"]
		<-watch.Events //文件监控
	}
}
