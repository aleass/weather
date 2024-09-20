package service

import (
	"fmt"
	"time"
	"weather/api/juhe"
	"weather/common"
)

func NewsRun(selectTime time.Duration) {
	defer func() {
		if err := recover(); err != nil {
			common.LogSend(fmt.Sprintf("panic err:%v", err), common.PanicType)
		}
		time.Sleep(selectTime)
		go NewsRun(selectTime)
	}()
	for {
		juhe.GetNews()
		time.Sleep(selectTime)
	}
}
