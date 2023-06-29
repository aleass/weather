package model

import "time"

type DfFundEarings struct {
	Id              int64     `gorm:"column:id"             desc:""`
	Code            string    `gorm:"column:code"              desc:"基金代码"`
	CumulativeNav   string    `gorm:"column:cumulative_nav"    desc:"累计净值"`
	DailyGrowthRate string    `gorm:"column:daily_growth_rate" desc:"日增长率"`
	Date            time.Time `gorm:"column:date"              desc:"日期"`
	Name            string    `gorm:"column:name"              desc:"基金简称"`
	NavPerUnit      string    `gorm:"column:nav_per_unit"      desc:"单位净值"`
	Past1Month      string    `gorm:"column:past_1_month"      desc:"近1个月增长率"`
	Past1Week       string    `gorm:"column:past_1_week"       desc:"近1周增长率"`
	Past1Year       string    `gorm:"column:past_1_year"       desc:"近1年增长率"`
	Past2Years      string    `gorm:"column:past_2_years"      desc:"近2年增长率"`
	Past3Months     string    `gorm:"column:past_3_months"     desc:"近3个月增长率"`
	Past3Years      string    `gorm:"column:past_3_years"      desc:"近3年增长率"`
	Past6Months     string    `gorm:"column:past_6_months"     desc:"近6个月增长率"`
	SinceInception  string    `gorm:"column:since_inception"   desc:"成立来增长率"`
	ThisYear        string    `gorm:"column:this_year"               desc:"今年来增长率"`
}

func (DfFundEarings) TableName() string {
	return "df_fund_earings"
}
