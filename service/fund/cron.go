package fund

import (
	"github.com/go-co-op/gocron"
	"time"
)

// 定时
func InitCron() {
	list := fundList{}
	star := fundStar{}
	day := fundDayEarnings{}
	earnings := fundEarnings{}
	earningsRank := FundEaringsRank{}
	timezone, _ := time.LoadLocation("Asia/Shanghai")
	c := gocron.NewScheduler(timezone)

	//日收益 0点
	_, err := c.Cron("0 0 * * 2-6").Do(day.GetData)
	if err != nil {
		panic("cron err :" + err.Error())
	}

	//排行榜 0点
	_, err = c.Cron("0 0 * * 2-6").Do(earningsRank.GetData)
	if err != nil {
		panic("cron err :" + err.Error())
	}

	//star 每个月
	_, err = c.Cron("0 0 1 * *").Do(star.GetData)
	if err != nil {
		panic("cron err :" + err.Error())
	}

	//list 每个月
	_, err = c.Cron("0 0 1 * *").Do(list.GetData)
	if err != nil {
		panic("cron err :" + err.Error())
	}

	//earing 每周
	_, err = c.Every(1).Sunday().Do(earnings.GetData)
	if err != nil {
		panic("cron err :" + err.Error())
	}

}
