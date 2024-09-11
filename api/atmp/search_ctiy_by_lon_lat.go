package atmp

import (
	"fmt"
	"weather/common"
)

const (
	searchByLonLacUrl = "https://restapi.amap.com/v3/geocode/regeo?location=%s&key=%s"
)

// 地区 全名 是否广州内
var used = [3]string{}

func SearchByLonLac(loc string) (string, string, bool) {
	if used[0] == loc {
		return used[1], used[2], true
	}

	key := common.MyConfig.Atmp.Key
	url := fmt.Sprintf(searchByLonLacUrl, loc, key)
	var resp SearchByLonLacResp
	_, err := common.HttpRequest(common.MapType, common.GetType, url, nil, nil, false, &resp)
	if err != nil {
		common.Logger.Error(err.Error())
		return "", "", false
	}
	var isOk = resp.Regeocode.AddressComponent.City == "广州市"
	used[0] = loc
	used[1], used[2] = resp.Regeocode.AddressComponent.District, resp.Regeocode.FormattedAddress
	return used[1], used[2], isOk
}

type SearchByLonLacResp struct {
	Status    string `json:"status"`
	Regeocode struct {
		AddressComponent struct {
			City         string `json:"city"`
			Province     string `json:"province"`
			Adcode       string `json:"adcode"`
			District     string `json:"district"`
			Towncode     string `json:"towncode"`
			StreetNumber struct {
				Number    string `json:"number"`
				Location  string `json:"location"`
				Direction string `json:"direction"`
				Distance  string `json:"distance"`
				Street    string `json:"street"`
			} `json:"streetNumber"`
			Country       string `json:"country"`
			Township      string `json:"township"`
			BusinessAreas []struct {
				Location string `json:"location"`
				Name     string `json:"name"`
				Id       string `json:"id"`
			} `json:"businessAreas"`
			Building struct {
				Name []interface{} `json:"name"`
				Type []interface{} `json:"type"`
			} `json:"building"`
			Neighborhood struct {
				Name []interface{} `json:"name"`
				Type []interface{} `json:"type"`
			} `json:"neighborhood"`
			Citycode string `json:"citycode"`
		} `json:"addressComponent"`
		FormattedAddress string `json:"formatted_address"`
	} `json:"regeocode"`
	Info     string `json:"info"`
	Infocode string `json:"infocode"`
}
