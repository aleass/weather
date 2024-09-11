package telegram

import (
	"fmt"
	"io"
	"weather/common"
)

const (
	pullUrl = "https://api.telegram.org/bot%s/"
)

func SendPhoto(file io.Reader, name string) {
	url := fmt.Sprintf(pullUrl+"sendPhoto", common.MyConfig.Telegram.Token)
	formInfo := [][2]string{
		{"chat_id", common.Int642Str(common.MyConfig.Telegram.ChatId)},
		{"caption", name},
	}
	err := common.UploadFile(url, formInfo, file, "photo")
	if err != nil {
		common.Logger.Error(err.Error())
		return
	}
}

type Msg struct {
	ChatId int64  `json:"chat_id"`
	Text   string `json:"text"`
	Offset int    `json:"offset"`
}
