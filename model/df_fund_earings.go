package model

type DfFundEarings struct {
	Date          int    `gorm:"column:date"             desc:""`
	Code          string `gorm:"column:code"           desc:""`
	Name          string `gorm:"column:name"           desc:""`
	UnitNV        string `gorm:"column:unit_NV"        desc:"单位净值"`
	TotalNV       string `gorm:"column:total_NV"       desc:"累计净值"`
	DayIncreVal   string `gorm:"column:day_incre_val"  desc:"日增长值"`
	DayIncreRate  string `gorm:"column:day_incre_rate" desc:"日增长率"`
	BuyStatus     string `gorm:"column:buy_status"     desc:"申购状态"`
	SellStatus    string `gorm:"column:sell_status"    desc:"赎回状态"`
	ServiceCharge string `gorm:"column:service_charge" desc:"手续费"`
}

func (DfFundEarings) TableName() string {
	return "df_fund_earings"
}
