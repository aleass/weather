package fund

import (
	"bytes"
	"io"
	"net/http"
	"time"
	"weather/common"
	"weather/model"
	"weather/service"
)

const ()

var (
	starDatasFormat = []byte(`var fundinfos = "`)
	zq              = "债券型"
	JG2Pjrq         = []byte("JG_2_pjrq = \"") //招商证券
	JG3Pjrq         = []byte("JG_3_pjrq = \"") //上海证券
	JG5Pjrq         = []byte("JG_5_pjrq = \"") //济安金信
)

type fundStar struct {
	data []byte
}

func (f *fundStar) GetData() {
	res, err := http.Get(common.StarUrl)
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
	index := bytes.Index(raw, starDatasFormat)
	if index == -1 {
		return
	}
	var i1, i2 int
	i1 = index + len(starDatasFormat)
	if index = bytes.IndexByte(raw[i1:], ';'); index == -1 {
		return
	}
	i2 = index + i1

	var (
		ZhaoShangDate = ""
		ShanghaiDate  = ""
		JiananDate    = ""
	)

	if index = bytes.Index(raw[i2:], JG2Pjrq); index != -1 {
		index += i2 + len(JG2Pjrq)
		if index2 := bytes.IndexByte(raw[index:], '"'); index2 != -1 {
			ZhaoShangDate = string((raw[index : index+index2]))
		}
	}

	if index = bytes.Index(raw[i2:], JG3Pjrq); index != -1 {
		index += i2 + len(JG3Pjrq)
		if index2 := bytes.IndexByte(raw[index:], '"'); index2 != -1 {
			ShanghaiDate = string(raw[index : index+index2])
		}
	}
	if index = bytes.Index(raw[i2:], JG5Pjrq); index != -1 {
		index += i2 + len(JG5Pjrq)
		if index2 := bytes.IndexByte(raw[index:], '"'); index2 != -1 {
			JiananDate = string((raw[index : index+index2]))
		}
	}

	//数据比较
	var dbData []model.DfFundStar
	var codeMap map[string]int64
	service.FuncDb.Model(&model.DfFundStar{}).Find(&dbData)
	if len(dbData) != 0 {
		//时间全一样,返回
		if dbData[0].JiananJinxinSecuritiesDate == JiananDate && dbData[0].ShanghaiSecuritiesDate == ShanghaiDate &&
			dbData[0].ZhaoShangSecuritiesDate == ZhaoShangDate {
			//return
		}
		codeMap = make(map[string]int64, len(dbData))
		for _, fund := range dbData {
			codeMap[fund.Code] = fund.Id
		}
	}

	f.startExtract(raw[i1:i2], ZhaoShangDate, ShanghaiDate, JiananDate, codeMap)
}

func (f *fundStar) startExtract(data []byte, ZhaoShangDate, ShanghaiDate, JiananDate string, fundCodeMap map[string]int64) {
	fundList := bytes.Split(data, []byte("_"))
	bufferDfFundEarnings := make([]model.DfFundStar, 0, 100)
	updateDfFundEarnings := make([]model.DfFundStar, 0, 100)
	now := time.Now()
	for _, fund := range fundList {
		info := bytes.Split(fund, []byte("|"))
		if len(info) < 25 {
			continue
		}
		//过滤
		if len(info[2]) <= len(zq) || string(info[2][:len(zq)]) != zq {
			continue
		}

		//buff full
		if len(bufferDfFundEarnings) > 100 {
			db := service.FuncDb.Create(&bufferDfFundEarnings)
			if err := db.Error; err != nil {
				common.Logger.Error(err.Error())
				return
			}
			bufferDfFundEarnings = bufferDfFundEarnings[:0]
		}
		models := model.DfFundStar{
			Code:                       string(info[0]),
			JiananJinxinSecuritiesDate: JiananDate,
			JiananJinxinTrend:          f.starHandler(info[17]),
			Name:                       string(info[1]),
			ShanghaiSecuritiesDate:     ShanghaiDate,
			ShanghaiSecuritiesTrend:    f.starHandler(info[13]),
			UpdateTime:                 now,
			ZhaoShangSecuritiesDate:    ZhaoShangDate,
			ZhaoShangSecuritiesTrend:   f.starHandler(info[11]),
			JiananJinxinStar:           common.DefaultVal(string(info[16])),
			ZhaoShangSecuritiesStar:    common.DefaultVal(string(info[10])),
			ShanghaiSecuritiesStar:     common.DefaultVal(string(info[12])),
		}

		//更新
		if id, ok := fundCodeMap[models.Code]; ok {
			models.Id = id
			updateDfFundEarnings = append(updateDfFundEarnings, models)
			continue
		}

		bufferDfFundEarnings = append(bufferDfFundEarnings, models)
	}

	if len(bufferDfFundEarnings) != 0 {
		db := service.FuncDb.Create(&bufferDfFundEarnings)
		if err := db.Error; err != nil {
			common.Logger.Error(err.Error())
			return
		}
	}

	if len(updateDfFundEarnings) > 0 {
		service.FuncDb.Save(updateDfFundEarnings)
	}
	f.data = f.data[:0]

}

func (f *fundStar) starHandler(_range []byte) string {
	if len(_range) == 0 {
		return ""
	}
	switch _range[0] {
	case '0':
		return ""
	case '-':
		return "down " + string(_range[1:])
	default:
		return "up " + string(_range)
	}
}
