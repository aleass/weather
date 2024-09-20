package he_feng

import (
	"fmt"
	"weather/common"
)

const (
	realTimeRain = host + "/v7/minutely/5m?location=%s&key=%s"
)

const titleTemp = common.SubStr + `%s
`

var nextTemp = common.SubStr + `%s  %s
`

func FiveMinRain() string {
	url := fmt.Sprintf(realTimeRain, common.MyConfig.Atmp.Loc, common.MyConfig.HeFeng.Key)
	var fiveMinRainRes FiveMinRainRes
	_, err := common.HttpRequest(common.WeatherType, common.GetType, url, nil, nil, false, &fiveMinRainRes)
	if err != nil {
		common.Logger.Error(err.Error())
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
	return "【五分钟降雨预报】\n" + title
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
