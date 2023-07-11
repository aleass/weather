package service

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"time"
	"weather/common"
)

const (
	expressageTypeUrl = "http://www.kuaidi.com/index-ajaxselectinfo-%s.html"
	expressageMsgUrl  = "http://www.kuaidi.com/index-ajaxselectcourierinfo-%s-%s-KUAIDICODE%d.html"
)

func GetExpressage(c *gin.Context) {
	nu := c.Query("nu")
	_type := c.Query("type")
	companyType, ok := companyTypes[_type]
	if !ok {
		c.JSON(200, "type  不存在"+_type)
		return
	}
	var msg map[string]companyInfo
	url := fmt.Sprintf(expressageTypeUrl, nu)
	res, err := http.Get(url)
	if err != nil {
		println(err.Error())
		return
	}
	data, err := io.ReadAll(res.Body)
	if err != nil {
		println(err.Error())
		return
	}
	json.Unmarshal(data, &msg)
	now := time.Now().Unix()
	for _, v := range msg {
		if companyType == v.Exname {
			go GetExpressageData(fmt.Sprintf(expressageMsgUrl, nu, v.Exname, now))
			c.JSON(200, "success")
			return
		}
	}
	c.JSON(200, "type  不存在"+_type)
}

func GetExpressageData(url string) {
	var status = -1
	var index = 0
	var unit = time.Hour
	for {
		data, err := common.HttpRequest(http.MethodPost, url, nil, [][2]string{{"Referer", "http://www.kuaidi.com/"}})
		if err != nil {
			println(err.Error())
			return
		}
		var msg expressageMsg
		json.Unmarshal(data, &msg)
		if !msg.Success {
			goto end
		}
		status = msg.Status
		for i := len(msg.Data) - 1; i >= index; i-- {
			common.Send(msg.Data[i].Time+" "+msg.Data[i].Context, GetWechatUrl("mine"))
		}
		index = len(msg.Data) - 1
		switch status {
		case 0, 2, 6, 5, 9, 7:
			common.Send(fmt.Sprintf("%s - %s 快递结束:%s", msg.Company, msg.Nu, expressageStatus[status]), GetWechatUrl("mine"))
			return
		case 8:
			unit = time.Minute * 10
		}
	end:
		time.Sleep(unit)
	}
}

var expressageStatus = map[int]string{
	0: "物流单暂无结果",
	1: "查询成功",
	2: "接口出现异常",
	3: "在途",
	4: "揽件",
	5: "疑难",
	6: "已签收",
	7: "退签",
	8: "派件中",
	9: "退回",
}

type expressageMsg struct {
	Success bool `json:"success"`
	//Ico         string `json:"ico"`         //
	//Phone       string `json:"phone"`       //
	//Url         string `json:"url"`         //
	Status int `json:"status"` //
	//Companytype string `json:"companytype"` //
	Nu      string `json:"nu"`      //
	Company string `json:"company"` //
	Reason  string `json:"reason"`  //
	Data    []struct {
		Time    string `json:"time"`    //时间
		Context string `json:"context"` //详情
	} `json:"data"`
	//Time     string `json:"time"`
	//Rank     string `json:"rank"`
	//Exceed   string `json:"exceed"`
	//Timeused string `json:"timeused"`
}

type companyInfo struct {
	//Name   string `json:"name"`
	Exname string `json:"exname"`
	//Ico    string `json:"ico"`
	//Url    string `json:"url"`
	//Phone  string `json:"phone"`
}

var companyTypes = map[string]string{
	"zt":  "zhongtong",
	"st":  "shentong",
	"bs":  "huitongkuaidi",
	"yt":  "yuantong",
	"yd":  "yunda",
	"tt":  "tiantian",
	"sf":  "shunfeng",
	"yz":  "youzhengguonei",
	"ems": "ems",
	"zjs": "zhaijisong",
	"qf":  "quanfengkuaidi",
	"fk":  "rufengda",
}
