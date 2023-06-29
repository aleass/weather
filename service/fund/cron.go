package fund

import (
	"github.com/go-co-op/gocron"
	"time"
)

// 定时
func InitCron() {
	var c = gocron.Scheduler{}
	timezone, _ := time.LoadLocation("Asia/Shanghai")
	gocron.NewScheduler(timezone)

	//日
	_, err := c.Cron("0 0 * * 1-5").Do(fundDayEarings)
	if err != nil {
		panic("cron err :" + err.Error())
	}

	//star
	_, err = c.Cron("0 0 1 * *").Do(GetStarData)
	if err != nil {
		panic("cron err :" + err.Error())
	}

	//earing
	_, err = c.Every(1).Sunday().Do(GetEaringsReq)
	if err != nil {
		panic("cron err :" + err.Error())
	}

}
