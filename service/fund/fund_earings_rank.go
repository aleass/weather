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

/*
排名和日收益
*/
const (
	async = 5
)

var (
	dataRateInSimilarType    = []byte(`Data_rateInSimilarType = `)  //同类排名走势
	dataRateInSimilarPersent = []byte(`Data_rateInSimilarPersent=`) //同类排名百分比
	unitNV                   = []byte(`Data_netWorthTrend = `)      //同类排名百分比
	totalNV                  = []byte(`Data_ACWorthTrend = `)       //累计净值
	//收益率
	syl1n = []byte(`syl_1n="`) //近一年收益率
	syl6y = []byte(`syl_6y="`) //近6月收益率
	syl3y = []byte(`syl_3y="`) //近三月收益率
	syl1y = []byte(`syl_1y="`) //近一月收益率
)

type rank struct {
	X  int64  `json:"x"`
	Y  int    `json:"y"`
	Sc string `json:"sc"`
}
type rankPrecent [][]float64
type FundEaringsRank struct {
	run string
}

type unitNVData struct {
	X int64   `json:"x"`
	Y float64 `json:"y"`
	//EquityReturn float64 `json:"equityReturn"` //净值回报
	//UnitMoney    string  `json:"unitMoney"` // 每份派送金
}

type equ struct {
	date string
	code string
}

func (f *FundEaringsRank) GetData() {
	var now = time.Now().Format("20060102")
	if now == f.run {
		return
	}
	var earningsRankMode []model.DfFundEarningsRank
	months := time.Now().AddDate(0, -2, 0).Format(common.UsualTime)
	service.FuncDb.Select("date,code").Model(&model.DfFundEarningsRank{}).Where("date >= ?", months).Find(&earningsRankMode)
	earningsRankMap := make(map[equ]struct{}, len(earningsRankMode))
	for _, data := range earningsRankMode {
		earningsRankMap[equ{data.Date, data.Code}] = struct{}{}
	}

	//收益率
	var dfEarningsMode []model.DfFundEarnings
	earningsMap := make(map[string]int64, len(earningsRankMode))
	service.FuncDb.Select("id,code").Model(&model.DfFundEarnings{}).Find(&dfEarningsMode)
	for _, data := range dfEarningsMode {
		earningsMap[data.Code] = data.Id
	}

	codes := make(chan [2]string, 1000)

	var df []model.DfFundList
	service.FuncDb.Model(&model.DfFundList{}).Where("`type` LIKE '债券型%'").Find(&df)

	var closeChan = make(chan struct{}, async)
	for i := 0; i < async; i++ {
		go f.getEaringsRankUrlData(codes, closeChan, &earningsRankMap, &earningsMap)
	}

	for _, fund := range df {
		codes <- [2]string{fund.Code, fund.Name}
	}
	close(codes)
	for i := 0; i < async; i++ {
		<-closeChan
	}

	close(closeChan)
	f.run = now
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

func (f *FundEaringsRank) getEaringsRankUrlData(codes chan [2]string, closeChan chan struct{}, earingsRankMap *map[equ]struct{}, earningsMap *map[string]int64) {
	refer := [][2]string{
		{"Referer", "http://fundf10.eastmoney.com/"},
		{"Host", "api.fund.eastmoney.com"},
	}

	//默认值
	defVal := model.DfFundEarningsRank{
		Rank:        0,
		RankPrecent: 0,
		TotalRate:   "0",
		KindAvgRate: "0",
	}

	var modelMap = make(map[string]*model.DfFundEarningsRank, 500)
	var list = make([]*model.DfFundEarningsRank, 0, 500)
	now := time.Now()
	for data := range codes {
		code := data[0]
		name := data[1]
		defVal.CreateTime = time.Now()

		//不做并发
		raw := f.GetUrlData(common.PostType, fmt.Sprintf(common.TotalEarningsUrl, code), refer)
		earnings := &totalEarnings{}
		json.Unmarshal(raw, earnings)
		if len(earnings.Data) < 2 {
			continue
		}
		fundData := earnings.Data[0]
		unix := fundData.Data[len(fundData.Data)-1][0]
		date := time.Unix(int64(unix)/1000, 0)
		if _, ok := (*earingsRankMap)[equ{date: date.Format(common.UsualTimeInt), code: code}]; ok {
			continue
		}
		//累计
		for _, datum := range earnings.Data[0].Data {
			_date := time.Unix(int64(datum[0])/1000, 0).Format(common.UsualTimeInt)
			_model := defVal
			_model.Name = name
			_model.Code = code
			_model.Date = _date
			_model.TotalRate = fmt.Sprintf("%f", datum[1])
			list = append(list, &_model)
			modelMap[_date] = &_model
		}

		//平均
		for _, datum := range earnings.Data[1].Data {
			_date := time.Unix(int64(datum[0])/1000, 0).Format(common.UsualTimeInt)
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

		//详情数据
		raw = f.GetUrlData(http.MethodGet, fmt.Sprintf(common.RankUrl, code), refer)
		sql := "UPDATE `df_fund_earnings` SET date = '%s',`past_1_month`=%s,`past_1_year`=%s,`past_3_months`=%s,`past_6_months`=%s WHERE `id` =%d"
		var past_1_month, past_1_year, past_3_months, past_6_months string

		id, ok := (*earningsMap)[code]
		if !ok {
			sql = "INSERT INTO `fund`.`df_fund_earnings`(`code`, `name`, `date`, `past_1_month`, `past_1_year`, " +
				"`past_3_months`, `past_6_months`) VALUES ('%s','%s',%s,%s,%s,%s,%s)"
		}

		//收益率
		var hasEarnings bool
		//近一年
		var pastTemp = f.extract2(raw, syl1n, []byte{'"'})
		if pastTemp != nil {
			hasEarnings = true
			past_1_year = string(pastTemp)
		}

		//近6月
		pastTemp = f.extract2(raw, syl6y, []byte{'"'})
		if pastTemp != nil {
			hasEarnings = true
			past_6_months = string(pastTemp)
		}

		//近3月
		pastTemp = f.extract2(raw, syl3y, []byte{'"'})
		if pastTemp != nil {
			hasEarnings = true
			past_3_months = string(pastTemp)
		}

		//近一月
		pastTemp = f.extract2(raw, syl1y, []byte{'"'})
		if pastTemp != nil {
			hasEarnings = true
			past_1_month = string(pastTemp)
		}

		if hasEarnings {
			if ok {
				sql = fmt.Sprintf(sql, now.Format("2006-01-02 15:04:05"), past_1_month, past_1_year, past_3_months, past_6_months, id)
			} else {
				sql = fmt.Sprintf(sql, code, name, now.Format("2006-01-02 15:04:05"), past_1_month, past_1_year, past_3_months, past_6_months)
			}
			service.FuncDb.Exec(sql)
		}

		//排名
		var grandTotal = f.extract2(raw, dataRateInSimilarType, []byte{';'})
		var tempData = []rank{}
		if grandTotal == nil {
			goto avg
		}
		json.Unmarshal(grandTotal, &tempData)

		for _, str := range tempData {
			_date := time.Unix(str.X/1000, 0).Format(common.UsualTimeInt)
			if _val, ok := modelMap[_date]; !ok {
				_model := defVal
				_model.Name = name
				_model.Code = code
				_model.Date = _date
				_model.Rank = str.Y
				list = append(list, &_model)
				modelMap[_date] = &_model
			} else {
				_val.Rank = str.Y
			}
		}

	avg:
		//排名百分比
		grandTotal = f.extract2(raw, dataRateInSimilarPersent, []byte{';'})
		var tempData2 = rankPrecent{}
		if grandTotal == nil {
			goto unit
		}
		json.Unmarshal(grandTotal, &tempData2)
		for _, _rankPre := range tempData2 {
			_date := time.Unix(int64(_rankPre[0]/1000), 0).Format(common.UsualTimeInt)
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
	unit:
		//单位净值
		grandTotal = f.extract2(raw, unitNV, []byte{';'})
		var unitList = []unitNVData{}
		if grandTotal == nil {
			goto total
		}
		json.Unmarshal(grandTotal, &unitList)
		for _, _rankPre := range unitList {
			_date := time.Unix(_rankPre.X/1000, 0).Format(common.UsualTimeInt)
			if val, ok := modelMap[_date]; !ok {
				_model := defVal
				_model.Name = name
				_model.Code = code
				_model.Date = _date
				_model.UnitNV = _rankPre.Y
				list = append(list, &_model)
				modelMap[_date] = &_model
			} else {
				val.UnitNV = _rankPre.Y
			}
		}

	total:
		//单位净值
		grandTotal = f.extract2(raw, totalNV, []byte{';'})
		var totalList = rankPrecent{}
		var last float64 = 0
		if grandTotal == nil {
			goto end
		}
		json.Unmarshal(grandTotal, &totalList)
		for _, _rankPre := range totalList {
			_date := time.Unix(int64(_rankPre[0]/1000), 0).Format(common.UsualTimeInt)
			val, ok := modelMap[_date]
			if !ok {
				_model := defVal
				_model.Name = name
				_model.Code = code
				_model.Date = _date
				list = append(list, &_model)
				modelMap[_date] = &_model
				val = &_model
			}
			val.TotalNV = _rankPre[1]
			val.DayIncreVal = _rankPre[1] - last
			val.DayIncreRate = val.DayIncreVal * 100
			last = val.TotalNV
		}

	end:
		if len(list) > 0 {
			service.FuncDb.CreateInBatches(list, 50)
			modelMap = make(map[string]*model.DfFundEarningsRank, 500)
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
