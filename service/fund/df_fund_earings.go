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
	earingsUrl = "http://fund.eastmoney.com/data/rankhandler.aspx?op=ph&dt=kf&ft=zq&rs=&gs=0&sc=1nzf&st=desc&qdii=%s|&tabSubtype=,,,,,&pi=1&pn=30000&dx=1&v=0.6135069950706549"
)

var (
	earingsForat = []byte("var rankData = {datas:[")
)

func GetEaringsReq() {
	var arrFund = [...]string{"041", "042"}
	for _, _type := range arrFund {
		url := fmt.Sprintf(earingsUrl, _type)
		GetEaringsData(url)
	}

}

func GetEaringsData(url string) {
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

	if index := bytes.Index(res, earingsForat); index != -1 {
		res = res[index+len(earingsForat)+1:]
		if index2 := bytes.IndexByte(res, ']'); index2 != -1 {
			earingsExtract(res[:index2-1])
		}
	}
}

func earingsExtract(raw []byte) {
	var bufferEarings []model.DfFundEarings
	var updateEarings []model.DfFundEarings
	service.FuncDb.Model(&model.DfFundEarings{}).Find(&bufferEarings)
	var earingsMap = make(map[string]int64, len(bufferEarings))
	for _, v := range bufferEarings {
		earingsMap[v.Code] = v.Id
	}

	bufferEarings = bufferEarings[:0]

	earList := bytes.Split(raw, []byte(`","`))
	now := time.Now()

	for _, v := range earList {
		val := bytes.Split(v, []byte(","))
		var earings = model.DfFundEarings{
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
		if id, ok := earingsMap[earings.Code]; ok {
			earings.Id = id
			updateEarings = append(updateEarings, earings)
			continue
		}
		bufferEarings = append(bufferEarings, earings)
	}
	if len(bufferEarings) > 0 {
		service.FuncDb.Create(bufferEarings)
	}
	if len(updateEarings) > 0 {
		service.FuncDb.Updates(updateEarings)
	}

}
