package fund

import (
	"fmt"
	"go.uber.org/zap/buffer"
	"time"
	"weather/common"
	"weather/service"
)

const (
	//每日基金排名
	daysPastTimeRankSql = "SELECT e.code,e.name, past_1_month ,past_3_months ,past_6_months ,past_1_year " +
		",buy,since_inception  FROM `fund`.`df_fund_earnings` e INNER JOIN  `df_fund_list` l  on e.code = l.code " +
		"and  type in ('债券型-长债','债券型-中短债') and buy ='开放申购' where past_1_month >= 0.3 and past_3_months >= 1.5 " +
		"and past_6_months >= 3 and (past_1_year >= 6 or past_1_year = 0) order by past_1_month desc"
)

type daysPastTimeRank struct {
	Buy            string  `gorm:"column:buy"`
	Code           string  `gorm:"column:code"`
	Name           string  `gorm:"column:name"`
	Past1Month     float64 `gorm:"column:past_1_month"`    //近一月
	Past3Months    float64 `gorm:"column:past_3_months"`   //近三月
	Past6Months    float64 `gorm:"column:past_6_months"`   //近六月
	Past1Year      float64 `gorm:"column:past_1_year"`     //近一年
	SinceInception float64 `gorm:"column:since_inception"` //成立至今
}

func (s *daysPastTimeRank) Send() {
	var list []daysPastTimeRank
	service.FuncDb.Raw(daysPastTimeRankSql).Find(&list)
	var msg = buffer.Buffer{}
	msg.WriteString("code name 近一月 近三月 近六月 近一年 成立至今\n")
	for _, info := range list {
		msg.WriteString(fmt.Sprintf("[%s,%s]%s %s %0.2f %0.2f %0.2f %0.2f  %0.2f \n", info.Buy, info.Code, info.Name, info.Past1Month,
			info.Past3Months, info.Past6Months, info.Past1Year, info.SinceInception))
	}

	for _, note := range service.MyConfig.Fund {
		common.Send(msg.String(), service.GetWechatUrl(note.Notes))
	}
}
