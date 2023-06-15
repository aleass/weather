package service

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
	"time"
	"weather/common"
)

var (
	ipUrls = [2]string{"https://ipapi.co/%s/json/", "http://ip-api.com/json/%s?lang=zh-CN"}
)

type IPResponse struct {
	Latitude   float64 `json:"latitude" desc:"ipapi"`
	Longitude  float64 `json:"longitude" desc:"ipapi"`
	Latitude1  float64 `json:"lat" desc:"ip-api"`
	Longitude2 float64 `json:"lon" desc:"ip-api"`
}

var (
	timeDump = make(chan struct{}, 1)
	localMap = map[string]*UrlInfo{}
	ipRes    = IPResponse{}
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

func GetIpAddress(ip string) (float64, float64, error) {
	for i := 0; i < 2; i++ {
		var url = fmt.Sprintf(ipUrls[i], ip)
		res, err := http.Get(url)
		if err != nil {
			return 0, 0, err
		}
		data, err := io.ReadAll(res.Body)
		if err != nil {
			return 0, 0, err
		}
		if err = json.Unmarshal(data, &ipRes); err != nil {
			return 0, 0, err
		}
		if ipRes.Longitude != 0 && ipRes.Latitude != 0 {
			return ipRes.Longitude, ipRes.Latitude, nil
		}
		if ipRes.Latitude1 != 0 && ipRes.Longitude2 != 0 {
			return ipRes.Longitude2, ipRes.Latitude1, nil
		}
	}
	return 0, 0, nil
}

func HttpRun() {
	r := gin.Default()
	file, _ := os.Create("access.log")
	r.Use(gin.LoggerWithWriter(file, ""))
	r.Use(Recover)
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
			ip                  string
			Latitude, Longitude float64
			err                 error
		)
		info, ok := localMap[name]
		if op == "del" {
			//任务退出
			if !ok || !info.isRun {
				context.JSON(200, name+" 不存在")
				return
			}
			info._switch <- struct{}{} //关闭一个任务
			info.isRun = false
			goto end
		}
		op = "add"

		//ip判断
		ip = context.ClientIP()
		if ip == "" {
			context.JSON(200, "ip == nil")
			return
		}
		//获取ip信息

		Longitude, Latitude, err = GetIpAddress(ip)
		if err != nil {
			context.JSON(200, err.Error())
			return
		}
		if Latitude == 0 || Longitude == 0 {
			context.JSON(200, "经纬度找不到")
			return
		}

		if !ok {
			info = getUrlInfo(name, fmt.Sprintf("%f,%f", Longitude, Latitude), wechatNote)
			info.isRun = true
			info.isUrlConfig = true
			timeDump <- struct{}{}
			localMap[name] = info
			<-timeDump
			go info.WatchWeather()
		} else {
			op = "edit"
			localMap[name].caiYunUrl = fmt.Sprintf(caiYunUrl, myConfig.CaiYun.Token, fmt.Sprintf("%f,%f", Longitude, Latitude))
		}
	end:
		_msg := fmt.Sprintf("%s %s %s ip:%s,经度:%f,纬度:%f", time.Now().Format("2006-01-02 15:04:05"),
			name, op, ip, Longitude, Latitude)
		context.JSON(200, _msg)
		common.LogSend(_msg, common.InfoErrorType)
	})
	r.Run(":8080")
}
