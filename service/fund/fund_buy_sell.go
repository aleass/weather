package fund

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
	"weather/common"
	"weather/model"
	"weather/service"
)

var (
	byteDatasFormatStart = []byte(`datas:`)
	byteDatasFormatEnd   = []byte(`,count`)
)

type FundBuySell struct {
}

// 获取数据
func (f *FundBuySell) GetData() {
	var temp [][]string
	res, err := http.Get(common.FundBuySellUrl)
	if err != nil {
		common.Logger.Error(err.Error())
		return
	}
	raw, err := io.ReadAll(res.Body)
	if err != nil {
		common.Logger.Error(err.Error())
		return
	}
	if len(raw) == 0 {
		return
	}
	if index := bytes.Index(raw, byteDatasFormatStart); index != -1 {
		raw = raw[index+len(byteDatasFormatStart):]
		if index = bytes.Index(raw, byteDatasFormatEnd); index != -1 {
			json.Unmarshal(raw[:index], &temp)
			f.extract(temp)
		}
	}

}

// 数据抽取
func (f *FundBuySell) extract(dataSlice [][]string) {
	//检查新增的基金
	var df []model.DfFundList
	var models = &model.DfFundList{}
	service.FuncDb.Model(models).Select("id,code").Find(&df)
	var codeMap = make(map[string]int64, len(df))
	for _, fund := range df {
		codeMap[fund.Code] = fund.Id
	}
	var types = map[string]struct{}{}
	updateBuff := strings.Builder{}
	sql := "UPDATE `df_fund_list` SET `buy`='%s',`sell`='%s',`date`='%s' WHERE `id` = %d;"
	var now = time.Now().Format("2006-01-02 15:04:05")
	for _, trend := range dataSlice {
		id, ok := codeMap[trend[0]]
		if !ok {
			continue
		}
		updateBuff.WriteString(fmt.Sprintf(sql, trend[9], trend[10], now, id))
		types[trend[9]] = struct{}{}
		types[trend[10]] = struct{}{}
	}
	if updateBuff.Len() != 0 {
		service.FuncDb.Exec(updateBuff.String())
	}
	for k := range types {
		println(k)
	}
}

/*
buy
限大额
暂停申购
开放申购

sell
封闭期
开放赎回
暂停赎回
*/
