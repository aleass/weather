package common

const (
	UsualTime = "20060102"
	UsualDate = "2006-01-02 15:04:05"
)

// fund
const (
	//阶段收益
	EarningsUrl = "http://fund.eastmoney.com/data/rankhandler.aspx?op=ph&dt=kf&ft=zq&rs=&gs=0&sc=1nzf&st=desc&tabSubtype=,,,,,&pi=1&pn=30000&dx=1&v=0.6135069950706549"

	//排名和日收益
	TotalEarningsUrl = "https://api.fund.eastmoney.com/pinzhong/LJSYLZS?fundCode=%s&type=se&indexcode=000300" //dt=单位，m 1月 ,se 成立来
	RankUrl          = "https://fund.eastmoney.com/pingzhongdata/%s.js"

	//基金列表
	FundListUrl    = "http://fund.eastmoney.com/js/fundcode_search.js?v=20230630094933"
	FundBuySellUrl = "http://fund.eastmoney.com/Data/Fund_JJJZ_Data.aspx?page=1,20000"

	//星
	StarUrl = "http://fund.eastmoney.com/data/fundrating.html"
)
