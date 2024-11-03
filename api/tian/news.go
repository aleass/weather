package tian

import (
	"fmt"
	"os"
	"services/api/telegram"
	"services/common"
)

const (
	newUrl = "https://apis.tianapi.com/guonei/index?num=50&key="
)

var (
	lastUniquekey = ""
)

func GetNews() {
	if lastUniquekey == "" {
		bytes, _ := os.ReadFile(common.FileKeyPath + "new2_key")
		lastUniquekey = string(bytes)
	}

	var newsResponse NewsResponse
	//var ss,_ = sysos.ReadFile("temp2.txt")
	//json.Unmarshal(ss, &newsResponse)
	raw, err := common.HttpRequest(common.WeatherType, common.GetType, newUrl+common.MyConfig.Tian.Key, nil, nil, false, &newsResponse)
	if err != nil {
		common.Logger.Error(err.Error())
		return
	}
	os.WriteFile(common.FileKeyPath+"temp2.txt", raw, 0777)

	if len(newsResponse.Result.Newslist) == 0 {
		common.Logger.Error(newsResponse.Msg)
		return
	}

	var msgList []string
	for _, datum := range newsResponse.Result.Newslist {
		if datum.Id == lastUniquekey {
			break
		}

		var newMsg = fmt.Sprintf("[%s](%s)\n\n**来源：%s**\n**时间：%s**", datum.Title, datum.Url, datum.Source, datum.Ctime) //2024-09-26 00:00
		msgList = append(msgList, newMsg)
	}

	if lastUniquekey != newsResponse.Result.Newslist[0].Id {
		lastUniquekey = newsResponse.Result.Newslist[0].Id
		os.WriteFile(common.FileKeyPath+"new2_key", []byte(lastUniquekey), 0777)
	}

	for i := len(msgList) - 1; i >= 0; i-- {
		telegram.SendMessage(msgList[i], common.MyConfig.Telegram.Token2)
	}
	return
}

type NewsResponse struct {
	Code   int    `json:"code"`
	Msg    string `json:"msg"`
	Result struct {
		Curpage  int `json:"curpage"`
		Allnum   int `json:"allnum"`
		Newslist []struct {
			Id          string `json:"id"`
			Ctime       string `json:"ctime"`
			Title       string `json:"title"`
			Description string `json:"description"`
			Source      string `json:"source"`
			PicUrl      string `json:"picUrl"`
			Url         string `json:"url"`
		} `json:"newslist"`
	} `json:"result"`
}
