package main

import (
	"weather/service"
	"weather/service/fund"
)

func main() {
	//初始化配置
	service.InitConfig()
	star := fund.FundEarnings{}
	star.GetData()
}
