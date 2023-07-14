package fund

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
	"time"
	"weather/common"
	"weather/model"
	"weather/service"
)

/*
阶段收益
*/

var (
	earningsFormat = []byte("var rankData = {datas:[")
)

type fundEarnings struct {
}

func (f *fundEarnings) GetData() {
	common.Logger.Info("执行 阶段收益任务")

	refer := [][2]string{
		{"Referer", "http://fund.eastmoney.com/data/fundranking.html"},
	}
	res, err := common.HttpRequest(http.MethodPost, common.EarningsUrl, nil, refer)
	if err != nil {
		common.Logger.Error(err.Error())
		return
	}

	if len(res) == 0 {
		return
	}

	if index := bytes.Index(res, earningsFormat); index != -1 {
		res = res[index+len(earningsFormat)+1:]
		if index2 := bytes.IndexByte(res, ']'); index2 != -1 {
			f.extract(res[:index2-1])
		}
	}
}

func (f *fundEarnings) extract(data []byte) {
	var bufferEarnings []model.DfFundEarnings
	var updateEarnings []model.DfFundEarnings
	service.FuncDb.Model(&model.DfFundEarnings{}).Find(&bufferEarnings)
	var earningsMap = make(map[string]int64, len(bufferEarnings))
	for _, v := range bufferEarnings {
		earningsMap[v.Code] = v.Id
	}

	bufferEarnings = bufferEarnings[:0]

	earList := bytes.Split(data, []byte(`","`))
	now := time.Now()
	nowDate := now.Format("2006-01-02 15:04:05")

	var df []model.DfFundList
	service.FuncDb.Model(&model.DfFundList{}).Select("id,code").Where("Inc_date is null").Find(&df)
	var codeMap = make(map[string]int64, len(df))
	for _, fund := range df {
		codeMap[fund.Code] = fund.Id
	}

	updateBuff := strings.Builder{}
	sql := "UPDATE `df_fund_list` SET `Inc_date`='%s',`date`='%s' WHERE `id` = %d;"

	for _, v := range earList {
		val := bytes.Split(v, []byte(","))
		var earnings = model.DfFundEarnings{
			Name:            string(val[1]),
			Code:            string(val[0]),
			Date:            now,
			DailyGrowthRate: common.DefaultVal(string(val[6])),
			CumulativeNav:   common.DefaultVal(string(val[5])),
			NavPerUnit:      common.DefaultVal(string(val[4])),
			Past1Month:      common.DefaultVal(string(val[8])),
			Past1Week:       common.DefaultVal(string(val[7])),
			Past1Year:       common.DefaultVal(string(val[11])),
			Past2Years:      common.DefaultVal(string(val[12])),
			Past3Months:     common.DefaultVal(string(val[9])),
			Past3Years:      common.DefaultVal(string(val[13])),
			Past6Months:     common.DefaultVal(string(val[10])),
			SinceInception:  common.DefaultVal(string(val[15])),
			ThisYear:        common.DefaultVal(string(val[14])),
		}
		//成立日
		if id, ok := codeMap[earnings.Code]; ok {
			updateBuff.WriteString(fmt.Sprintf(sql, string(val[16]), nowDate, id))
		}

		if id, ok := earningsMap[earnings.Code]; ok {
			earnings.Id = id
			updateEarnings = append(updateEarnings, earnings)
			continue
		}
		bufferEarnings = append(bufferEarnings, earnings)

	}
	if len(bufferEarnings) > 0 {
		service.FuncDb.CreateInBatches(bufferEarnings, 100)
	}
	if len(updateEarnings) > 0 {
		service.FuncDb.Updates(updateEarnings)
	}

	if updateBuff.Len() != 0 {
		service.FuncDb.Exec(updateBuff.String())
	}

}
