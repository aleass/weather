package fund

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"weather/common"
	"weather/model"
	"weather/service"
)

const (
	totalEarningsUrl = "https://api.fund.eastmoney.com/pinzhong/LJSYLZS?fundCode=%s&type=se&indexcode=000300" //dt=单位，m 1月 ,se 成立来
	rankUrl          = "https://fund.eastmoney.com/pingzhongdata/%s.js"
	async            = 5
)

var (
	dataRateInSimilarType    = []byte(`var Data_rateInSimilarType = `)  /*同类排名走势*/
	dataRateInSimilarPersent = []byte(`var Data_rateInSimilarPersent=`) /*同类排名百分比*/
)

type rank struct {
	X  int64  `json:"x"`
	Y  int    `json:"y"`
	Sc string `json:"sc"`
}
type rankPrecent [][]float64
type FundEaringsRank struct {
}

type equ struct {
	date string
	code string
}

func (f *FundEaringsRank) GetData() {
	var earningsRankMode []model.FundEaringsRank
	months := time.Now().AddDate(0, -2, 0).Format("20060102")
	service.FuncDb.Select("date,code").Model(&model.FundEaringsRank{}).Where("date >= ?", months).Find(&earningsRankMode)
	earningsRankMap := make(map[equ]struct{}, len(earningsRankMode))
	for _, data := range earningsRankMode {
		earningsRankMap[equ{data.Date, data.Code}] = struct{}{}
	}
	codes := make(chan [2]string, 1000)

	var df []model.DfFundList
	service.FuncDb.Model(&model.DfFundList{}).Where("`type` LIKE '债券型-%'").Find(&df)

	//
	var closeChan = make(chan struct{}, async)
	for i := 0; i < async; i++ {
		go f.getEaringsRankUrlData(codes, closeChan, &earningsRankMap)
	}

	for _, fund := range df {
		codes <- [2]string{fund.Code, fund.Name}
	}
	//codes <- [2]string{`006011`, `中信保诚稳鸿A`}
	close(codes)
	for i := 0; i < async; i++ {
		<-closeChan
	}
	close(closeChan)
}

type EaringsRankRes struct {
	Data string `json:"Data"`
}

type totalEarnings struct {
	Data []struct {
		Data [][]float64 `json:"data"`
		Name string      `json:"name"`
	} `json:"Data"`
}

func (f *FundEaringsRank) getEaringsRankUrlData(codes chan [2]string, closeChan chan struct{}, earingsRankMap *map[equ]struct{}) {
	refer := [][2]string{
		{"Referer", "http://fundf10.eastmoney.com/"},
		{"Host", "api.fund.eastmoney.com"},
	}

	//默认值
	defVal := model.FundEaringsRank{
		Rank:        0,
		RankPrecent: 0,
		TotalRate:   "0",
		KindAvgRate: "0",
	}

	var modelMap = make(map[string]*model.FundEaringsRank, 500)
	var list = make([]*model.FundEaringsRank, 0, 500)
	for data := range codes {
		code := data[0]
		name := data[1]
		defVal.CreateTime = time.Now()

		println(code, name)
		//不做并发
		raw := f.GetUrlData(common.PostType, fmt.Sprintf(totalEarningsUrl, code), refer)
		earnings := &totalEarnings{}
		json.Unmarshal(raw, earnings)
		if len(earnings.Data) < 2 {
			return
		}
		fundData := earnings.Data[0]
		unix := fundData.Data[len(fundData.Data)-1][0]
		date := time.Unix(int64(unix)/1000, 0)
		if _, ok := (*earingsRankMap)[equ{date: date.Format("20060102"), code: code}]; ok {
			goto total
		}

		for _, datum := range earnings.Data[0].Data {
			_date := time.Unix(int64(datum[0])/1000, 0).Format("20060102")
			_model := defVal
			_model.Name = name
			_model.Code = code
			_model.Date = _date
			_model.TotalRate = fmt.Sprintf("%f", datum[1])
			list = append(list, &_model)
			modelMap[_date] = &_model
		}

		//avg
		for _, datum := range earnings.Data[1].Data {
			_date := time.Unix(int64(datum[0])/1000, 0).Format("20060102")
			if val, ok := modelMap[_date]; !ok {
				_model := defVal
				_model.Name = name
				_model.Code = code
				_model.Date = _date
				_model.TotalRate = fmt.Sprintf("%f", datum[1])
				list = append(list, &_model)
				modelMap[_date] = &_model
			} else {
				val.KindAvgRate = fmt.Sprintf("%f", datum[1])
			}
		}

	total:
		raw = f.GetUrlData(http.MethodGet, fmt.Sprintf(rankUrl, code), refer)
		var grandTotal = f.extract2(raw, dataRateInSimilarType, []byte{';'})
		var tempData = []rank{}
		if grandTotal == nil {
			goto avg
		}
		json.Unmarshal(grandTotal, &tempData)

		for _, str := range tempData {
			date := time.Unix(str.X/1000, 0).Format("20060102")
			if val, ok := modelMap[date]; !ok {
				_model := defVal
				_model.Name = name
				_model.Code = code
				_model.Date = date
				_model.Rank = str.Y
				list = append(list, &_model)
				modelMap[date] = &_model
			} else {
				val.Rank = str.Y
			}
		}

	avg:
		grandTotal = f.extract2(raw, dataRateInSimilarPersent, []byte{';'})
		var tempData2 = rankPrecent{}
		if grandTotal == nil {
			goto end
		}
		json.Unmarshal(grandTotal, &tempData2)
		for _, _rankPre := range tempData2 {
			_date := time.Unix(int64(_rankPre[0]/1000), 0).Format("20060102")
			if val, ok := modelMap[_date]; !ok {
				_model := defVal
				_model.Name = name
				_model.Code = code
				_model.Date = _date
				_model.RankPrecent = _rankPre[1]
				list = append(list, &_model)
				modelMap[_date] = &_model
			} else {
				val.RankPrecent = _rankPre[1]
			}
		}

	end:
		if len(list) > 0 {
			service.FuncDb.CreateInBatches(list, 50)
			modelMap = make(map[string]*model.FundEaringsRank, 500)
			list = list[:0]
		}
	}
	closeChan <- struct{}{}
}
func (f *FundEaringsRank) extract2(data []byte, startStr, endStr []byte) []byte {
	if index := bytes.Index(data, startStr); index != -1 {
		data = data[len(startStr)+index:]
		if index = bytes.Index(data, endStr); index != -1 {
			return data[:index]
		}
	}
	return nil
}

func (f *FundEaringsRank) GetUrlData(_type common.HttpMethod, url string, refer [][2]string) []byte {
	res, err := common.HttpRequest(_type, url, nil, refer)
	if err != nil {
		common.Logger.Error(err.Error())
		return nil
	}
	return res
}
