package he_feng

import (
	"fmt"
	"services/common"
)

const (
	weatherUrl = host + "/v7/weather/now?location=%s&key=%s"
)

// 地区 全名 是否广州内
func WeatherInfo(loc string) string {
	url := fmt.Sprintf(weatherUrl, loc, common.MyConfig.HeFeng.Key)
	var weatherResp WeatherResp
	_, err := common.HttpRequest(common.WeatherType, common.GetType, url, nil, nil, false, &weatherResp)
	if err != nil {
		common.Logger.Error(err.Error())
		return ""
	}

	var weather = "【天气】\n%s%s° %s %s%sm/s 能见度：%s 湿度：%s%%\n\n"

	return fmt.Sprintf(weather, common.SubStr, weatherResp.Now.Temp, weatherResp.Now.Text, weatherResp.Now.WindDir, weatherResp.Now.WindSpeed, weatherResp.Now.Vis, weatherResp.Now.Humidity)
}

type WeatherResp struct {
	Code       string `json:"code"`
	UpdateTime string `json:"updateTime"`
	FxLink     string `json:"fxLink"`
	Now        struct {
		ObsTime   string `json:"obsTime"`
		Temp      string `json:"temp"`
		FeelsLike string `json:"feelsLike"`
		Icon      string `json:"icon"`
		Text      string `json:"text"`
		Wind360   string `json:"wind360"`
		WindDir   string `json:"windDir"`
		WindScale string `json:"windScale"`
		WindSpeed string `json:"windSpeed"`
		Humidity  string `json:"humidity"`
		Precip    string `json:"precip"`
		Pressure  string `json:"pressure"`
		Vis       string `json:"vis"`
		Cloud     string `json:"cloud"`
		Dew       string `json:"dew"`
	} `json:"now"`
	Refer struct {
		Sources []string `json:"sources"`
		License []string `json:"license"`
	} `json:"refer"`
}
