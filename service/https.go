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

var timeDump = make(chan struct{}, 1)
var localMap = map[string]*UrlInfo{}

func HttpRun() {
	r := gin.Default()
	r.GET("/", func(context *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				common.LogSend(fmt.Sprintf("panic err:%v", err), common.PanicType)
			}
		}()
		var ip = context.ClientIP()
		var url = fmt.Sprintf("https://ipapi.co/%s/json/", ip)
		res, err := http.Get(url)
		if err != nil {
			context.JSON(200, err.Error()+url)
		}
		if res == nil {
			context.JSON(200, "res == nil")
			return
		} else {
			ipRes := IPResponse{}
			data, err := io.ReadAll(res.Body)
			if err != nil {
				context.JSON(200, err.Error())
				return
			}
			if err = json.Unmarshal(data, &ipRes); err != nil {
				context.JSON(200, err.Error())
				return
			}
			if ipRes.Latitude == 0 || ipRes.Longitude == 0 {
				context.JSON(200, "经纬度找不到")
				return
			}

			name := context.Query("name")
			weachetNote := context.Query("note")
			info, ok := localMap[name]
			if !ok {
				info = &UrlInfo{
					name:        name,
					caiYunUrl:   fmt.Sprintf(caiYunUrl, myConfig.CaiYun.Token, fmt.Sprintf("%f,%f", ipRes.Longitude, ipRes.Latitude)),
					weChatUrl:   wechatUrl + notes[weachetNote],
					_switch:     make(chan struct{}),
					watchTime:   5, //默认10分钟
					isUrlConfig: true,
				}
				timeDump <- struct{}{}
				localMap[name] = info
				<-timeDump
				go info.WatchWeather()
			} else {
				localMap[name].name = name + ":" + ipRes.City
				localMap[name].caiYunUrl = fmt.Sprintf(caiYunUrl, myConfig.CaiYun.Token, fmt.Sprintf("%f,%f", ipRes.Longitude, ipRes.Latitude))
			}
			context.JSON(200, fmt.Sprintf("name:%s,ip:%s,经度:%f,纬度:%f", info.name, ip, ipRes.Longitude, ipRes.Latitude))
		}
	})
	r.Run(":8080")
}
