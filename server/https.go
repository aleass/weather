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
		common.Logger.Error("null addr:" + loc)
		return
	}

	//地址
	common.MyConfig.Home.Loc, common.MyConfig.Home.Addr = common.CheckAddrOrLoc(loc)
	if common.MyConfig.Home.Addr != "" {
		goto end
	}

	//经纬度
	if !common.CheckLoc(loc) {
		common.Logger.Error("error addr:" + loc)
		return
	}
	common.MyConfig.Home.Loc = loc

end:
	common.Logger.Warn("update addr" + loc)
	NewAddr <- true
}
