package typhoon

import (
	"fmt"
	"weather/common"
)

const (
	typhoonPathUrl = "https://typhoon.slt.zj.gov.cn/Api/TyphoonInfo/"
)

func TyphoonPath(tyId string) ([]string, *PointsInfo) {
	var typhoonResp TyphoonResp
	var url = typhoonPathUrl + tyId
	_, err := common.HttpRequest(common.OtherType, common.GetType, url, nil, header, false, &typhoonResp)
	if err != nil {
		common.Logger.Error(err.Error())
		return nil, nil
	}
	if len(typhoonResp.Points) == 0 {
		return nil, nil
	}
	var loctions []string
	for _, s := range typhoonResp.Points {
		loctions = append(loctions, fmt.Sprintf("%s,%s", s.Lng, s.Lat))
	}
	return loctions, &typhoonResp.Points[len(typhoonResp.Points)-1]
}

type TyphoonResp struct {
	Tfid      string        `json:"tfid"`
	Name      string        `json:"name"`
	Enname    string        `json:"enname"`
	Isactive  string        `json:"isactive"`
	Starttime string        `json:"starttime"`
	Endtime   string        `json:"endtime"`
	Warnlevel string        `json:"warnlevel"`
	Centerlng string        `json:"centerlng"`
	Centerlat string        `json:"centerlat"`
	Land      []interface{} `json:"land"`
	Points    []PointsInfo  `json:"points"`
}
type PointsInfo struct {
	Time          string `json:"time"`
	Lng           string `json:"lng"`
	Lat           string `json:"lat"`
	Strong        string `json:"strong"`
	Power         string `json:"power"`
	Speed         string `json:"speed"`
	Pressure      string `json:"pressure"`
	Movespeed     string `json:"movespeed"`
	Movedirection string `json:"movedirection"`
	Radius7       string `json:"radius7"`
	Radius10      string `json:"radius10"`
	Radius12      string `json:"radius12"`
	//Forecast      []struct {
	//	Tm             string `json:"tm"`
	//	Forecastpoints []struct {
	//		Time     string    `json:"time"`
	//		Lng      string    `json:"lng"`
	//		Lat      string    `json:"lat"`
	//		Strong   string    `json:"strong"`
	//		Power    string    `json:"power"`
	//		Speed    string    `json:"speed"`
	//		Pressure string    `json:"pressure"`
	//		Tm       string    `json:"tm,omitempty"`
	//		Ybsj     time.Time `json:"ybsj,omitempty"`
	//	} `json:"forecastpoints"`
	//} `json:"forecast"`
	Ckposition *string `json:"ckposition"`
	Jl         *string `json:"jl"`
}
