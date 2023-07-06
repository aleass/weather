package fund

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
	"unsafe"
	"weather/common"
	"weather/model"
	"weather/service"
)

var fundListFormat = []byte("var r = ")

type fundList struct {
	data [][]string
}

func (f *fundList) GetData() {
	raw, err := common.HttpRequest(http.MethodGet, common.FundListUrl, nil, nil)
	if err != nil {

	}
	index := bytes.Index(raw, fundListFormat)
	if index == -1 {
		return
	}
	raw = raw[index+len(fundListFormat) : len(raw)-1]
	json.Unmarshal(raw, &f.data)
	f.extract()
}

type memequal struct {
	Code       string `gorm:"column:code"        desc:"代码"`
	Name       string `gorm:"column:name"        desc:"名字"`
	Pinyin     string `gorm:"column:pinyin"      desc:"拼音"`
	AbbrPinyin string `gorm:"column:abbr_pinyin" desc:"拼音简写"`
	Type       string `gorm:"column:type"        desc:"基金类型"`
}

func (f *fundList) extract() {
	//检查新增的基金
	var df []model.DfFundList
	service.FuncDb.Model(&model.DfFundList{}).Find(&df)
	var codeMap = make(map[string]*model.DfFundList, len(df))
	for i, fund := range df {
		codeMap[fund.Code] = &df[i]
	}
	var newFund = make([]model.DfFundList, 0, 100)
	var updateFund = make([]model.DfFundList, 0, 100)
	now := time.Now()
	for _, fund := range f.data {
		temp := model.DfFundList{
			AbbrPinyin: fund[1],
			Code:       fund[0],
			Name:       fund[2],
			Pinyin:     fund[4],
			Type:       fund[3],
			Date:       now,
		}
		if _fund, ok := codeMap[fund[0]]; ok {
			e1 := *(*memequal)(unsafe.Pointer(&temp))
			e2 := *(*memequal)(unsafe.Pointer(_fund))
			if e2 == e1 {
				continue
			}
			temp.Id = _fund.Id
			updateFund = append(updateFund, temp)
			continue
		}

		newFund = append(newFund, temp)
	}

	//新增基金
	if len(newFund) > 0 {
		db := service.FuncDb.CreateInBatches(newFund, 100)
		if err := db.Error; err != nil {
			common.Logger.Error(err.Error())
			return
		}
	}
	//更新基金
	if len(updateFund) > 0 {
		db := service.FuncDb.Save(updateFund)
		if err := db.Error; err != nil {
			common.Logger.Error(err.Error())
			return
		}
	}
	//清空
	f.data = f.data[:0]
}
