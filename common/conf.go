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

	Atmp struct {
		Key string `mapstructure:"key"`
		Loc string `mapstructure:"loc"`
	} `mapstructure:"atmp"`

	Telegram struct {
		Token       string `mapstructure:"token"`
		AddresToken string `mapstructure:"addres_token"`
		ChatId      int64  `mapstructure:"chat_id"`
	} `mapstructure:"telegram"`
}

var (
	MyConfig = Config{}
)

func init() {
	var (
		vip  = viper.New()
		path = "/Users/tuski/code/src/weather/pkg/config.yaml"
	)
	// 使用 os.Stat 函数获取文件的信息
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
		panic("key not exist")
	}
}
