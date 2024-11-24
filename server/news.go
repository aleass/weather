package service

import (
	"services/api/juhe"
	"services/api/tian"
	"services/common"
	"time"
)

func NewsRun(selectTime time.Duration) {
	defer common.RecoverWithStackTrace(RunWeather, selectTime)

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
