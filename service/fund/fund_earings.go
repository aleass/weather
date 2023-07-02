package fund

import (
	"bytes"
	"fmt"
	"net/http"
	"time"
	"weather/common"
	"weather/model"
	"weather/service"
)

const (
	earningsUrl = "http://fund.eastmoney.com/data/rankhandler.aspx?op=ph&dt=kf&ft=zq&rs=&gs=0&sc=1nzf&st=desc&qdii=%s|&tabSubtype=,,,,,&pi=1&pn=30000&dx=1&v=0.6135069950706549"
)

var (
	earningsFormat = []byte("var rankData = {datas:[")
)

type fundEarnings struct {
	data []byte
}

func (f *fundEarnings) GetData() {
	var arrFund = [...]string{"041", "042"}
	for _, _type := range arrFund {
		url := fmt.Sprintf(earningsUrl, _type)
		f.getUrlData(url)
	}

}

func (f *fundEarnings) getUrlData(url string) {
	refer := [][2]string{
		{"Referer", "http://fund.eastmoney.com/data/fundranking.html"},
	}
	res, err := common.HttpRequest(http.MethodPost, url, nil, refer)
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
			f.data = res[:index2-1]
			f.extract()
		}
	}
}

func (f *fundEarnings) extract() {
	var bufferEarnings []model.DfFundEarnings
	var updateEarnings []model.DfFundEarnings
	service.FuncDb.Model(&model.DfFundEarnings{}).Find(&bufferEarnings)
	var earningsMap = make(map[string]int64, len(bufferEarnings))
	for _, v := range bufferEarnings {
		earningsMap[v.Code] = v.Id
	}

	bufferEarnings = bufferEarnings[:0]

	earList := bytes.Split(f.data, []byte(`","`))
	now := time.Now()

	for _, v := range earList {
		val := bytes.Split(v, []byte(","))
		var earnings = model.DfFundEarnings{
			Name:            string(val[1]),
			Code:            string(val[0]),
			Date:            now,
			DailyGrowthRate: defaultVal(string(val[6])),
			CumulativeNav:   defaultVal(string(val[5])),
			NavPerUnit:      defaultVal(string(val[4])),
			Past1Month:      defaultVal(string(val[8])),
			Past1Week:       defaultVal(string(val[7])),
			Past1Year:       defaultVal(string(val[11])),
			Past2Years:      defaultVal(string(val[12])),
			Past3Months:     defaultVal(string(val[9])),
			Past3Years:      defaultVal(string(val[13])),
			Past6Months:     defaultVal(string(val[10])),
			SinceInception:  defaultVal(string(val[15])),
			ThisYear:        defaultVal(string(val[14])),
		}
		if id, ok := earningsMap[earnings.Code]; ok {
			earnings.Id = id
			updateEarnings = append(updateEarnings, earnings)
			continue
		}
		bufferEarnings = append(bufferEarnings, earnings)
	}
	if len(bufferEarnings) > 0 {
		service.FuncDb.Create(bufferEarnings)
	}
	if len(updateEarnings) > 0 {
		service.FuncDb.Updates(updateEarnings)
	}
	f.data = f.data[:0]
}
