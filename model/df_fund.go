package model

type DfFund struct {
	Id     int64  `gorm:"column:id"     desc:""`
	Code   string `gorm:"column:code"   desc:"代码"`
	Name   string `gorm:"column:name"   desc:"名字"`
	Date   int    `gorm:"column:date"   desc:"日期"`
	Pinyin string `gorm:"column:pinyin" desc:""`
}

func (DfFund) TableName() string {
	return "df_fund"
}
