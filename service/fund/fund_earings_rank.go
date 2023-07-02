package fund

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
	"weather/common"
	"weather/model"
	"weather/service"
)

const (
	totalEaringsUrl = "http://api.fund.eastmoney.com/f10/FundLJSYLZS/?bzdm=%s&dt=month"             //dt=单位，all全部
	rankUrl         = "http://api.fund.eastmoney.com/f10/FundTLPM/?bzdm=%s&dt=month&rt=threemonth"  //dt=单位，all全部
	precentRankUrl  = "http://api.fund.eastmoney.com/f10/FundBFBPM/?bzdm=%s&dt=month&rt=threemonth" //dt=单位，all全部
	async           = 5
)

type FundEaringsRank struct {
}

type equ struct {
	date string
	code string
}

func (f *FundEaringsRank) GetData() {
	var earingsRankMode []model.FundEaringsRank
	months := time.Now().AddDate(0, -2, 0).Format("20060102")
	service.FuncDb.Select("date,code").Model(&model.FundEaringsRank{}).Where("date >= ?", months).Find(&earingsRankMode)
	earingsRankMap := make(map[equ]struct{}, len(earingsRankMode))
	for _, data := range earingsRankMode {
		earingsRankMap[equ{data.Date, data.Code}] = struct{}{}
	}
	codes := make(chan [2]string, 1000)

	var df []model.DfFundList
	service.FuncDb.Model(&model.DfFundList{}).Where("`type` LIKE '债券型-%'").Find(&df)

	//
	var closeChan = make(chan struct{}, async)
	for i := 0; i < async; i++ {
		go f.getEaringsRankUrlData(codes, closeChan, &earingsRankMap)
	}

	for _, fund := range df {
		codes <- [2]string{fund.Code, fund.Name}
	}
	close(codes)
	for i := 0; i < async; i++ {
		<-closeChan
	}
}

type EaringsRankRes struct {
	Data string `json:"Data"`
}

func (f *FundEaringsRank) getEaringsRankUrlData(codes chan [2]string, closeChan chan struct{}, earingsRankMap *map[equ]struct{}) {
	refer := [][2]string{
		{"Referer", "http://fundf10.eastmoney.com/"},
		{"Host", "api.fund.eastmoney.com"},
	}

	//默认值
	defVal := model.FundEaringsRank{
		Rank:        "0",
		RankPrecent: "0",
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
		totalEaringStr := f.GetUrlData(code, fmt.Sprintf(totalEaringsUrl, code), refer, earingsRankMap)
		if len(totalEaringStr) == 0 {
			continue
		}
		for _, str := range totalEaringStr {
			info := strings.Split(str, "_")
			if len(info) < 3 {
				goto end1
			}
			date := strings.ReplaceAll(info[0], "/", "")
			_model := defVal
			_model.Name = name
			_model.Code = code
			_model.Date = date
			_model.TotalRate = defaultVal(info[2])
			_model.KindAvgRate = defaultVal(info[1])

			list = append(list, &_model)
			modelMap[date] = &_model
		}
	end1:
		rankStr := f.GetUrlData(code, fmt.Sprintf(rankUrl, code), refer, earingsRankMap)
		for _, str := range rankStr {
			info := strings.Split(str, "_")
			if len(info) < 2 {
				goto end2
			}
			date := strings.ReplaceAll(info[0], "/", "")
			if val, ok := modelMap[date]; !ok {
				_model := defVal
				_model.Name = name
				_model.Code = code
				_model.Date = date
				_model.Rank = defaultVal(info[1])

				list = append(list, &_model)
				modelMap[date] = &_model
				continue
			} else {
				val.Rank = defaultVal(info[1])
			}
		}
	end2:

		rankPrecentStr := f.GetUrlData(code, fmt.Sprintf(precentRankUrl, code), refer, earingsRankMap)
		for _, str := range rankPrecentStr {
			info := strings.Split(str, "_")
			if len(info) < 2 {
				goto end3
			}
			date := strings.ReplaceAll(info[0], "/", "")
			if val, ok := modelMap[date]; !ok {
				_model := defVal
				_model.Name = name
				_model.Code = code
				_model.Date = date
				_model.RankPrecent = defaultVal(info[1])

				list = append(list, &_model)
				modelMap[date] = &_model
				continue
			} else {
				val.RankPrecent = defaultVal(info[1])
			}
		}
	end3:
		if len(list) > 0 {
			service.FuncDb.CreateInBatches(list, 50)
			modelMap = make(map[string]*model.FundEaringsRank, 500)
			list = list[:0]
		}
	}
	closeChan <- struct{}{}
}

func (f *FundEaringsRank) GetUrlData(code, url string, refer [][2]string, earingsRankMap *map[equ]struct{}) []string {
	res, err := common.HttpRequest(http.MethodPost, url, nil, refer)
	if err != nil {
		common.Logger.Error(err.Error())
		return nil
	}
	data := &EaringsRankRes{}
	json.Unmarshal(res, data)
	return f.extract(code, data.Data, earingsRankMap)
}

type sortMap struct {
	info map[string]string
	list []string
}

func (f *FundEaringsRank) extract(code, data string, earingsRankMap *map[equ]struct{}) []string {
	dateList := strings.Split(data, "|")
	//2016/10/21_17.3500
	info := strings.Split(dateList[len(dateList)-1], "_")
	lastDate := strings.ReplaceAll(info[0], "/", "")
	if _, ok := (*earingsRankMap)[equ{date: lastDate, code: code}]; ok {
		return nil
	}
	//var maps = sortMap{
	//	make(map[string]string,len(dateList)),
	//	make([]string,len(dateList)),
	//}
	//for _, d := range dateList {
	//
	//}
	//

	return dateList
}
