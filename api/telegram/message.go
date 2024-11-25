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
func GetMessage() (ok, isTem bool) {
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

	message := resp.Result[l-1]
	if messagedId == message.UpdateId || time.Now().Unix()-message.Message.Date > 15*60 || message.Message.Text == "" {
		return
	}

	//检查是否临时
	var text = message.Message.Text
	//填入数据
	if text[0] == '-' {
		common.MyConfig.TemHome.Loc, common.MyConfig.TemHome.Addr = common.CheckAddrOrLoc(text[1:])
		isTem = true
	} else {
		common.MyConfig.Home.Loc, common.MyConfig.Home.Addr = common.CheckAddrOrLoc(text)
	}

	messagedId = message.UpdateId
	ok = true
	return
}
