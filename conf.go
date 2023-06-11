package main

import (
	"fmt"
	"github.com/spf13/viper"
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

func init() {
	viper.SetConfigFile("pkg/config.yaml")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("无法读取配置文件: %w", err))
	}

	if err := viper.Unmarshal(&MyConfig); err != nil {
		panic(fmt.Errorf("无法解析配置文件: %w", err))
	}

	Location = make([]url_info, len(MyConfig.CaiYun.Addres))
	tempStatus = make([]string, len(MyConfig.CaiYun.Addres))

	var notes = make(map[string]string, len(MyConfig.Wechat))
	for _, v := range MyConfig.Wechat {
		notes[v.Notes] = v.Urls
	}

	for i, v := range MyConfig.CaiYun.Addres {
		if !v.Switch {
			continue
		}
		Location[i] = url_info{
			addr:      v.Name,
			caiyunUrl: fmt.Sprintf(caiyunUrl, MyConfig.CaiYun.Token, v.Coordinate),
			weChatUrl: notes[v.WechatNotes],
		}
	}

}
