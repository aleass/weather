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

type Users struct {
	Name        string `desc:"名字"`
	ConfigGroup map[string]*BaseWeatherInfo
}

type BaseWeatherInfo struct {
	Ip          string   `json:"ip"`
	Op          string   `json:"op" desc:"当前操作"`
	AddrCodes   string   `json:"adcodes" desc:"经纬度"`
	AllowNight  bool     `desc:"晚上是否运行运行"`
	Main        string   `desc:"配置信息"`
	AllowWeek   *[7]bool `desc:"不为空则指定星期运行"`
	StartTime   int64    `desc:"任务创建时间"`
	Msg         strings.Builder
	WatchTime   time.Duration `desc:"监控时间:分钟" json:"watch_time"`
	IsRun       bool          `desc:"是否运行" json:"is_run"`
	Name        string        `desc:"名字"`
	Notes       string        `desc:"备注"`
	Address     string        `desc:"url 配置的地址"`
	CaiYunUrl   string        `desc:"caiyun url" json:"cai_yun_url"`
	WeChatUrl   string        `desc:"wechat url" json:"we_chat_url"`
	Switch      bool
	IsUrlConfig bool
}

var (
	myConfig       = common.Config{}
	wechatNoteMap  = make(map[string]string, len(myConfig.Wechat))
	allowUrlConfig = make(map[string]string, len(myConfig.UrlConfigPass))
	taskMap        = map[string]*Users{} //任务控制

)

// 运行
func Run() {
	defer func() {
		if err := recover(); err != nil {
			stack := common.Stack(3)
			common.Logger.Error(string(stack))
			common.LogSend(fmt.Sprintf("panic err:%v", err), common.PanicType)
		}
		Run()
	}()
	var (
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

		//map
		common.SetToken(myConfig.QqMapToken, myConfig.GeoMapToken)

		for _, v := range myConfig.Wechat {
			wechatNoteMap[v.Notes] = v.Token
		}
		for _, v := range myConfig.UrlConfigPass {
			allowUrlConfig[v.Name] = v.Notes
		}

		for _, v := range myConfig.CaiYun.Addres {
			user, ok := taskMap[v.Name]
			if !ok {
				user = &Users{
					Name:        v.Name,
					ConfigGroup: map[string]*BaseWeatherInfo{},
				}
				taskMap[v.Name] = user
				go user.WatchWeather() //生成一个监控任务
			}

			config, ok := user.ConfigGroup[v.Addr]
			if !ok {
				config = getUrlInfo(v.Addr, v.Coordinate, v.Name, v.AllowWeek, 5, v.Switch)
				user.ConfigGroup[v.Addr] = config
			} else {
				updateUrlInfo(config, v.Addr, v.Coordinate, v.Name, v.AllowWeek, 5)
			}

			if !v.Switch && config.IsRun {
				config.IsRun = false
			}
		}
		common.ErrorUrl = wechatUrl + wechatNoteMap["error"]
		<-watch.Events //文件监控
	}
}

func updateUrlInfo(weather *BaseWeatherInfo, addr, coordinate, wechatNotes, allowWeek string, watchTime time.Duration) {
	weather.Address = addr
	weather.CaiYunUrl = fmt.Sprintf(caiYunUrl, myConfig.CaiYun.Token, coordinate)
	weather.WeChatUrl = wechatUrl + wechatNoteMap[wechatNotes]
	weather.WatchTime = watchTime //默认10分
	weather.Notes = wechatNotes
	if allowWeek != "" {
		weather.AllowWeek = &[7]bool{}
		for _, w := range strings.Split(allowWeek, ",") {
			week, err := strconv.Atoi(w)
			if err != nil || week < 0 || week > 6 {
				panic("invial  week")
			}
			weather.AllowWeek[week] = true
		}
	}
}

func getUrlInfo(addr, coordinate, name, allowWeek string, watchTime time.Duration, _switch bool) *BaseWeatherInfo {
	user := &BaseWeatherInfo{
		Address:   addr,
		CaiYunUrl: fmt.Sprintf(caiYunUrl, myConfig.CaiYun.Token, coordinate),
		WeChatUrl: wechatUrl + wechatNoteMap[name],
		Switch:    _switch,
		WatchTime: watchTime, //默认10分钟
		Notes:     name,
		StartTime: time.Now().Unix(),
	}
	if allowWeek != "" {
		user.AllowWeek = &[7]bool{}
		for _, w := range strings.Split(allowWeek, ",") {
			week, err := strconv.Atoi(w)
			if err != nil || week < 0 || week > 6 {
				panic("invial  week")
			}
			user.AllowWeek[week] = true
		}
	}
	return user
}
