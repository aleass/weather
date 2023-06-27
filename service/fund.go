package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
	"weather/common"
	"weather/model"
)

const (
	dfUrl = "http://fund.eastmoney.com/Data/Fund_JJJZ_Data.aspx?t=1&lx=13&letter=&gsid=&text=&sort=zdf,desc&page=1,20000&dt=%d&atfc=&onlySale=0"
)

var (
	byteDatasFormatStart = []byte(`datas:`)
	byteDatasFormatEnd   = []byte(`,count:`)
	byteshowdayFormat    = []byte(`showday:`)
)

func FundRun() {
	time.Sleep(time.Second)
	for {
		now := time.Now()
		if week := now.Weekday() - 1; now.Hour() == 0 && week > 0 && week < 6 {
			DfFund()
		}
		time.Sleep(time.Hour)
	}
}

func DfFund() {
	now := time.Now()
	millisecond := now.UnixMilli()
	url := fmt.Sprintf(dfUrl, millisecond)
	res, err := http.Get(url)
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
		var data, showDay []byte
		raw = raw[index+len(byteDatasFormatStart):]
		if index = bytes.Index(raw, byteshowdayFormat); index != -1 {
			showDay = raw[index+len(byteshowdayFormat) : len(raw)-1]
		}
		if index = bytes.Index(raw, byteDatasFormatEnd); index != -1 {
			data = raw[:index]
		}
		extract(data, showDay)
	}
}

func extract(data, showDay []byte) {
	showDay = showDay[2:12]
	copy(showDay[4:], showDay[5:])
	copy(showDay[6:], showDay[7:])
	date, _ := strconv.Atoi(string(showDay[:len(showDay)-2]))
	//检查日期
	var dfDateModel = &model.DfDateNote{}
	db := funcDb.Model(dfDateModel).Where("date = ?", date).Find(dfDateModel)
	var err = db.Error
	if err != nil {
		common.Logger.Error(err.Error())
		return
	}

	if dfDateModel.Date > 0 {
		//return
	}

	//检查新增的基金
	var df []model.DfFund
	funcDb.Model(&model.DfFund{}).Find(&df)
	var codeMap = make(map[string]struct{}, len(df))
	for _, fund := range df {
		codeMap[fund.Code] = struct{}{}
	}
	//复用
	df = df[:0]

	dataSlice := [][]string{}
	if err = json.Unmarshal(data, &dataSlice); err != nil {
		common.Logger.Error(err.Error())
		return
	}
	var msgBuffer = strings.Builder{}
	var msgBufferMinus = strings.Builder{}
	var i, j int
	for _, trend := range dataSlice {
		if len(trend) < 18 {
			continue
		}
		models := model.DfFundEarings{
			Date:          date,
			Code:          trend[0],
			Name:          trend[1],
			UnitNV:        defaultval(trend[3]),
			TotalNV:       defaultval(trend[4]),
			DayIncreVal:   defaultval(trend[5]),
			DayIncreRate:  defaultval(trend[8]),
			BuyStatus:     trend[9],
			SellStatus:    trend[10],
			ServiceCharge: "0",
		}
		if len(trend[18]) > 0 {
			models.ServiceCharge = trend[18][:len(trend[18])-1]
		}

		//db = funcDb.Create(&models)
		//if err = db.Error; err != nil {
		//	common.Logger.Error(err.Error())
		//	return
		//}

		if _, ok := codeMap[models.Code]; !ok {
			df = append(df, model.DfFund{
				Code:   models.Code,
				Name:   models.Name,
				Date:   date,
				Pinyin: trend[2],
			})
		}
		if models.BuyStatus != "开放申购" || models.UnitNV == models.TotalNV {
			continue
		}
		f, _ := strconv.ParseFloat(models.DayIncreRate, 64)
		if f >= 0.1 {
			i++
			msgBuffer.WriteString(fmt.Sprintf("%s %s  涨率:%s 涨值:%s\r\n", models.Code, models.Name, models.DayIncreRate, models.DayIncreVal))
		}
		if f <= -0.1 {
			j++
			msgBufferMinus.WriteString(fmt.Sprintf("%s %s  涨率:%s 涨值:%s\r\n", models.Code, models.Name, models.DayIncreRate, models.DayIncreVal))
		}

	}
	//插入更新的日期
	//dfDateModel.Date = date
	//db = funcDb.Create(dfDateModel)
	//if err = db.Error; err != nil {
	//	common.Logger.Error(err.Error())
	//	return
	//}

	println(i, j)
	//新增基金
	if len(df) > 0 {
		db = funcDb.Create(df)
		if err = db.Error; err != nil {
			common.Logger.Error(err.Error())
			return
		}
	}

	if str := msgBuffer.String(); str != "" {
		for _, note := range myConfig.Fund {
			common.Send(str, wechatUrl+wechatNoteMap[note.Notes])
		}
	}

	if str := msgBufferMinus.String(); str != "" {
		for _, note := range myConfig.Fund {
			common.Send(str, wechatUrl+wechatNoteMap[note.Notes])
		}
	}

}
func defaultval(val string) string {
	if val == "" {
		return "0"
	}
	return val
}
