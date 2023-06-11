package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
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
		taskMap  = map[string]*url_info{} //任务控制
		myConfig = Config{}
		vip      = viper.New()
		path     = "pkg/config.yaml"
	)
	vip.SetConfigFile(path)
	vip.SetConfigType("yaml")

	//创建一个监控对象
	watch, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	//添加要监控的对象，文件或文件夹
	err = watch.Add(path)
	if err != nil {
		log.Fatal(err)
	}
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
				task := &url_info{
					name:      v.Name,
					caiyunUrl: fmt.Sprintf(caiyunUrl, myConfig.CaiYun.Token, v.Coordinate),
					weChatUrl: notes[v.WechatNotes],
					_switch:   make(chan struct{}),
				}
				taskMap[v.Name] = task
				info = task
			}
			if !v.Switch && info.isrun {
				info._switch <- struct{}{} //关闭一个任务
				info.isrun = false
			} else if v.Switch && !info.isrun {
				info.isrun = true
				go watch_weather(info) //生成一个监控任务
			}
		}
		<-watch.Events //文件监控
	}
	watch.Close()
}
