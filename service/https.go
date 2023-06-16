package service

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
	"runtime"
	"time"
	"weather/common"
)

var (
	ipUrls = [...]string{"https://qifu-api.baidubce.com/ip/geo/v1/district?ip=%s", "https://ipapi.co/%s/json/", "http://ip-api.com/json/%s?lang=zh-CN"}
)

type IPResponse struct {
	Latitude   float64 `json:"latitude" desc:"ipapi"`
	Longitude  float64 `json:"longitude" desc:"ipapi"`
	Latitude1  float64 `json:"lat" desc:"ip-api"`
	Longitude1 float64 `json:"lon" desc:"ip-api"`
}

type BaiduIpRes struct {
	Data struct {
		//Continent string `json:"continent"`
		//Country   string `json:"country"`
		//Zipcode   string `json:"zipcode"`
		//Timezone  string `json:"timezone"`
		//Accuracy  string `json:"accuracy"`
		//Owner     string `json:"owner"`
		//Isp       string `json:"isp"`
		//Source    string `json:"source"`
		//Areacode  string `json:"areacode"`
		//Adcode    string `json:"adcode"`
		//Asnumber  string `json:"asnumber"`
		Lat string `json:"lat"`
		Lng string `json:"lng"`
		//Radius    string `json:"radius"`
		//Prov      string `json:"prov"`
		//City      string `json:"city"`
		//District  string `json:"district"`
	} `json:"data"`
	//Charge bool `json:"charge"`
	Code string `json:"code"`
}

var (
	timeDump = make(chan struct{}, 1)
	localMap = map[string]*urlInfo{}
)

// Recover 错误
func Recover(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			stack := common.Stack(3)
			common.LogSend(fmt.Sprintf("api 发生panic：%v,stack:%s", err, string(stack)), common.PanicType)
			c.JSON(500, "服务器发生错误,请稍后再试")
			c.Abort()
		}
	}()
	c.Next()
}

func GetIpAddress(ip string) (string, error) {
	var (
		ipRes      = IPResponse{}
		ipBaiduRes = BaiduIpRes{}
	)
	for i, url := range ipUrls {
		url = fmt.Sprintf(url, ip)
		res, err := http.Get(url)
		if err != nil {
			return "", err
		}
		data, err := io.ReadAll(res.Body)
		if err != nil {
			return "", err
		}
		if i == 0 {
			if err = json.Unmarshal(data, &ipBaiduRes); err != nil {
				return "", err
			}
			//百度
			if ipBaiduRes.Code == "Success" && ipBaiduRes.Data.Lng != "" && ipBaiduRes.Data.Lat != "" {
				return fmt.Sprintf("%s,%s", ipBaiduRes.Data.Lng, ipBaiduRes.Data.Lat), nil
			}
			continue
		}

		if err = json.Unmarshal(data, &ipRes); err != nil {
			return "", err
		}

		//ipapi
		if ipRes.Longitude != 0 && ipRes.Latitude != 0 {
			return fmt.Sprintf("%f,%f", ipRes.Longitude, ipRes.Latitude), nil
		}
		//ip-api
		if ipRes.Latitude1 != 0 && ipRes.Longitude1 != 0 {
			return fmt.Sprintf("%f,%f", ipRes.Longitude1, ipRes.Latitude1), nil
		}
	}
	return "", nil
}

func HttpRun() {
	r := gin.Default()
	file, _ := os.Create("access.log")
	r.Use(gin.LoggerWithWriter(file, ""))
	r.Use(Recover)
	r.GET("/list", func(context *gin.Context) {
		var userList = []UrlInfo{}
		for _, info := range localMap {
			userList = append(userList, UrlInfo{
				Address: info.address,
				Name:    info.Name,
				IsRun:   info.IsRun,
			})
		}
		context.JSON(200, runtime.NumGoroutine())
		context.JSON(200, userList)
	})

	r.GET("/set", func(context *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				common.LogSend(fmt.Sprintf("panic err:%v", err), common.PanicType)
			}
		}()
		//判断是否存在
		name := context.Query("name")
		wechatNote := context.Query("note")
		op := context.Query("op") //del 清除
		if name == "" || wechatNote == "" {
			context.JSON(200, "name == nil")
			return
		}
		if allowUrlConfig[name] != wechatNote {
			context.JSON(200, name+" 非法")
			return
		}
		var (
			ip      = context.ClientIP()
			adcodes string
			err     error
		)
		info, ok := localMap[name]
		if op == "del" {
			//任务退出
			if !ok || !info.IsRun {
				context.JSON(200, name+" 不存在")
				return
			}
			//直接删除,否则会导致状态不一致
			delete(localMap, name)
			info.Switch <- struct{}{} //关闭一个任务
			goto end
		}
		op = "add"
		if ip == "" {
			context.JSON(200, "ip == nil")
			return
		}
		//ip判断
		adcodes, err = GetIpAddress(ip)
		if err != nil {
			context.JSON(200, ip+err.Error())
			return
		}
		if adcodes == "" {
			context.JSON(200, ip+":经纬度找不到")
			return
		}

		if !ok {
			info = getUrlInfo(name, adcodes, wechatNote)
			info.IsUrlConfig = true
			timeDump <- struct{}{}
			localMap[name] = info
			<-timeDump
			go info.WatchWeather()
		} else {
			op = "edit"
			localMap[name].CaiYunUrl = fmt.Sprintf(caiYunUrl, myConfig.CaiYun.Token, adcodes)
			if !info.IsRun { //启动
				go info.WatchWeather()
			}
		}
	end:
		_msg := fmt.Sprintf("%s %s 操作:%s ip:%s 坐标:%s", time.Now().Format("2006-01-02 15:04:05"),
			name, op, ip, adcodes)
		context.JSON(200, _msg)
		common.LogSend(_msg, common.InfoErrorType)
	})
	r.Run(":8080")
}
