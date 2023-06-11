package main

import (
	"fmt"
	"github.com/spf13/viper"
	"time"
)

var MyConfig = Config{}

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

func getConfig() {
	var (
		vip     = viper.New()
		taskMap = map[string]*url_info{}
	)
	vip.SetConfigFile("pkg/config.yaml")
	vip.SetConfigType("yaml")
	//循环更新
	for {
		if err := vip.ReadInConfig(); err != nil {
			panic(fmt.Errorf("无法读取配置文件: %w", err))
		}

		if err := vip.Unmarshal(&MyConfig); err != nil {
			panic(fmt.Errorf("无法解析配置文件: %w", err))
		}

		var notes = make(map[string]string, len(MyConfig.Wechat))
		for _, v := range MyConfig.Wechat {
			notes[v.Notes] = v.Urls
		}

		for _, v := range MyConfig.CaiYun.Addres {
			info, ok := taskMap[v.Name]
			if !ok {
				task := &url_info{
					addr:      v.Name,
					caiyunUrl: fmt.Sprintf(caiyunUrl, MyConfig.CaiYun.Token, v.Coordinate),
					weChatUrl: notes[v.WechatNotes],
					_switch:   v.Switch,
				}
				taskMap[v.Name] = task
				go watch_weather(task)
				info = task
				continue
			}
			if !v.Switch {
				info._switch = false
			}
		}
		time.Sleep(time.Minute * 10)
	}
}
