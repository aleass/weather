package fund

import (
	"github.com/go-co-op/gocron"
	"time"
)

// 定时
func InitCron() {
	list := fundList{}
	star := fundStar{}
	earnings := fundEarnings{}
	earningsRank := FundEaringsRank{}
	timezone, _ := time.LoadLocation("Asia/Shanghai")
	c := gocron.NewScheduler(timezone)
	//收益排行 0点
	_, err := c.Every(1).Day().At("00:00").Do(earningsRank.GetData)
	if err != nil {
		panic("cron err :" + err.Error())
	}

	//star 每个月
	_, err = c.Every(1).Month().At("00:00").Do(star.GetData)
	if err != nil {
		panic("cron err :" + err.Error())
	}

	//list 每个月
	_, err = c.Every(1).Month().At("01:00").Do(list.GetData)
	if err != nil {
		panic("cron err :" + err.Error())
	}

	//earing 每周
	_, err = c.Every(1).Sunday().Do(earnings.GetData)
	if err != nil {
		panic("cron err :" + err.Error())
	}
	c.StartAsync()
}
