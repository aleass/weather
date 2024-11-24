package he_feng

import (
	"fmt"
	"services/common"
)

const (
	warningUrl = host + "/v7/warning/now?key=%s&location=%s"
)

// 地区 全名 是否广州内
func CityWarning(loc string) (string, string) {
	url := fmt.Sprintf(warningUrl, common.MyConfig.HeFeng.Key, loc)
	var warningRes WarningResp
	_, err := common.HttpRequest(common.WeatherType, common.GetType, url, nil, nil, false, &warningRes)
	if err != nil {
		common.Logger.Error(err.Error())
		return "", ""
	}
	if len(warningRes.Warning) == 0 {
		return "", ""
	}

	var title, text = "【预警】", ""
	for _, s := range warningRes.Warning {
		title += s.Title + " "
		text += s.Text + "\n\n"
	}
	title += "\n\n"
	return title, text
}

type WarningResp struct {
	Code       string `json:"code"`
	UpdateTime string `json:"updateTime"`
	FxLink     string `json:"fxLink"`
	Warning    []struct {
		Id            string `json:"id"`
		Sender        string `json:"sender"`
		PubTime       string `json:"pubTime"`
		Title         string `json:"title"`
		StartTime     string `json:"startTime"`
		EndTime       string `json:"endTime"`
		Status        string `json:"status"`
		Level         string `json:"level"`
		Severity      string `json:"severity"`
		SeverityColor string `json:"severityColor"`
		Type          string `json:"type"`
		TypeName      string `json:"typeName"`
		Urgency       string `json:"urgency"`
		Certainty     string `json:"certainty"`
		Text          string `json:"text"`
		Related       string `json:"related"`
	} `json:"warning"`
	Refer struct {
		Sources []string `json:"sources"`
		License []string `json:"license"`
	} `json:"refer"`
}
