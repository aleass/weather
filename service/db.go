package service

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"weather/common"
)

var funcDb *gorm.DB

func getMysql() {
	if funcDb != nil {
		return
	}
	config := myConfig.DB
	dsn := fmt.Sprintf(
		"%s:%s@%s(%s)/%s?charset=%s&multiStatements=true&parseTime=True&loc=Local",
		config.User,
		config.Password,
		"tcp",
		config.Host,
		config.DbName,
		"utf8mb4",
	)

	var (
		err error
	)

	funcDb, err = gorm.Open(mysql.Open(dsn))
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
