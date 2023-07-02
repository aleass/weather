package main

import (
	"weather/service"
	"weather/service/fund"
)

func main() {
	go service.HttpRun()
	go fund.InitCron()
	//获取配置并运行
	service.Run()
}
