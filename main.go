package main

import (
	"weather/service"
)

func main() {
	go service.HttpRun()
	//基金获取
	go service.FundRun()
	//获取配置并运行
	service.Run()
}
