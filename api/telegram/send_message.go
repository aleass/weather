package telegram

import (
	"fmt"
	"services/common"
	"time"
)

const (
	pullUrl = "https://api.telegram.org/bot%s/"
)

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

func GetAddress() (ok bool) {
	url := fmt.Sprintf(pullUrl+"getUpdates", common.MyConfig.Telegram.AddresToken)
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
	if messagedId == message.UpdateId || time.Now().Unix()-message.Message.Date > 15*60 {
		return
	}

	//填入数据
	common.MyConfig.Home.Loc, common.MyConfig.Home.Addr = common.CheckAddrOrLoc(message.Message.Text)
	messagedId = message.UpdateId
	ok = true
	return
}

type getUpdatesResp struct {
	Ok     bool `json:"ok"`
	Result []struct {
		UpdateId int `json:"update_id"`
		Message  struct {
			MessageId int `json:"message_id"`
			From      struct {
				Id           int    `json:"id"`
				IsBot        bool   `json:"is_bot"`
				FirstName    string `json:"first_name"`
				LastName     string `json:"last_name"`
				Username     string `json:"username"`
				LanguageCode string `json:"language_code"`
			} `json:"from"`
			Chat struct {
				Id        int    `json:"id"`
				FirstName string `json:"first_name"`
				LastName  string `json:"last_name"`
				Username  string `json:"username"`
				Type      string `json:"type"`
			} `json:"chat"`
			Date int64  `json:"date"`
			Text string `json:"text"`
		} `json:"message"`
	} `json:"result"`
}

type Msg struct {
	ChatId    int64  `json:"chat_id"`
	Text      string `json:"text"`
	Offset    int    `json:"offset"`
	ParseMode string `json:"parse_mode"`
}

type Response struct {
	Ok     bool `json:"ok"`
	Result struct {
		MessageId int64 `json:"message_id"`
		From      struct {
			Id        int64  `json:"id"`
			IsBot     bool   `json:"is_bot"`
			FirstName string `json:"first_name"`
			Username  string `json:"username"`
		} `json:"from"`
		Chat struct {
			Id        int    `json:"id"`
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
			Username  string `json:"username"`
			Type      string `json:"type"`
		} `json:"chat"`
		Date int    `json:"date"`
		Text string `json:"text"`
	} `json:"result"`
}
