package service

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"weather/common"
)

type IPResponse struct {
	//Ip                 string      `json:"ip"`
	//Network            string      `json:"network"`
	//Version            string      `json:"version"`
	City string `json:"city"`
	//Region             string      `json:"region"`
	//RegionCode         string      `json:"region_code"`
	//Country            string      `json:"country"`
	//CountryName        string      `json:"country_name"`
	//CountryCode        string      `json:"country_code"`
	//CountryCodeIso3    string      `json:"country_code_iso3"`
	//CountryCapital     string      `json:"country_capital"`
	//CountryTld         string      `json:"country_tld"`
	//ContinentCode      string      `json:"continent_code"`
	//InEu               bool        `json:"in_eu"`
	//Postal             interface{} `json:"postal"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	//Timezone           string      `json:"timezone"`
	//UtcOffset          string      `json:"utc_offset"`
	//CountryCallingCode string      `json:"country_calling_code"`
	//Currency           string      `json:"currency"`
	//CurrencyName       string      `json:"currency_name"`
	//Languages          string      `json:"languages"`
	//CountryArea        int         `json:"country_area"`
	//CountryPopulation  int         `json:"country_population"`
	//Asn                string      `json:"asn"`
	//Org                string      `json:"org"`
}

var localMap = map[string]*UrlInfo{}

func HttpRun() {
	r := gin.Default()
	r.GET("/", func(context *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				common.LogSend(fmt.Sprintf("panic err:%v", err), common.PanicType)
			}
		}()
		var url = fmt.Sprintf("https://ipapi.co/%s/json/", context.ClientIP())
		res, err := http.Get(url)
		if err != nil {
			context.JSON(200, err.Error()+url)
		}
		if res == nil {
			println("res == nil")
		} else {
			ipRes := IPResponse{}
			data, err := io.ReadAll(res.Body)
			if err != nil {
				context.JSON(200, err.Error())
			}
			if err = json.Unmarshal(data, &ipRes); err != nil {
				context.JSON(200, err.Error())
			}
			ipRes.Latitude = 23.118100
			ipRes.Longitude = 113.253900
			name := context.Query("name")
			weachetNote := context.Query("note")
			info, ok := localMap[name]
			if !ok {
				localMap["guan"] = &UrlInfo{
					name:      name + ":" + ipRes.City,
					caiYunUrl: fmt.Sprintf(caiYunUrl, myConfig.CaiYun.Token, fmt.Sprintf("%f,%f", ipRes.Longitude, ipRes.Latitude)),
					weChatUrl: wechatUrl + notes[weachetNote],
					_switch:   make(chan struct{}),
					watchTime: 5, //默认10分钟
				}
				info = localMap["guan"]
				go info.WatchWeather()
			} else {
				localMap[name].name = name + ":" + ipRes.City
				localMap[name].caiYunUrl = fmt.Sprintf(caiYunUrl, "TAkhjf8d1nlSlspN", fmt.Sprintf("%f,%f", ipRes.Longitude, ipRes.Latitude))
			}
			context.JSON(200, info.name)
		}
	})
	r.Run(":8080")
}
