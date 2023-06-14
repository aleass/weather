package main

import (
	"weather/service"
)

func main() {
	go service.HttpRun()
	//获取配置并运行
	service.Run()
}
