package service

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"os"
	"strconv"
	"strings"
	"time"
	"weather/common"
)

const (
	wechatUrl = "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key="
	caiYunUrl = "https://api.caiyunapp.com/v2.6/%s/%s/weather?alert=true&dailysteps=1&hourlysteps=24&unit=metric:v2"
)

type UrlInfo struct {
	Name    string `desc:"地址" json:"name" `
	Address string `json:"address" desc:"url 配置的地址"`
	IsRun   bool   `desc:"是否运行" json:"is_run"`
}

type configInfo struct {
	IsUrlConfig bool   `desc:"是否url配置" json:"is_url_config"`
	Ip          string `json:"ip"`
	Op          string `json:"op" desc:"当前操作"`
	Adcodes     string `json:"adcodes" desc:"经纬度"`
	AllowNight  bool   `desc:"晚上是否运行运行"`
}

type urlInfo struct {
	Name        string        `desc:"名字"`
	address     string        `desc:"url 配置的地址"`
	CaiYunUrl   string        `desc:"caiyun url" json:"cai_yun_url"`
	WeChatUrl   string        `desc:"wechat url" json:"we_chat_url"`
	Switch      chan struct{} `desc:"开关" json:"__switch"`
	IsRun       bool          `desc:"是否运行" json:"is_run"`
	WatchTime   time.Duration `desc:"监控时间:分钟" json:"watch_time"`
	msg         strings.Builder
	IsUrlConfig bool `desc:"是否url配置" json:"is_url_config"`
}

var (
	myConfig       = common.Config{}
	wechatNoteMap  = make(map[string]string, len(myConfig.Wechat))
	allowUrlConfig = make(map[string]string, len(myConfig.UrlConfigPass))
)

// 运行
func Run() {
	var (
		taskMap = map[string]*urlInfo{} //任务控制
		vip     = viper.New()
		path    = "pkg/config.yaml"
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
			wechatNoteMap[v.Notes] = v.Token
		}
		for _, v := range myConfig.UrlConfigPass {
			allowUrlConfig[v.Name] = v.Notes
		}

		for _, v := range myConfig.CaiYun.Addres {
			info, ok := taskMap[v.Name]
			if !ok {
				//生成一个任务
				task := getUrlInfo(v.Name, v.Coordinate, v.WechatNotes, v.AllowWeek, 5)
				taskMap[v.Name] = task
				info = task
			}
			if !v.Switch && info.IsRun {
				info.Switch <- struct{}{} //关闭一个任务
				info.IsRun = false
			} else if v.Switch && !info.IsRun {
				go info.WatchWeather() //生成一个监控任务
			}
		}
		common.ErrorUrl = wechatUrl + wechatNoteMap["error"]
		<-watch.Events //文件监控
	}
}

func getUrlInfo(name, coordinate, wechatNotes, allowWeek string, watchTime time.Duration) *urlInfo {
	info := &urlInfo{
		Name:      name,
		CaiYunUrl: fmt.Sprintf(caiYunUrl, myConfig.CaiYun.Token, coordinate),
		WeChatUrl: wechatUrl + wechatNoteMap[wechatNotes],
		Switch:    make(chan struct{}, 1),
		WatchTime: watchTime, //默认10分钟
	}
	return info
}
