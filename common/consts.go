package common

const (
	UsualTime     = "2006-01-02"
	UsualTimeInt  = "20060102"
	UsualTimeHour = "15:04:05"
	UsualDate     = "2006-01-02 15:04:05"
)

// fund
const (
	//阶段收益
	EarningsUrl = "http://fund.eastmoney.com/data/rankhandler.aspx?op=ph&dt=kf&ft=zq&rs=&gs=0&sc=1nzf&st=desc&tabSubtype=,,,,,&pi=1&pn=30000&dx=1&v=0.6135069950706549"

	//排名和日收益
	TotalEarningsUrl = "https://api.fund.eastmoney.com/pinzhong/LJSYLZS?fundCode=%s&type=se&indexcode=000300"
	//基金成立以来数据
	RankUrl          = "https://fund.eastmoney.com/pingzhongdata/%s.js"

	//基金列表
	FundListUrl    = "http://fund.eastmoney.com/js/fundcode_search.js?v=20230630094933"
	FundBuySellUrl = "http://fund.eastmoney.com/Data/Fund_JJJZ_Data.aspx?page=1,20000"

	//星
	StarUrl = "http://fund.eastmoney.com/data/fundrating.html"
)

// db
const (
	//每日基金排名
	DaysPastTimeRankSql = "SELECT e.code,e.name, past_1_month ,past_3_months ,past_6_months ,past_1_year " +
		",Inc_date,buy,since_inception  FROM `fund`.`df_fund_earnings` e INNER JOIN  `df_fund_list` l  on e.code = l.code " +
		"and  type in ('债券型-长债','债券型-中短债') and buy ='开放申购' where past_1_month >= 0.3 and past_3_months >= 1.5 " +
		"and past_6_months >= 3 and (past_1_year >= 6 or past_1_year = 0) order by past_1_month desc"
)

type DaysPastTimeRank struct {
	Buy            string `gorm:"column:buy"`
	Code           string `gorm:"column:code"`
	Name           string `gorm:"column:name"`
	Past1Month     string `gorm:"column:past_1_month"`    //近一月
	Past3Months    string `gorm:"column:past_3_months"`   //近三月
	Past6Months    string `gorm:"column:past_6_months"`   //近六月
	Past1Year      string `gorm:"column:past_1_year"`     //近一年
	SinceInception string `gorm:"column:since_inception"` //成立至今
}
