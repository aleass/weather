package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"os"
	"weather/common"
)

type Config struct {
	// 结构映射
	Wechat []struct {
		Urls  string `mapstructure:"url"`
		Notes string `mapstructure:"note"`
	} `mapstructure:"wechat"`
	CaiYun struct {
		Token  string `json:"token"`
		Addres []struct {
			Name        string `json:"name"`
			WechatNotes string `json:"wechatNotes"`
			Coordinate  string `json:"coordinate"`
			Switch      bool   `json:"switch" desc:"开关"`
		} `json:"addres"`
	} `json:"caiyun"`
}

// 运行
func run() {
	var (
		taskMap  = map[string]*urlInfo{} //任务控制
		myConfig = Config{}
		vip      = viper.New()
		path     = "pkg/config.yaml"
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
		common.Logger.Error("err:" + err.Error())
	}
	//添加要监控的对象，文件或文件夹

	if err = watch.Add(path); err != nil {
		common.Logger.Error("err:" + err.Error())
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

		var notes = make(map[string]string, len(myConfig.Wechat))
		for _, v := range myConfig.Wechat {
			notes[v.Notes] = v.Urls
		}

		for _, v := range myConfig.CaiYun.Addres {
			info, ok := taskMap[v.Name]
			if !ok {
				//生成一个任务
				task := &urlInfo{
					name:      v.Name,
					caiyunUrl: fmt.Sprintf(caiYunUrl, myConfig.CaiYun.Token, v.Coordinate),
					weChatUrl: notes[v.WechatNotes],
					_switch:   make(chan struct{}),
					watchTime: 5, //默认10分钟
				}
				taskMap[v.Name] = task
				info = task
			}
			if !v.Switch && info.isrun {
				info._switch <- struct{}{} //关闭一个任务
				info.isrun = false
			} else if v.Switch && !info.isrun {
				info.isrun = true
				go watchWeather(info) //生成一个监控任务
			}
		}
		<-watch.Events //文件监控
	}
}
