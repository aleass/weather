package telegram

import (
	"fmt"
	"services/common"
	"time"
)

func WebHoods() (ok bool) {
	url := fmt.Sprintf(pullUrl+"setWebhook", common.MyConfig.Telegram.AddresToken)
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
	if messagedId == message.UpdateId || time.Now().Unix()-message.Message.Date > 5*60 {
		return
	}

	common.MyConfig.Home.Loc = message.Message.Text
	messagedId = message.UpdateId
	ok = true
	return
}

type webHoodsParamet struct {
	Url       string `json:"url"`
	Text      string `json:"text"`
	Offset    int    `json:"offset"`
	ParseMode string `json:"parse_mode"`
}
