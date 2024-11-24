package he_feng

import (
	"fmt"
	"services/common"
)

const (
	realTimeRain = host + "/v7/minutely/5m?location=%s&key=%s"
)

var titleTemp = common.SubStr + `%s
`

var nextTemp = common.SubStr + `%s  %s
`

func FiveMinRain(loc string) string {
	url := fmt.Sprintf(realTimeRain, loc, common.MyConfig.HeFeng.Key)
	var fiveMinRainRes FiveMinRainRes
	_, err := common.HttpRequest(common.WeatherType, common.GetType, url, nil, nil, false, &fiveMinRainRes)
	if err != nil {
		common.Logger.Error(err.Error())
		return ""
	}
	if len(fiveMinRainRes.Minutely) == 0 || fiveMinRainRes.Summary == "未来两小时无降水" {
		return ""
	}

	//var title = fmt.Sprintf(titleTemp, fiveMinRainRes.Summary, fiveMinRainRes.UpdateTime[11:16])
	var title = fmt.Sprintf(titleTemp, fiveMinRainRes.Summary)
	var lastTime string
	var max, curr = 5, 0
	for _, s := range fiveMinRainRes.Minutely {
		if s.Type != "rain" || s.Precip == "0.00" {
			continue
		}
		//过滤同十分钟
		var now = s.FxTime[11:15]
		if now == lastTime {
			continue
		}
		lastTime = now
		title += fmt.Sprintf(nextTemp, s.FxTime[11:16], s.Precip)
		curr++
		if max > curr {
			break
		}
	}
	return "【五分钟降雨预报】\n" + title + "\n"
}

type FiveMinRainRes struct {
	Code       string `json:"code"`
	UpdateTime string `json:"updateTime"`
	FxLink     string `json:"fxLink"`
	Summary    string `json:"summary"`
	Minutely   []struct {
		FxTime string `json:"fxTime"`
		Precip string `json:"precip"`
		Type   string `json:"type"`
	} `json:"minutely"`
	Refer struct {
		Sources []string `json:"sources"`
		License []string `json:"license"`
	} `json:"refer"`
}
