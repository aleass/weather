package telegram

import (
	"fmt"
	"io"
	"weather/common"
)

func SendPhoto(file io.Reader, name string, messageId string) {
	url := fmt.Sprintf(pullUrl+"sendPhoto", common.MyConfig.Telegram.Token2)
	formInfo := [][2]string{
		{"chat_id", common.Int642Str(common.MyConfig.Telegram.ChatId)},
		{"caption", name},
		{"parse_mode", "Markdown"},
		{"reply_to_message_id", messageId},
	}
	err := common.UploadFile(url, formInfo, file, "photo")
	if err != nil {
		common.Logger.Error(err.Error())
		return
	}
}
