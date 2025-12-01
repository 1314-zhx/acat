package controller

import (
	"acat/dao"
	"acat/dao/db"
	"acat/logic"
	"acat/model"
	"acat/model/code"
	"acat/setting"
	"acat/util/auth"
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"log"
)

func AdminLoginHandler(c *gin.Context) {
	var adminLog logic.AdminLog
	var err error
	err = c.ShouldBindJSON(&adminLog)
	if err != nil {
		zap.L().Info("controller/admin.go AdminLoginHandler failed shouldBindJSON error : ", zap.Error(err))
		log.Println("controller/admin.go AdminLoginHandler failed shouldBindJSON error : ", err)
		c.JSON(400, ErrorResponse(errors.New("ShouldBindJSON error")))
		return
	}
	res := adminLog.Login()
	if res.Status != 200 {
		c.JSON(400, ErrorResponse(errors.New("管理员登录失败")))
		return
	}
	adminDao := dao.NewAdminDao(db.DB)
	var admin *model.AdminView
	admin, err = adminDao.GetAdmin(adminLog.Phone)
	if err != nil {
		zap.L().Info("controller/admin.go AdminLoginHandler failed GetAdmin error : ", zap.Error(err))
		log.Println("controller/admin.go AdminLoginHandler failed GetAdmin error : ", err)
		c.JSON(400, ErrorResponse(errors.New("GetAdmin error")))
		return
	}
	if res.Status == code.Success {
		token, err := auth.GenerateToken(admin.Aid, admin.Name, 1)
		if err != nil {
			c.JSON(500, ErrorResponse(errors.New("生成 token 失败")))
			return
		}
		secure := setting.Conf.WebMode == "release"
		// Token 是15分钟，cookie也要是15分钟，path 为"/admin"只有管理员路径可以用通用
		c.SetCookie("token", token, 15*60, "/admin", "", secure, true)
	}
	c.JSON(200, res)
}

func AdminSetScheduleHandler(c *gin.Context) {
	//TODO:
}
func AdminSetInterviewResultHandler(c *gin.Context) {
	//TODO:
}
func AdminPostEmailHandler(c *gin.Context) {
	//TODO:
}
