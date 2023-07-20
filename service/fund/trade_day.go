package fund

import (
	"encoding/json"
	"fmt"
	"net/http"
	"weather/common"
	"weather/model"
	"weather/service"
)

type TradeDayResp struct {
	Data struct {
		Klines []string `json:"klines"`
	} `json:"data"`
}

func TradeDay() {
	common.Logger.Info("执行 获取交易日")
	var db model.TradeDay
	err := service.FuncDb.Model(&model.TradeDay{}).Order("date desc").Limit(1).Find(&db).Error
	if err != nil {
		common.Logger.Error(err.Error())
		return
	}

	var beg int64 = 0
	if db.Date != 0 {
		beg = db.Date
	}
	url := fmt.Sprintf(common.TradeDayUrl, beg)
	raw, err := common.HttpRequest(http.MethodGet, url, nil, nil)
	if err != nil {
		common.Logger.Error(err.Error())
		return
	}
	var resp TradeDayResp
	json.Unmarshal(raw, &resp)

	_date := make([]byte, 8)
	var tradeList = make([]model.TradeDay, 0, len(resp.Data.Klines))
	for _, date := range resp.Data.Klines {
		i := copy(_date, date[:4])
		i += copy(_date[i:], date[5:7])
		copy(_date[i:], date[8:])
		var dateInt = common.Str2Int64(string(_date))
		if db.Date == dateInt {
			continue
		}
		tradeList = append(tradeList, model.TradeDay{dateInt})
	}
	err = service.FuncDb.CreateInBatches(tradeList, 10000).Error
	if err != nil {
		common.Logger.Error(err.Error())
		return
	}

	common.Logger.Info("结束 获取交易日")
}
