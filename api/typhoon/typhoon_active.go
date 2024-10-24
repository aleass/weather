package typhoon

import (
	"fmt"
	"strings"
	"time"
	"weather/api/atmp"
	"weather/api/telegram"
	"weather/common"
)

const (
	typhoonActiveUrl = "https://typhoon.slt.zj.gov.cn/Api/TyhoonActivity"
	temp             = `【台风】%s级%s-%s
     ● 风速：%sm/s，距离：%0.2fkm，风圈：%skm，%s
`
)

var (
	lastUpdate string
	lastPoint  = &PointsInfo{}
)

// 有台风1小时 否则1天
func TyphoonActive() string {
	var typhoonActiveResp []TyphoonActiveResp
	_, err := common.HttpRequest(common.OtherType, common.GetType, typhoonActiveUrl, nil, common.Header, false, &typhoonActiveResp)
	if err != nil {
		common.Logger.Error(err.Error())
		return ""
	}
	if len(typhoonActiveResp) == 0 {
		return ""
	}
	var lonSelf, latSelf = common.LocStr2float(common.MyConfig.Atmp.Loc)
	var message string
	for _, v := range typhoonActiveResp {
		//计算距离，大于1000公里跳过
		dis := common.Haversine(common.Str2Float64(v.Lng), common.Str2Float64(v.Lat), lonSelf, latSelf)
		if dis > 1000 {
			continue
		}

		//如果数据没变化不更新
		if v.Timeformate == lastUpdate {
			message += fmt.Sprintf(temp, v.Power, v.Strong, v.Name, v.Speed, dis, radiusHanlder(lastPoint.Radius7), lastPoint.Time[11:16])
			continue
		}
		lastUpdate = v.Timeformate

		//获取数据
		loc, lastes := TyphoonPath(v.Tfid)
		lastPoint = lastes
		//获取级别
		radius7 := radiusHanlder(lastes.Radius7)
		message += fmt.Sprintf(temp, v.Power, v.Strong, v.Name, v.Speed, dis, radius7, lastes.Time[11:16])
		radius12 := radiusHanlder(lastes.Radius12)
		radius10 := radiusHanlder(lastes.Radius10)
		resp := atmp.CreatePhoto(loc, dis, common.Str2Float64(radius7), common.Str2Float64(radius10), common.Str2Float64(radius12))
		if resp == nil {
			continue
		}
		telegram.SendPhoto(resp, v.Name, "")
	}
	return message
}

func radiusHanlder(data string) string {
	if data == "" {
		return ""
	}
	parts := strings.Split(data, "|")
	return parts[len(parts)-1]
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
