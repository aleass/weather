package common

const (
	UsualTime = "2006-01-02"
	UsualDate = "2006-01-02 15:04:05"
)

type DaysPastTimeRank struct {
	Buy            string `gorm:"column:buy"`
	Code           string `gorm:"column:code"`
	Name           string `gorm:"column:name"`
	Past1Month     string `gorm:"column:past_1_month"`    //近一月
	Past3Months    string `gorm:"column:past_3_months"`   //近三月
	Past6Months    string `gorm:"column:past_6_months"`   //近六月
	Past1Year      string `gorm:"column:past_1_year"`     //近一年
	SinceInception string `gorm:"column:since_inception"` //成立至今
}
