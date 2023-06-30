package fund

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
	"weather/common"
	"weather/model"
	"weather/service"
)

const (
	dayEarningsUrl = "http://fund.eastmoney.com/Data/Fund_JJJZ_Data.aspx?t=1&lx=13&feature=%s&gsid=&text=&sort=zdf,desc&page=1,20000&dt=%d&atfc=&onlySale=0"
)

var (
	byteDataFormatStart = []byte(`datas:`)
	byteDataFormatEnd   = []byte(`,count:`)
	byteShowdayFormat   = []byte(`showday:`)
)

type fundTeyp struct {
	code string
	name string
}

type fundDayEarnings struct {
	data [][]string
}

// 日收益
func (f *fundDayEarnings) GetData() {
	now := time.Now()
	millisecond := now.UnixMilli()
	var arrFund = [...]fundTeyp{{"041", "长期纯债"}, {"042", "短期纯债"}}
	for _, data := range arrFund {
		url := fmt.Sprintf(dayEarningsUrl, data.code, millisecond)
		f.getDfDayFund(url, data.name)
	}
}

// 获取数据
func (f *fundDayEarnings) getDfDayFund(url, name string) {
	var showDay []byte
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
	if index := bytes.Index(raw, byteDataFormatStart); index != -1 {
		raw = raw[index+len(byteDataFormatStart):]
		if index = bytes.Index(raw, byteShowdayFormat); index != -1 {
			showDay = raw[index+len(byteShowdayFormat) : len(raw)-1]
		}
		if index = bytes.Index(raw, byteDataFormatEnd); index != -1 {
			json.Unmarshal(raw[:index], &f.data)
			f.extract(showDay, name)
		}
	}

}

// 数据抽取
func (f *fundDayEarnings) extract(showDay []byte, name string) {
	if len(showDay) < 12 {
		return
	}
	showDay = showDay[2:12]
	copy(showDay[4:], showDay[5:])
	copy(showDay[6:], showDay[7:])
	date := common.Str2Int64(string(showDay[:len(showDay)-2]))
	//检查日期
	var dfDateModel = &model.DfFundDayEarnings{}
	db := service.FuncDb.Model(dfDateModel).Where("date = ? and type = ?", date, name).Find(dfDateModel).Limit(1)
	var err = db.Error
	if err != nil {
		common.Logger.Error(err.Error())
		return
	}
	//是否
	if dfDateModel.Date > 0 {
		return
	}
	var (
		msgBuffer            = strings.Builder{}
		msgBufferMinus       = strings.Builder{}
		i, j                 int
		bufferDfFundEarnings = make([]model.DfFundDayEarnings, 0, len(f.data))
	)

	for _, trend := range f.data {
		if len(trend) < 18 || trend[9] != "开放申购" {
			continue
		}
		//increVal, _ := strconv.ParseFloat(defaultVal(trend[7]), 64)
		increRate, _ := strconv.ParseFloat(defaultVal(trend[8]), 64)
		models := model.DfFundDayEarnings{
			Date:          date,
			Code:          trend[0],
			Name:          trend[1],
			UnitNV:        defaultVal(trend[3]),
			TotalNV:       defaultVal(trend[4]),
			DayIncreVal:   defaultVal(trend[7]),
			DayIncreRate:  increRate,
			BuyStatus:     trend[9],
			SellStatus:    trend[10],
			ServiceCharge: "0",
			Type:          name,
		}
		if len(trend[18]) > 0 {
			models.ServiceCharge = trend[18][:len(trend[18])-1]
		}
		bufferDfFundEarnings = append(bufferDfFundEarnings, models)

		if increRate > 0 {
			i++
		}
		if increRate < 0 {
			j++
		}
	}

	if len(bufferDfFundEarnings) == 0 {
		return
	}

	sort.Slice(bufferDfFundEarnings, func(i, j int) bool {
		return bufferDfFundEarnings[i].DayIncreRate > bufferDfFundEarnings[j].DayIncreRate
	})

	for _, models := range bufferDfFundEarnings[:5] {
		msgBuffer.WriteString(fmt.Sprintf("%s %s  涨率:%0.2f 涨值:%0.4f\r\n", models.Code, models.Name, models.DayIncreRate, models.DayIncreVal))
	}

	for _c := 0; _c < 5; _c++ {
		var models = bufferDfFundEarnings[len(bufferDfFundEarnings)-_c-1]
		msgBuffer.WriteString(fmt.Sprintf("%s %s  涨率:%0.2f 涨值:%0.4f\r\n", models.Code, models.Name, models.DayIncreRate, models.DayIncreVal))
	}

	db = service.FuncDb.CreateInBatches(&bufferDfFundEarnings, 1000)
	if err = db.Error; err != nil {
		common.Logger.Error(err.Error())
		return
	}

	if str := msgBuffer.String(); str != "" {
		for _, note := range service.MyConfig.Fund {
			common.Send(fmt.Sprintf("%s 正日增率: %d个\r\n", name, i)+str, service.GetWechatUrl(note.Notes))
		}
	}

	if str := msgBufferMinus.String(); str != "" {
		for _, note := range service.MyConfig.Fund {
			common.Send(fmt.Sprintf("%s 负日增率: %d个\n", name, j)+str, service.GetWechatUrl(note.Notes))
		}
	}
	f.data = f.data[:0]
}

func defaultVal(val string) string {
	if val == "" {
		return "0"
	}
	return val
}
