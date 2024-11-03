package service

import (
	"net/http"
	"services/common"
)

func WebService() {
	http.HandleFunc("/addr", updateAddr)
	var err = http.ListenAndServe(":6868", nil)
	if err != nil {
		common.Logger.Error(err.Error())
	}
}

func updateAddr(w http.ResponseWriter, r *http.Request) {
	// 判断参数是否是Get请求，并且参数解析正常
	if r.Method != "GET" {
		return
	}

	// 接收参数
	loc := r.FormValue("loc")
	if loc == "" {
		return
	}
	if !common.CheckLoc(loc) {
		return
	}
	common.MyConfig.Home.Loc = loc
	NewAddr <- true
}
