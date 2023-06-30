package model

import "time"

type DfFundList struct {
	Code       string    `gorm:"column:code"        desc:"代码"`
	Name       string    `gorm:"column:name"        desc:"名字"`
	Pinyin     string    `gorm:"column:pinyin"      desc:"拼音"`
	AbbrPinyin string    `gorm:"column:abbr_pinyin" desc:"拼音简写"`
	Type       string    `gorm:"column:type"        desc:"基金类型"`
	Date       time.Time `gorm:"column:date"        desc:"日期"`
	Id         int64     `gorm:"column:id"          desc:""`
}

func (DfFundList) TableName() string {
	return "df_fund_list"
}
