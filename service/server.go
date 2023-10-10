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
	Name        string `desc:"地址" json:"name" `
	Address     string `json:"address" desc:"url 配置的地址"`
	IsRun       bool   `desc:"是否运行" json:"is_run"`
	IsUrlConfig bool   `desc:"是否url配置" json:"is_url_config"`
}

type configInfo struct {
	IsUrlConfig bool   `desc:"是否url配置" json:"is_url_config"`
	Ip          string `json:"ip"`
	Op          string `json:"op" desc:"当前操作"`
	Adcodes     string `json:"adcodes" desc:"经纬度"`
	AllowNight  bool   `desc:"晚上是否运行运行"`
	Main        string `desc:"配置信息"`
}

type urlInfo struct {
	Name      string        `desc:"名字"`
	Notes     string        `desc:"备注"`
	address   string        `desc:"url 配置的地址"`
	CaiYunUrl string        `desc:"caiyun url" json:"cai_yun_url"`
	WeChatUrl string        `desc:"wechat url" json:"we_chat_url"`
	Switch    chan struct{} `desc:"开关" json:"__switch"`
	IsRun     bool          `desc:"是否运行" json:"is_run"`
	WatchTime time.Duration `desc:"监控时间:分钟" json:"watch_time"`
	msg       strings.Builder
	AllowWeek *[7]bool `desc:"不为空则指定星期运行"`
	isEdit    bool     `desc:"发送了修改"`
	RunTime   int64    `desc:"任务创建时间"`
	configInfo
}

var (
	MyConfig       = common.Config{}
	wechatNoteMap  = make(map[string]string, len(MyConfig.Wechat))
	allowUrlConfig = make(map[string]string, len(MyConfig.UrlConfigPass))
	taskMap        = map[string]*urlInfo{} //任务控制

)

func GetWechatUrl(note string) string {
	return wechatUrl + wechatNoteMap[note]
}

var (
	vip  = viper.New()
	path = "pkg/config.yaml"
)

// 运行
func InitConfig() {
	// 使用 os.Stat 函数获取文件的信息
	_, err := os.Stat(path)
	// 检查文件是否存在
	if os.IsNotExist(err) {
		path = "config.yaml"
	}
	vip.SetConfigFile(path)
	vip.SetConfigType("yaml")
	if err = vip.ReadInConfig(); err != nil {
		panic(fmt.Errorf("无法读取配置文件: %w", err))
	}

	if err = vip.Unmarshal(&MyConfig); err != nil {
		panic(fmt.Errorf("无法解析配置文件: %w", err))
	}
	//mysql
	InitMysql()

	//map
	common.SetToken(MyConfig.QqMapToken, MyConfig.GeoMapToken)
}

func Run() {
	defer func() {
		if err := recover(); err != nil {
			common.LogSend(fmt.Sprintf("panic err:%v", err), common.PanicType)
		}
		time.Sleep(time.Minute * 10)
		Run()
	}()

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
		InitConfig()
		//if len(MyConfig.CaiYun.Token) == 0 || len(MyConfig.CaiYun.Addres) == 0 || len(MyConfig.Wechat) == 0 {
		//	panic("token，任务，发送url为空")
		//}

		//
		//for _, v := range MyConfig.Wechat {
		//	wechatNoteMap[v.Notes] = v.Token
		//}
		//for _, v := range MyConfig.UrlConfigPass {
		//	allowUrlConfig[v.Name] = v.Notes
		//}
		//
		//for _, v := range MyConfig.CaiYun.Addres {
		//	info, ok := taskMap[v.Addr]
		//	if !ok {
		//		//生成一个任务
		//		task := getUrlInfo(v.Addr, v.Coordinate, v.WechatNotes, v.AllowWeek, 5)
		//		taskMap[v.Addr] = task
		//		info = task
		//	} else {
		//		updateUrlInfo(info, v.Addr, v.Coordinate, v.WechatNotes, v.AllowWeek, 5)
		//	}
		//	if !v.Switch && info.IsRun {
		//		info.Switch <- struct{}{} //关闭一个任务
		//		info.IsRun = false
		//	} else if v.Switch && !info.IsRun {
		//		go info.WatchWeather() //生成一个监控任务
		//	}
		//}
		//common.ErrorUrl = wechatUrl + wechatNoteMap["error"]
		//healUrl = wechatUrl + wechatNoteMap["mine"]
		<-watch.Events //文件监控
	}
}

func updateUrlInfo(user *urlInfo, name, coordinate, wechatNotes, allowWeek string, watchTime time.Duration) {
	user.Name = name
	user.CaiYunUrl = fmt.Sprintf(caiYunUrl, MyConfig.CaiYun.Token, coordinate)
	user.WeChatUrl = wechatUrl + wechatNoteMap[wechatNotes]
	user.WatchTime = watchTime //默认10分
	user.Notes = wechatNotes
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
}

func getUrlInfo(name, coordinate, wechatNotes, allowWeek string, watchTime time.Duration) *urlInfo {
	user := &urlInfo{
		Name:      name,
		CaiYunUrl: fmt.Sprintf(caiYunUrl, MyConfig.CaiYun.Token, coordinate),
		WeChatUrl: wechatUrl + wechatNoteMap[wechatNotes],
		Switch:    make(chan struct{}, 1),
		WatchTime: watchTime, //默认10分钟
		Notes:     wechatNotes,
		RunTime:   time.Now().Unix(),
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
