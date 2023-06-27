package model

type DfDateNote struct {
	Date int `gorm:"column:date"     desc:""`
}

func (DfDateNote) TableName() string {
	return "df_date_note"
}
