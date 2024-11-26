package telegram

import (
	"fmt"
	"services/common"
	"time"
)

const (
	pullUrl = "https://api.telegram.org/bot%s/"
)

// 发生信息
func SendMessage(info, token string) int64 {
	url := fmt.Sprintf(pullUrl+"sendMessage", token)
	var msg = Msg{
		ChatId:    common.MyConfig.Telegram.ChatId,
		Text:      info,
		ParseMode: "Markdown",
	}
	var resp Response
	_, err := common.HttpRequest(common.OtherType, common.PostType, url, msg, nil, true, &resp)
	if err != nil {
		common.Logger.Error(err.Error())
		return 0
	}
	if !resp.Ok {
		common.Logger.Error(resp.Result.Text)
	}
	return resp.Result.MessageId
}

var messagedId = 0

// 接受信息
func GetMessage(acceptChan chan string) (ok bool) {
	url := fmt.Sprintf(pullUrl+"getUpdates", common.MyConfig.Telegram.WeatherToken)
	var msg = Msg{
		Offset: messagedId,
	}
	var resp getUpdatesResp
	_, err := common.HttpRequest(common.OtherType, common.PostType, url, msg, nil, true, &resp)
	if err != nil {
		common.Logger.Error(err.Error())
		return
	}
	var l = len(resp.Result)
	if l == 0 {
		return
	}

	for _, message := range resp.Result {
		if messagedId >= message.UpdateId || time.Now().Unix()-message.Message.Date > 5*60+1 || message.Message.Text == "" {
			continue
		}
		//检查是否临时
		var text = message.Message.Text
		//填入数据
		if text[0] == '-' {
			acceptChan <- text[1:]
		} else {
			common.MyConfig.Home.Loc, common.MyConfig.Home.Addr = common.CheckAddrOrLoc(text)
			ok = true
		}
		messagedId = message.UpdateId
	}

	return
}
