/*
在main文件中初始化组件
*/
package main

import (
	"acat/dao/db"
	"acat/logger"
	"acat/redislock"
	"acat/router"
	"acat/setting"
	"fmt"
	"go.uber.org/zap"
	"log"
)

func main() {
	r := router.NewRouter()
	// 初始化配置文件，很重要，没有无法运行，使用失败必须退出
	if err := setting.Init(); err != nil {
		log.Fatal("[PANIC] in main setting.Init() failed , error: ", err)
	}
	// 初始化日志文件，失败与否对程序影响较小，可以强行启动，无须退出
	if err := logger.Init(setting.Conf.LogConf, "dev"); err != nil {
		log.Println("[INFO] log file db failed && WHERE main.go ,error : ", err)
	}
	zap.L().Info("日志文件初始化完成")
	// 初始化mysql数据库
	if err := db.InitDB(); err != nil {
		log.Fatal("[PANIC] in main db.InitDB() failed , error: ", err)
	}
	// 初始化全局唯一redis实例，默认没有错误
	redislock.Init(setting.Conf.RedisHost+":"+setting.Conf.RedisPort, setting.Conf.RedisDb, setting.Conf.RedisPoolSize)
	zap.L().Info("数据库表同步完成")
	_ = r.Run(fmt.Sprintf(":%s", setting.Conf.WebPort))
}
