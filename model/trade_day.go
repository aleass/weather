package model


type TradeDay struct {
	Date              int64     `gorm:"column:date"                desc:""`
}

func (TradeDay) TableName() string {
	return "trade_day"
}
