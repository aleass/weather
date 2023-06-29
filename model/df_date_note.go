package model

type DfDateNote struct {
	Date int64  `gorm:"column:date" desc:""`
	Type string `gorm:"column:type"        desc:"类型"`
}

func (DfDateNote) TableName() string {
	return "df_date_note"
}
