package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
	"runtime"
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
	var userList = []UrlInfo{}
	for _, info := range taskMap {
		userList = append(userList, UrlInfo{
			Address:     info.address,
			Name:        info.Name,
			IsRun:       info.IsRun,
			IsUrlConfig: info.IsUrlConfig,
		})
	}
	context.JSON(200, runtime.NumGoroutine())
	context.JSON(200, userList)

}
func UserHandler(context *gin.Context) {
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
		ip                  = context.ClientIP()
		adcodes, addr, main string
		err                 error
	)
	ip = "43.138.231.123"
	info, ok := taskMap[name]
	if op == "del" {
		//任务退出
		if !ok || !info.IsRun {
			context.JSON(200, name+" 不存在")
			return
		}
		//直接删除,否则会导致状态不一致
		delete(taskMap, name)
		info.Switch <- struct{}{} //关闭一个任务
		goto end
	}
	op = "add"
	//地址搜索
	addr = context.Query("addr")
	if addr != "" {
		adcodes, addr, err = common.GetKeyWordAddr(addr)
		if err != nil {
			context.JSON(200, "ip or adcodes == nil")
			return
		}
		if adcodes != "" {
			main = addr
			goto start
		}
	}
	adcodes = context.Query("adcodes")

	//优先经纬度
	if adcodes != "" {
		main = adcodes
		goto start
	}

	//ip判断
	if ip == "" {
		context.JSON(200, "ip nil")
		return
	}

	adcodes, err = common.GetIpAddress(ip)
	if err != nil {
		context.JSON(200, ip+err.Error())
		return
	}
	if adcodes == "" {
		context.JSON(200, ip+":经纬度找不到")
		return
	}

	main = ip
start:
	if !ok {
		info = getUrlInfo(name, adcodes, wechatNote, "", 1)
		info.IsUrlConfig = true
		timeDump <- struct{}{}
		taskMap[name] = info
		<-timeDump
		go info.WatchWeather()
	} else {
		op = "edit"
		taskMap[name].CaiYunUrl = fmt.Sprintf(caiYunUrl, myConfig.CaiYun.Token, adcodes)
		if !info.IsRun { //启动
			go info.WatchWeather()
		}
	}
	info.Ip = ip
	info.Op = op
	info.Adcodes = adcodes
	info.Main = main
	info.AllowNight = context.Query("night") == "true"
end:
	_msg := fmt.Sprintf("%s %s 操作:%s-%s ip:%s 坐标:%s ", time.Now().Format("2006-01-02 15:04:05"),
		name, op, main, ip, adcodes)
	context.JSON(200, _msg)
	common.LogSend(_msg, common.InfoErrorType)
}
