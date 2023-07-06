package fund

import (
	cron "github.com/robfig/cron/v3"
)

// 定时
func InitCron() {
	list := fundList{}
	star := fundStar{}
	earnings := fundEarnings{}
	earningsRank := FundEaringsRank{}
	buySell := FundBuySell{}
	task := daysPastTimeRank{}

	c := cron.New()
	//收益排行 0点
	_, err := c.AddFunc("0 0 * * 2-6", earningsRank.GetData)
	if err != nil {
		panic("cron err :" + err.Error())
	}
	//基金购买情况 0点
	_, err = c.AddFunc("10 0 * * 2-6", buySell.GetData)
	if err != nil {
		panic("cron err :" + err.Error())
	}

	//star 每个月
	_, err = c.AddFunc("0 0 * */1 *", star.GetData)
	if err != nil {
		panic("cron err :" + err.Error())
	}

	//list 每个月
	_, err = c.AddFunc("0 0 * */1 *", list.GetData)
	if err != nil {
		panic("cron err :" + err.Error())
	}

	//earing 每周日
	_, err = c.AddFunc("0 0 * * 0", earnings.GetData)
	if err != nil {
		panic("cron err :" + err.Error())
	}

	//wechat
	//基金购买情况 0点
	_, err = c.AddFunc("30 0 * * 2-6", task.Send)
	if err != nil {
		panic("cron err :" + err.Error())
	}
	c.Start()
}
