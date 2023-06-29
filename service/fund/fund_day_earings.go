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
	dayEaringsUrl = "http://fund.eastmoney.com/Data/Fund_JJJZ_Data.aspx?t=1&lx=13&feature=%s&gsid=&text=&sort=zdf,desc&page=1,20000&dt=%d&atfc=&onlySale=0"
)

var (
	byteDatasFormatStart = []byte(`datas:`)
	byteDatasFormatEnd   = []byte(`,count:`)
	byteshowdayFormat    = []byte(`showday:`)
)

type fundTeyp struct {
	code string
	name string
}

// 日收益
func fundDayEarings() {
	now := time.Now()
	millisecond := now.UnixMilli()
	var arrFund = [...]fundTeyp{{"041", "长期纯债"}, {"042", "短期纯债"}}
	for _, data := range arrFund {
		url := fmt.Sprintf(dayEaringsUrl, data.code, millisecond)
		getDfDayFund(url, data.name)
	}
}

// 获取数据
func getDfDayFund(url, name string) {
	var showDay []byte
	var temp [][]string
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
		raw = raw[index+len(byteDatasFormatStart):]
		if index = bytes.Index(raw, byteshowdayFormat); index != -1 {
			showDay = raw[index+len(byteshowdayFormat) : len(raw)-1]
		}
		if index = bytes.Index(raw, byteDatasFormatEnd); index != -1 {
			json.Unmarshal(raw[:index], &temp)
			extract(temp, showDay, name)
		}
	}

}

// 数据抽取
func extract(dataSlice [][]string, showDay []byte, name string) {
	if len(showDay) < 12 {
		return
	}
	showDay = showDay[2:12]
	copy(showDay[4:], showDay[5:])
	copy(showDay[6:], showDay[7:])
	date := common.Str2Int64(string(showDay[:len(showDay)-2]))
	//检查日期
	var dfDateModel = &model.DfFundDayEarings{}
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

	//检查新增的基金
	var df []model.DfFundList
	service.FuncDb.Model(&model.DfFundList{}).Find(&df)
	var codeMap = make(map[string]struct{}, len(df))
	for _, fund := range df {
		codeMap[fund.Code] = struct{}{}
	}
	//复用
	df = df[:0]

	var (
		msgBuffer           = strings.Builder{}
		msgBufferMinus      = strings.Builder{}
		i, j                int
		bufferDfFundEarings = make([]model.DfFundDayEarings, 0, len(dataSlice))
	)

	for _, trend := range dataSlice {
		if len(trend) < 18 || trend[9] != "开放申购" {
			continue
		}
		increVal, _ := strconv.ParseFloat(defaultVal(trend[7]), 64)
		increRate, _ := strconv.ParseFloat(defaultVal(trend[8]), 64)
		models := model.DfFundDayEarings{
			Date:          date,
			Code:          trend[0],
			Name:          trend[1],
			UnitNV:        defaultVal(trend[3]),
			TotalNV:       defaultVal(trend[4]),
			DayIncreVal:   increVal,
			DayIncreRate:  increRate,
			BuyStatus:     trend[9],
			SellStatus:    trend[10],
			ServiceCharge: "0",
			Type:          name,
		}
		if len(trend[18]) > 0 {
			models.ServiceCharge = trend[18][:len(trend[18])-1]
		}
		bufferDfFundEarings = append(bufferDfFundEarings, models)

		if _, ok := codeMap[models.Code]; !ok {
			df = append(df, model.DfFundList{
				Code:   models.Code,
				Name:   models.Name,
				Date:   date,
				Pinyin: trend[2],
			})
		}
		if increRate > 0 {
			i++
		}
		if increRate < 0 {
			j++
		}
	}

	if len(bufferDfFundEarings) == 0 {
		return
	}

	sort.Slice(bufferDfFundEarings, func(i, j int) bool {
		return bufferDfFundEarings[i].DayIncreRate > bufferDfFundEarings[j].DayIncreRate
	})

	for _, models := range bufferDfFundEarings[:5] {
		msgBuffer.WriteString(fmt.Sprintf("%s %s  涨率:%0.2f 涨值:%0.4f\r\n", models.Code, models.Name, models.DayIncreRate, models.DayIncreVal))
	}

	for _c := 0; _c < 5; _c++ {
		var models = bufferDfFundEarings[len(bufferDfFundEarings)-_c-1]
		msgBuffer.WriteString(fmt.Sprintf("%s %s  涨率:%0.2f 涨值:%0.4f\r\n", models.Code, models.Name, models.DayIncreRate, models.DayIncreVal))
	}

	db = service.FuncDb.CreateInBatches(&bufferDfFundEarings, 1000)
	if err = db.Error; err != nil {
		common.Logger.Error(err.Error())
		return
	}

	//新增基金
	if len(df) > 0 {
		db = service.FuncDb.Create(df)
		if err = db.Error; err != nil {
			common.Logger.Error(err.Error())
			return
		}
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

}

func defaultVal(val string) string {
	if val == "" {
		return "0"
	}
	return val
}
