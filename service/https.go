package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
	"time"
	"weather/common"
)

var (
	timeDump = make(chan struct{}, 1)
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

func HttpRun() {
	r := gin.Default()
	file, _ := os.Create("access.log")
	r.Use(gin.LoggerWithWriter(file, ""))
	r.Use(Recover)
	r.GET("/list", ListConfigUser)
	r.GET("/set", UserHandler)

	r.Run(":8080")
}

func ListConfigUser(context *gin.Context) {
	//var userList = []BaseWeatherInfo{}
	//for _, info := range taskMap {
	//for _,data := range info.ConfigGroup {
	//	userList = append(userList, *data)
	//}
	context.JSON(200, taskMap)
	//}
}

//
//func GetGroupWeatherData(context *gin.Context) {
//	defer func() {
//		if err := recover(); err != nil {
//			common.LogSend(fmt.Sprintf("panic err:%v", err), common.PanicType)
//		}
//	}()
//	//判断是否存在
//	name := context.Query("name")
//	wechatNote := context.Query("note")
//	if name == "" || wechatNote == "" {
//		context.JSON(200, "name == nil")
//		return
//	}
//	if allowUrlConfig[name] != wechatNote {
//		context.JSON(200, name+" 非法")
//		return
//	}
//	info, ok := taskMap[name]
//
//	if !ok {
//		context.JSON(200, name+" 不存在")
//		return
//	}
//
//}

func UserHandler(context *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			common.LogSend(fmt.Sprintf("panic err:%v", err), common.PanicType)
		}
	}()
	//判断是否存在
	name := context.Query("name")
	wechatNote := context.Query("note")
	key := context.Query("key")
	unit := context.Query("unit")
	op := context.Query("op") //del 清除
	if name == "" || wechatNote == "" || key == "" {
		context.JSON(200, "name == nil")
		return
	}
	if allowUrlConfig[name] != wechatNote {
		context.JSON(200, name+" 非法")
		return
	}

	info, ok := taskMap[name]
	if !ok {
		context.JSON(200, name+" 不存在")
		return
	}

	var (
		ip, addr, addrCodes, main string
		err                       error
		data                      *BaseWeatherInfo
	)
	//任务退出
	if op == "del" {
		weather, ok := info.ConfigGroup[key]
		if !ok {
			context.JSON(200, key+" 不存在")
			return
		}
		weather.Switch = false
		goto end
	}

	op = "add"
	switch unit {
	case "ip":
		ip = context.ClientIP()
	case "addr":
		addr = key
	case "addr_codes":
		addrCodes = key
	}

	//地址搜索
	if addr != "" {
		addrCodes, addr, err = common.GetKeyWordAddr(addr)
		if err != nil {
			context.JSON(200, "ip or addrCodes == nil")
			return
		}
		if addrCodes != "" {
			main = addr
			goto start
		}
	}

	//经纬度
	if addrCodes != "" {
		main = addrCodes
		goto start
	}

	//ip判断
	if ip == "" {
		context.JSON(200, "ip nil")
		return
	}

	addrCodes, err = common.GetIpAddress(ip)
	if err != nil {
		context.JSON(200, ip+err.Error())
		return
	}
	if addrCodes == "" {
		context.JSON(200, ip+":经纬度找不到")
		return
	}

	main = ip
start:
	data, ok = info.ConfigGroup[key]
	if !ok {
		weather := getUrlInfo(name, addrCodes, wechatNote, "", 1, true)
		weather.IsUrlConfig = true
		timeDump <- struct{}{}
		taskMap[name] = info
		<-timeDump
		go info.WatchWeather()
	} else {
		op = "edit"
		data.CaiYunUrl = fmt.Sprintf(caiYunUrl, myConfig.CaiYun.Token, addrCodes)
	}
	data.Ip = ip
	data.Op = op
	data.AddrCodes = addrCodes
	data.Main = main
	data.AllowNight = context.Query("night") == "true"

end:
	_msg := fmt.Sprintf("%s %s 操作:%s-%s ip:%s 坐标:%s ", time.Now().Format("2006-01-02 15:04:05"),
		name, op, main, ip, addrCodes)
	context.JSON(200, _msg)
	common.LogSend(_msg, common.InfoErrorType)
}
