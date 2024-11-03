package he_feng

import (
	"fmt"
	"services/common"
)

const (
	LookupUrl = "https://geoapi.qweather.com/v2/city/lookup?location=%s&key=%s"
)

var used = [3]string{}

// 地区 全名 是否广州内
func Lookup() (string, string, bool) {
	if used[0] == common.MyConfig.Home.Loc {
		return used[1], used[2], true
	}

	url := fmt.Sprintf(LookupUrl, common.MyConfig.Home.Loc, common.MyConfig.HeFeng.Key)
	var cityRes CityResponse
	_, err := common.HttpRequest(common.WeatherType, common.GetType, url, nil, nil, false, &cityRes)
	if err != nil {
		common.Logger.Error(err.Error())
		return "", "", false
	}

	for _, v := range cityRes.Location {
		used[2] = v.Adm1 + v.Adm2 + v.Name
		if v.Adm2 == "广州" { //目前只支持广州
			used[0] = common.MyConfig.Home.Loc
			used[1] = v.Name
			return used[1], used[2], true
		}
	}
	return used[1], used[2], false
}

type CityResponse struct {
	Code     string `json:"code"`
	Location []struct {
		Name      string `json:"name"`
		Id        string `json:"id"`
		Lat       string `json:"lat"`
		Lon       string `json:"lon"`
		Adm2      string `json:"adm2"`
		Adm1      string `json:"adm1"`
		Country   string `json:"country"`
		Tz        string `json:"tz"`
		UtcOffset string `json:"utcOffset"`
		IsDst     string `json:"isDst"`
		Type      string `json:"type"`
		Rank      string `json:"rank"`
		FxLink    string `json:"fxLink"`
	} `json:"location"`
	Refer struct {
		Sources []string `json:"sources"`
		License []string `json:"license"`
	} `json:"refer"`
}
