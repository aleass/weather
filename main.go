package main

import (
	"weather/service"
	"weather/service/fund"
)

type DfFundList struct {
	AbbrPinyin string `gorm:"column:abbr_pinyin" desc:"拼音简写"`
	Code       string `gorm:"column:code"        desc:"代码"`
	Name       string `gorm:"column:name"        desc:"名字"`
	Pinyin     string `gorm:"column:pinyin"      desc:"拼音"`
	Type       string `gorm:"column:type"        desc:"基金类型"`
}

func main() {
	go service.HttpRun()
	go fund.InitCron()
	//获取配置并运行
	service.Run()
}
