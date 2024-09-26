package service

import (
	"fmt"
	"time"
	"weather/api/juhe"
	"weather/api/tian"
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
		var now = time.Now()
		var h = now.Hour()

		switch {
		case h < 6:
			goto sleep
		case h > 19:
			time.Sleep(selectTime)
		}

		juhe.GetNews()
		tian.GetNews()

	sleep:
		time.Sleep(selectTime)
	}
}
