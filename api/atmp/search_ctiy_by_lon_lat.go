package atmp

import (
	"fmt"
	"services/common"
)

const (
	searchByLonLacUrl = "https://restapi.amap.com/v3/geocode/regeo?location=%s&key=%s"
	searchByAddrUrl   = "https://restapi.amap.com/v3/geocode/geo?address=%s&key=%s"
)

// 经纬度 市级 全名
var used = [3]string{}

// 搜搜地址
// 经纬度  全名
func SearchAddrs(addr, loc string) (string, string, bool) {
	if used[0] == loc && addr == "" {
		return used[0], used[2], used[1] == "广州市"
	}
	if addr == "" {
		return SearchByLonLac(loc)
	}
	return SearchByLonAddr(addr)
}

// 搜索经纬度
func SearchByLonAddr(addr string) (string, string, bool) {
	key := common.MyConfig.Atmp.Key
	url := fmt.Sprintf(searchByAddrUrl, addr, key)
	var resp SearchByAddrResp
	_, err := common.HttpRequest(common.MapType, common.GetType, url, nil, nil, false, &resp)
	if err != nil {
		common.Logger.Error(err.Error())
		return "", "", false
	}
	if len(resp.Geocodes) == 0 {
		common.Logger.Error(addr)
		return "", "", false
	}
	var content = resp.Geocodes[0]

	common.MyConfig.Home.Addr = ""
	used[0] = content.Location
	used[1], used[2] = content.City, content.FormattedAddress
	return used[0], used[2], used[1] == "广州市"
}

func SearchByLonLac(loc string) (string, string, bool) {
	key := common.MyConfig.Atmp.Key
	url := fmt.Sprintf(searchByLonLacUrl, loc, key)
	var resp SearchByLonLacResp
	_, err := common.HttpRequest(common.MapType, common.GetType, url, nil, nil, false, &resp)
	if err != nil {
		common.Logger.Error(err.Error())
		return "", "", false
	}
	used[0] = loc
	used[1], used[2] = resp.Regeocode.AddressComponent.City, resp.Regeocode.FormattedAddress
	return loc, used[2], used[1] == "广州市"
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
type SearchByAddrResp struct {
	Status   string `json:"status"`
	Info     string `json:"info"`
	Infocode string `json:"infocode"`
	Count    string `json:"count"`
	Geocodes []struct {
		FormattedAddress string `json:"formatted_address"`
		Country          string `json:"country"`
		Province         string `json:"province"`
		Citycode         string `json:"citycode"`
		City             string `json:"city"`
		District         string `json:"district"`
		//Township         []interface{} `json:"township"`
		//Neighborhood     struct {
		//	Name []interface{} `json:"name"`
		//	Type []interface{} `json:"type"`
		//} `json:"neighborhood"`
		//Building struct {
		//	Name []interface{} `json:"name"`
		//	Type []interface{} `json:"type"`
		//} `json:"building"`
		Adcode string `json:"adcode"`
		//Street   []interface{} `json:"street"`
		//Number   []interface{} `json:"number"`
		Location string `json:"location"`
		Level    string `json:"level"`
	} `json:"geocodes"`
}
