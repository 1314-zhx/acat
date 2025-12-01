package db

import (
	"acat/model"
	"acat/setting"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"sync"
)

var DB *gorm.DB
var once sync.Once

// InitDB 初始化 GORM 数据库连接
func InitDB() error {
	var err error
	once.Do(func() {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			setting.Conf.MySQLUser,
			setting.Conf.MySQLPwd,
			setting.Conf.MySQLHost,
			setting.Conf.MySQLPort,
			setting.Conf.MySQLDbName)

		DB, err = gorm.Open(mysql.Open(dsn))
		if err != nil {
			log.Println("GORM 初始化失败 error=", err)
			zap.L().Error("GORM 初始化失败", zap.Error(err))
			return
		}

		err = autoMigrate()
		if err != nil {
			log.Println("自动迁移表结构失败 error=", err)
			zap.L().Error("自动迁移表结构失败", zap.Error(err))
		}
	})
	return err
}

// autoMigrate 自动迁移表结构
func autoMigrate() error {
	return DB.AutoMigrate(
		&model.UserModel{},
		&model.AdminModel{},
		&model.InterviewSlot{},
		&model.InterviewAssignment{},
		&model.InterviewResult{},
		&model.Message{},
	)
}
