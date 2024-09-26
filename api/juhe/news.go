package juhe

import (
	"fmt"
	"os"
	"time"
	"weather/api/telegram"
	"weather/common"
)

const (
	newUrl = "http://v.juhe.cn/toutiao/index?key="
)

var (
	lastUniquekey = ""
	category      = map[string]struct{}{
		"头条": {},
		"国内": {},
		"国际": {},
		"军事": {},
		"科技": {},
		"财经": {},
		"汽车": {},
		"健康": {},
	}
)

func GetNews() {
	if lastUniquekey == "" {
		bytes, _ := os.ReadFile("new_key")
		lastUniquekey = string(bytes)
	}

	var newsResponse NewsResponse
	//var ss,_ = os.ReadFile("temp.txt")
	//json.Unmarshal(ss, &newsResponse)
	raw, err := common.HttpRequest(common.WeatherType, common.GetType, newUrl+common.MyConfig.JuHe.Key, nil, nil, false, &newsResponse)
	if err != nil {
		common.Logger.Error(err.Error())
		return
	}
	os.WriteFile("/Users/tuski/code/src/weather/temp.txt", raw, 0777)

	if len(newsResponse.Result.Data) == 0 {
		common.Logger.Error(newsResponse.Reason)
		return
	}

	var msgList []string
	for _, datum := range newsResponse.Result.Data {
		if datum.Uniquekey == lastUniquekey {
			break
		}

		if _, ok := category[datum.Category]; !ok {
			continue
		}

		var newMsg = fmt.Sprintf("[%s](%s)\n\n**来源：%s（%s）**\n**时间：%s**", datum.Title, datum.Url, datum.AuthorName, datum.Category, datum.Date[11:16])
		msgList = append(msgList, newMsg)

		//if  datum.ThumbnailPicS != "" {
		//	resp, err := common.HttpRequest(common.OtherType, common.GetType, datum.ThumbnailPicS, nil, header, false, nil)
		//	if err != nil {
		//		common.Logger.Error(err.Error())
		//		return
		//	}
		//	telegram.SendPhoto(bytes.NewReader(resp),"",common.Int642Str(messageId))
		//}
	}

	if lastUniquekey != newsResponse.Result.Data[0].Uniquekey {
		lastUniquekey = newsResponse.Result.Data[0].Uniquekey
		os.WriteFile("new_key", []byte(lastUniquekey), 0777)
	}

	telegram.SendMessage("————————————————————"+time.Now().Format("2006-01-02 15:04:05"), common.MyConfig.Telegram.Token2)
	for i := len(msgList) - 1; i >= 0; i-- {
		telegram.SendMessage(msgList[i], common.MyConfig.Telegram.Token2)
	}
	return
}

type NewsResponse struct {
	Reason string `json:"reason"`
	Result struct {
		//Stat string `json:"stat"`
		Data []struct {
			Uniquekey  string `json:"uniquekey"`
			Title      string `json:"title"`
			Date       string `json:"date"`
			Category   string `json:"category"`
			AuthorName string `json:"author_name"`
			Url        string `json:"url"`
			//ThumbnailPicS   string `json:"thumbnail_pic_s"`
			//ThumbnailPicS02 string `json:"thumbnail_pic_s02,omitempty"`
			//IsContent       string `json:"is_content"`
		} `json:"data"`
		//Page     string `json:"page"`
		//PageSize string `json:"pageSize"`
	} `json:"result"`
	//ErrorCode int `json:"error_code"`
}
