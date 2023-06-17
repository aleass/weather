package common

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

var (
	ipUrls = [...]string{"https://qifu-api.baidubce.com/ip/geo/v1/district?ip=%s", "https://ipapi.co/%s/json/", "http://ip-api.com/json/%s?lang=zh-CN"}
)

type IPResponse struct {
	Latitude   float64 `json:"latitude" desc:"ipapi"`
	Longitude  float64 `json:"longitude" desc:"ipapi"`
	Latitude1  float64 `json:"lat" desc:"ip-api"`
	Longitude1 float64 `json:"lon" desc:"ip-api"`
}

type BaiduIpRes struct {
	Data struct {
		//Continent string `json:"continent"`
		//Country   string `json:"country"`
		//Zipcode   string `json:"zipcode"`
		//Timezone  string `json:"timezone"`
		//Accuracy  string `json:"accuracy"`
		//Owner     string `json:"owner"`
		//Isp       string `json:"isp"`
		//Source    string `json:"source"`
		//Areacode  string `json:"areacode"`
		//Adcode    string `json:"adcode"`
		//Asnumber  string `json:"asnumber"`
		Lat string `json:"lat"`
		Lng string `json:"lng"`
		//Radius    string `json:"radius"`
		//Prov      string `json:"prov"`
		//City      string `json:"city"`
		//District  string `json:"district"`
	} `json:"data"`
	//Charge bool `json:"charge"`
	Code string `json:"code"`
}

func GetIpAddress(ip string) (string, error) {
	var (
		ipRes      = IPResponse{}
		ipBaiduRes = BaiduIpRes{}
	)
	for i, url := range ipUrls {
		url = fmt.Sprintf(url, ip)
		res, err := http.Get(url)
		if err != nil {
			return "", err
		}
		data, err := io.ReadAll(res.Body)
		if err != nil {
			return "", err
		}
		if i == 0 {
			if err = json.Unmarshal(data, &ipBaiduRes); err != nil {
				return "", err
			}
			//百度
			if ipBaiduRes.Code == "Success" && ipBaiduRes.Data.Lng != "" && ipBaiduRes.Data.Lat != "" {
				return fmt.Sprintf("%s,%s", ipBaiduRes.Data.Lng, ipBaiduRes.Data.Lat), nil
			}
			continue
		}

		if err = json.Unmarshal(data, &ipRes); err != nil {
			return "", err
		}

		//ipapi
		if ipRes.Longitude != 0 && ipRes.Latitude != 0 {
			return fmt.Sprintf("%f,%f", ipRes.Longitude, ipRes.Latitude), nil
		}
		//ip-api
		if ipRes.Latitude1 != 0 && ipRes.Longitude1 != 0 {
			return fmt.Sprintf("%f,%f", ipRes.Longitude1, ipRes.Latitude1), nil
		}
	}
	return "", nil
}
