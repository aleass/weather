package service

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"time"
	"weather/common"
)

var FuncDb *gorm.DB

func InitMysql() {
	if FuncDb != nil {
		return
	}
	config := MyConfig.DB
	dsn := fmt.Sprintf(
		"%s:%s@%s(%s:%s)/%s?charset=%s&multiStatements=true&parseTime=True&loc=Local",
		config.User,
		config.Password,
		"tcp",
		config.Host,
		config.Port,
		config.DbName,
		"utf8mb4",
	)
	var err error
	var _log = slowLog{}
	level := logger.Warn
	newLogger := logger.New(
		log.New(_log, "\r\n", log.LstdFlags),
		//log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Millisecond * 300, // 慢 SQL 阈值
			LogLevel:                  level,                  // 日志级别
			IgnoreRecordNotFoundError: true,                   // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  false,                  // 禁用彩色打印
		},
	)
	FuncDb, err = gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: newLogger})
	if err != nil {
		panic("mysql 启动失败!,原因:" + err.Error())
	}
}

type slowLog struct {
}

func (slowLog) Write(p []byte) (n int, err error) {
	common.Logger.Debug(string(p))
	return len(p), nil
}
