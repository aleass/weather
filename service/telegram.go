package service

import (
	"fmt"
	"weather/common"
)

const (
	pullUrl = "https://api.telegram.org/bot%s/"
)

func SendMessage(info string) {
	url := fmt.Sprintf(pullUrl+"sendMessage", common.MyConfig.Telegram.Token)
	var msg = Msg{
		ChatId: common.MyConfig.Telegram.ChatId,
		Text:   info,
	}
	var resp Response
	_, err := common.HttpRequest(common.OtherType, common.PostType, url, msg, nil, true, &resp)
	if err != nil {
		common.Logger.Error(err.Error())
		return
	}
	if !resp.Ok {
		common.Logger.Error(resp.Result.Text)
	}
}

type Msg struct {
	ChatId int64  `json:"chat_id"`
	Text   string `json:"text"`
	Offset int    `json:"offset"`
}

type Response struct {
	Ok     bool `json:"ok"`
	Result struct {
		MessageId int `json:"message_id"`
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

var messagedId = 0

func GetAddress() {
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
	for _, s := range resp.Result {
		if messagedId == s.UpdateId {
			continue
		}
		common.MyConfig.Atmp.Loc = s.Message.Text
		messagedId = s.UpdateId
	}
	return
}

type getUpdatesResp struct {
	Ok     bool `json:"ok"`
	Result []struct {
		UpdateId int `json:"update_id"`
		Message  struct {
			MessageId int `json:"message_id"`
			//From      struct {
			//	Id           int    `json:"id"`
			//	IsBot        bool   `json:"is_bot"`
			//	FirstName    string `json:"first_name"`
			//	LastName     string `json:"last_name"`
			//	Username     string `json:"username"`
			//	LanguageCode string `json:"language_code"`
			//} `json:"from"`
			//Chat struct {
			//	Id        int    `json:"id"`
			//	FirstName string `json:"first_name"`
			//	LastName  string `json:"last_name"`
			//	Username  string `json:"username"`
			//	Type      string `json:"type"`
			//} `json:"chat"`
			//Date int    `json:"date"`
			Text string `json:"text"`
		} `json:"message"`
	} `json:"result"`
}
