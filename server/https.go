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
	text := r.FormValue("addr")
	if text == "" {
		common.Logger.Error("null addr:" + text)
		return
	}

	//地址
	//检查是否临时
	//填入数据
	if text[0] == '-' {
		//isTem = true
		NewTempAddr <- text[1:]
	} else {
		common.MyConfig.Home.Loc, common.MyConfig.Home.Addr = common.CheckAddrOrLoc(text)
		NewAddr <- true
	}

	common.Logger.Warn("update addr" + text)
}
