package fund

import (
	"fmt"
	"go.uber.org/zap/buffer"
	"weather/common"
	"weather/service"
)

type daysPastTimeRank struct {
}

func (s *daysPastTimeRank) Send() {
	common.Logger.Info("执行 基金购买情况")
	var list []common.DaysPastTimeRank
	if db := service.FuncDb.Raw(common.DaysPastTimeRankSql).Find(&list); db.Error != nil {
		common.Logger.Error(db.Error.Error())
	}
	var msg = buffer.Buffer{}

	for _, info := range list {
		msg.WriteString(fmt.Sprintf("%s %s %s\n", info.Past1Month[:len(info.Past1Month)-2], info.Code, info.Name))
	}
	msg.WriteString(fmt.Sprintf("\n\n查看数据点击 http://%s:8080/fund/day", service.MyConfig.Fund.Host))

	for _, note := range service.MyConfig.Fund.Notes {
		common.Send(msg.String(), service.GetWechatUrl(note))
	}
}
