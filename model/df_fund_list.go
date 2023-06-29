package model

type DfFundList struct {
	Code   string `gorm:"column:code"   desc:"代码"`
	Date   int64  `gorm:"column:date"   desc:"日期"`
	Id     int64  `gorm:"column:id"     desc:""`
	Name   string `gorm:"column:name"   desc:"名字"`
	Pinyin string `gorm:"column:pinyin" desc:""`
}

func (DfFundList) TableName() string {
	return "df_fund_list"
}
