package typhoon

import (
	"fmt"
	"services/api/atmp"
	"services/api/telegram"
	"services/common"
	"strings"
	"time"
)

var (
	typhoonActiveUrl = "https://typhoon.slt.zj.gov.cn/Api/TyhoonActivity"
	title            = "【台风】\n"
	temp             = common.SubStr + "%s级%s:%s  风速:%sm/s  距离:%0.2fkm %s\n"
)

var (
	lastUpdate = map[string][2]string{}
)

// 有台风1小时 否则1天
func TyphoonActive(loc string) string {
	var typhoonActiveResp []TyphoonActiveResp
	_, err := common.HttpRequest(common.OtherType, common.GetType, typhoonActiveUrl, nil, common.Header, false, &typhoonActiveResp)
	if err != nil {
		common.Logger.Error(err.Error())
		return ""
	}
	//置空
	if len(typhoonActiveResp) == 0 {
		lastUpdate = map[string][2]string{}
		return ""
	}

	var lonSelf, latSelf = common.LocStr2float(loc)
	var message string
	for _, v := range typhoonActiveResp {
		//计算距离，大于1000公里跳过
		dis := common.Haversine(common.Str2Float64(v.Lng), common.Str2Float64(v.Lat), lonSelf, latSelf)
		//if dis > 1000 {
		//	continue
		//}

		//如果数据没变化不更新图片
		if v.Timeformate == lastUpdate[v.Name][0] { //11月9日18时
			message += lastUpdate[v.Name][1]
			continue
		}

		//获取数据
		locArr, forecasts, forecastsName, forecastDate, lastes := TyphoonPath(v.Tfid)

		//贫血信息
		msg := fmt.Sprintf(temp, v.Power, v.Strong, v.Name, v.Speed, dis, lastes.Time[11:16])
		message += msg
		lastUpdate[v.Name] = [2]string{v.Timeformate, msg}

		//获取级别
		radius7 := radiusHanlder(lastes.Radius7)
		radius12 := radiusHanlder(lastes.Radius12)
		radius10 := radiusHanlder(lastes.Radius10)
		resp := atmp.CreatePhoto(loc, locArr, forecasts, forecastsName, forecastDate, dis, radius7, radius10, radius12)
		if resp == nil {
			continue
		}
		telegram.SendPhoto(resp, v.Name, "")
	}
	return title + message + "\n"
}

func radiusHanlder(data string) float64 {
	if data == "" {
		return 0
	}
	parts := strings.Split(data, "|")
	var max float64
	for _, part := range parts {
		var dis = common.Str2Float64(part)
		if dis > max {
			max = dis
		}
	}
	return max
}

type TyphoonActiveResp struct {
	Enname        string    `json:"enname"`
	Lat           string    `json:"lat"`
	Lng           string    `json:"lng"`
	Movedirection string    `json:"movedirection"`
	Movespeed     string    `json:"movespeed"`
	Name          string    `json:"name"`
	Power         string    `json:"power"`
	Pressure      string    `json:"pressure"`
	Radius10      *string   `json:"radius10"`
	Radius7       string    `json:"radius7"`
	Speed         string    `json:"speed"`
	Strong        string    `json:"strong"`
	Tfid          string    `json:"tfid"`
	Time          time.Time `json:"time"`
	Timeformate   string    `json:"timeformate"`
}
