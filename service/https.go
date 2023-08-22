package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"os"
	"runtime"
	"strings"
	"time"
	"weather/common"
)

var (
	timeDump = make(chan struct{}, 1)
	healUrl  string
)

// Recover 错误
func Recover(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			stack := common.Stack(2)
			common.LogSend(fmt.Sprintf("api 发生panic：%v,stack:%s", err, string(stack)), common.PanicType)
			c.JSON(500, "服务器发生错误,请稍后再试")
			c.Abort()
		}
	}()
	c.Next()
}

type HeartRate struct {
	Heart string `json:"heart"`
	Time  string `json:"time"`
}

func HttpRun() {
	r := gin.Default()
	file, _ := os.Create("access.log")
	r.Use(gin.LoggerWithWriter(file, ""))
	r.Use(Recover)
	r.GET("/list", ListConfigUser)
	r.GET("/set", UserHandler)

	//健康
	h := r.Group("/health")
	{
		h.POST("/heart_rate", heartRate)
	}

	//fund
	f := r.Group("fund")
	{
		f.GET("day", GetTodayEFund)
	}

	//other
	r.GET("kd", GetExpressage)

	r.Run(":8080")
}

var last int64

func heartRate(c *gin.Context) {
	var h HeartRate
	c.Bind(&h)
	if h.Time == "" {
		return
	}
	heartsStr := strings.Split(h.Heart, "\n")
	timeStr := strings.Split(h.Time, "\n")
	now := time.Now()
	heartMsg := ""
	var _last int64
	for i, heart := range heartsStr {
		//保留间隔5分钟以上的
		heartTime, _ := time.ParseInLocation("2006年1月2日 15:04", timeStr[i], time.Local)
		if heartTime.Unix() <= last || _last-heartTime.Unix() < 300 && _last != 0 {
			continue
		}
		_last = heartTime.Unix()
		if index := strings.Index(heart, "."); index != -1 {
			heart = heart[:index]
		}
		heartMsg += "\n" + heartTime.Format(common.UsualTimeDay) + " " + heart
	}
	heartTime, _ := time.ParseInLocation("2006年1月2日 15:04", timeStr[0], time.Local)
	if heartMsg == "" {
		return
	}
	//设置curr最新的数据
	last = heartTime.Unix()
	common.Send(now.Format(common.UsualTimeHour)+heartMsg, healUrl)
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
		info.isEdit = true //发生了变化发送一次
		info.CaiYunUrl = fmt.Sprintf(caiYunUrl, MyConfig.CaiYun.Token, adcodes)
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

func GetTodayEFund(c *gin.Context) {
	var list []common.DaysPastTimeRank
	FuncDb.Raw(common.DaysPastTimeAverSql).Find(&list)
	buff := strings.NewReader(common.AdjustData(list))
	io.Copy(c.Writer, buff)

	FuncDb.Raw(common.DaysPastTimeRankSql).Find(&list)
	buff = strings.NewReader(common.AdjustData(list))
	io.Copy(c.Writer, buff)
}
