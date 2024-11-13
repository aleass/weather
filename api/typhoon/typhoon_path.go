package typhoon

import (
	"fmt"
	"services/common"
	"time"
)

const (
	typhoonPathUrl = "https://typhoon.slt.zj.gov.cn/Api/TyphoonInfo/"
)

func TyphoonPath(tyId string) ([]string, []string, [][2]string, [][2]string, *PointsInfo) {
	var typhoonResp TyphoonResp
	var url = typhoonPathUrl + tyId
	_, err := common.HttpRequest(common.OtherType, common.GetType, url, nil, common.Header, false, &typhoonResp)
	if err != nil {
		common.Logger.Error(err.Error())
		return nil, nil, nil, nil, nil
	}
	if len(typhoonResp.Points) == 0 {
		return nil, nil, nil, nil, nil
	}
	var (
		loctions      []string
		forecasts     []string
		forecastsName [][2]string
		forecastsDate [][2]string
	)

	var l = len(typhoonResp.Points) - 1
	for i, s := range typhoonResp.Points {
		loctions = append(loctions, fmt.Sprintf("%s,%s", s.Lng, s.Lat))
		if i == l {
			var max, maxLen int
			for j, forecast := range s.Forecast {
				l = len(forecast.Forecastpoints)
				if l == 0 {
					continue
				}
				for _, points := range forecast.Forecastpoints {
					forecasts = append(forecasts, fmt.Sprintf("%s,%s", points.Lng, points.Lat))
				}
				if len(forecast.Forecastpoints) > max {
					max = len(forecast.Forecastpoints)
					maxLen = j
				}

				index := l - 1
				forecastsName = append(forecastsName, [2]string{forecast.Forecastpoints[index].Tm, fmt.Sprintf("%s,%s", forecast.Forecastpoints[index].Lng, forecast.Forecastpoints[index].Lat)})
				for ; index >= 0; index-- {
					forecasts = append(forecasts, fmt.Sprintf("%s,%s", forecast.Forecastpoints[index].Lng, forecast.Forecastpoints[index].Lat))
				}
			}

			for j, points := range s.Forecast[maxLen].Forecastpoints[1:] {
				if len(forecastsDate) == 9 && j != max-1 {
					continue
				}
				forecastsDate = append(forecastsDate, [2]string{points.Time[8:13], fmt.Sprintf("%s,%s", points.Lng, points.Lat)})
			}

		}
	}
	return loctions, forecasts, forecastsName, forecastsDate, &typhoonResp.Points[len(typhoonResp.Points)-1]
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
	Forecast      []struct {
		Tm             string `json:"tm"`
		Forecastpoints []struct {
			Time     string    `json:"time"`
			Lng      string    `json:"lng"`
			Lat      string    `json:"lat"`
			Strong   string    `json:"strong"`
			Power    string    `json:"power"`
			Speed    string    `json:"speed"`
			Pressure string    `json:"pressure"`
			Tm       string    `json:"tm,omitempty"`
			Ybsj     time.Time `json:"ybsj,omitempty"`
		} `json:"forecastpoints"`
	} `json:"forecast"`
	Ckposition *string `json:"ckposition"`
	Jl         *string `json:"jl"`
}
