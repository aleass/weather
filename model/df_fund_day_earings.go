package model

type DfFundDayEarings struct {
	BuyStatus     string  `gorm:"column:buy_status"     desc:"申购状态"`
	Code          string  `gorm:"column:code"           desc:""`
	Date          int64   `gorm:"column:date"           desc:"日期"`
	DayIncreRate  float64 `gorm:"column:day_incre_rate" desc:"日增长率"`
	DayIncreVal   float64 `gorm:"column:day_incre_val"  desc:"日增长值"`
	Id            int64   `gorm:"column:id"             desc:""`
	Name          string  `gorm:"column:name"           desc:""`
	SellStatus    string  `gorm:"column:sell_status"    desc:"赎回状态"`
	ServiceCharge string  `gorm:"column:service_charge" desc:"手续费"`
	TotalNV       string  `gorm:"column:total_NV"       desc:"累计净值"`
	UnitNV        string  `gorm:"column:unit_NV"        desc:"单位净值"`
	Type          string  `gorm:"column:type"        desc:"类型"`
}

func (DfFundDayEarings) TableName() string {
	return "df_fund_day_earings"
}
