package model

import "time"

type FundEaringsRank struct {
	Id          int64     `gorm:"column:id"                desc:""`
	Name        string    `gorm:"column:name"              desc:"基金简称"`
	Code        string    `gorm:"column:code"              desc:"基金代码"`
	Date        string    `gorm:"column:date"              desc:"日期"`
	Rank        string    `gorm:"column:rank"      desc:"排名"`
	RankPrecent string    `gorm:"column:rank_precent"      desc:"排名百分比"`
	TotalRate   string    `gorm:"column:total_rate"      desc:"排名百分比"`
	KindAvgRate string    `gorm:"column:kind_avg_rate"      desc:"排名百分比"`
	CreateTime  time.Time `gorm:"column:create_time"      desc:"创建时间"`
}

func (FundEaringsRank) TableName() string {
	return "fund_earings_rank"
}
