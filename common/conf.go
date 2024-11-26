package common

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
)

type Config struct {
	// 结构映射
	HeFeng struct {
		Key string `mapstructure:"key"`
	} `mapstructure:"hefeng"`

	JuHe struct {
		Key string `mapstructure:"key"`
	} `mapstructure:"juhe"`

	Tian struct {
		Key string `mapstructure:"key"`
	} `mapstructure:"tian"`

	Atmp struct {
		Key string `mapstructure:"key"`
	} `mapstructure:"atmp"`

	Home struct {
		Loc  string `mapstructure:"loc"`
		Addr string `mapstructure:"addr"`
	} `mapstructure:"home"`

	System struct {
		RootPath   string `mapstructure:"root_path"`
		IsProxy    bool   `mapstructure:"is_proxy"`
		WatchPower bool   `mapstructure:"watch_power"`
	} `mapstructure:"system"`

	Telegram struct {
		WeatherToken string `mapstructure:"weather_token"`
		NewToken     string `mapstructure:"new_token"`
		AddresToken  string `mapstructure:"addres_token"`
		ChatId       int64  `mapstructure:"chat_id"`
	} `mapstructure:"telegram"`
}

var (
	MyConfig = Config{}
)

func init() {
	var (
		vip  = viper.New()
		path = FileKeyPath + "pkg/config.yaml"
	)

	// 使用 sysos.Stat 函数获取文件的信息
	_, err := os.Stat(path)
	// 检查文件是否存在
	if os.IsNotExist(err) {
		path = "config.yaml"
	}
	vip.SetConfigFile(path)
	vip.SetConfigType("yaml")
	vip.ReadInConfig()
	if err = vip.Unmarshal(&MyConfig); err != nil {
		panic(fmt.Errorf("无法解析配置文件: %w", err))
	}
	key := MyConfig.Atmp.Key
	if key == "" {
		wd, _ := os.Getwd()
		panic("key not exist:" + wd)
	}

	if MyConfig.System.RootPath != "" {
		FileKeyPath = MyConfig.System.RootPath
	}
}
