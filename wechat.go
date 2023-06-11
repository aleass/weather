package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

type Text struct {
	Msgtype string  `json:"msgtype"`
	Text    Content `json:"text"`
}
type Content struct {
	Content string `json:"content"`
}

var message = Text{
	Msgtype: "text",
	Text:    Content{},
}

func Send(msg, url string) {
	println(msg)
	message.Text.Content = msg
	muta, _ := json.Marshal(message)
	req, err := http.NewRequest("POST", url, bytes.NewReader(muta))
	if err != nil {
		log.Println("err:" + err.Error())
		return
	}
	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		log.Println("err:" + err.Error())
	}
}
