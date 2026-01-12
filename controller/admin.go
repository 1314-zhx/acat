package controller

import (
	"acat/dao"
	"acat/dao/db"
	"acat/logic"
	"acat/model"
	"acat/model/code"
	"acat/setting"
	"acat/util/auth"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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
	var setSlot logic.SetSlot
	err := c.ShouldBindJSON(&setSlot)
	if err != nil {
		zap.L().Info("controller/admin.go AdminSetScheduleHandler failed shouldBindJSON error : ", zap.Error(err))
		log.Println("controller/admin.go AdminSetScheduleHandler failed shouldBindJSON error : ", err)
		c.JSON(400, ErrorResponse(errors.New("ShouldBindJSON error")))
		return
	}
	res := setSlot.SetSchedule()
	if res.Error != "" {
		c.JSON(400, ErrorResponse(errors.New(res.Error)))
		return
	}
	c.JSON(http.StatusOK, res)
}
func AdminSetInterviewResultHandler(c *gin.Context) {
	var setResult logic.SetResult
	err := c.ShouldBindJSON(&setResult)
	if err != nil {
		zap.L().Info("controller/admin.go AdminSetInterviewResultHandler failed shouldBindJSON error : ", zap.Error(err))
		log.Println("controller/admin.go AdminSetInterviewResultHandler failed shouldBindJSON error : ", err)
		c.JSON(400, ErrorResponse(errors.New("ShouldBindJSON error")))
		return
	}
	res := setResult.SetUserResult(c.Request.Context())
	if res.Error != "" {
		c.JSON(400, ErrorResponse(errors.New(res.Error)))
		return
	}
	c.JSON(200, res)
}
func SetPassHandler(c *gin.Context) {
	var setPass logic.SetPass
	err := c.ShouldBindJSON(&setPass)
	rawClaims, exists := c.Get("claims")
	if !exists {
		err := errors.New("认证信息缺失")
		zap.L().Warn("ResultHandler: claims 不存在", zap.Error(err))
		c.JSON(400, ErrorResponse(err))
		return
	}
	// 类型断言
	claims, ok := rawClaims.(*auth.JwtClaims)
	if !ok {
		err := errors.New("claims 类型错误")
		zap.L().Error("ResultHandler: claims 类型异常", zap.Error(err))
		c.JSON(400, ErrorResponse(err))
		return
	}
	if err != nil {
		zap.L().Info("controller/admin.go SetPassHandler failed shouldBindJSON error : ", zap.Error(err))
		log.Println("controller/admin.go SetPassHandler failed shouldBindJSON error : ", err)
		c.JSON(400, ErrorResponse(errors.New("ShouldBindJSON error")))
		return
	}
	res := setPass.SetUserPass(claims.UserID)
	if res.Error != "" {
		c.JSON(400, ErrorResponse(errors.New(res.Error)))
		return
	}
	c.JSON(200, res)
}

// 管理员公布面试结果
func AdminPublishHandler(c *gin.Context) {
	var passUser logic.PassUser
	err := c.ShouldBindJSON(&passUser)
	if err != nil {
		zap.L().Info("logic/admin.go AdminPublishHandler failed shouldBindJSON error : ", zap.Error(err))
		log.Println("logic/admin.go AdminPublishHandler failed shouldBindJSON error : ", err)
		c.JSON(400, ErrorResponse(errors.New("ShouldBindJSON error")))
		return
	}
	res := passUser.Publish()
	if res.Error != "" {
		c.JSON(400, ErrorResponse(errors.New(res.Error)))
		return
	}
	c.JSON(200, res)
}
func GetPassUserHandler(c *gin.Context) {
	var round logic.PublishRound
	err := c.ShouldBindJSON(&round)
	if err != nil {
		zap.L().Info("logic/admin.go GetPassUserHandler failed shouldBindJSON error : ", zap.Error(err))
		log.Println("logic/admin.go GetPassUserHandler failed shouldBindJSON error : ", err)
		c.JSON(400, ErrorResponse(errors.New("ShouldBindJSON error")))
		return
	}
	res := round.GetPassUser()
	if res.Error != "" {
		c.JSON(400, ErrorResponse(errors.New(res.Error)))
		return
	}
	c.JSON(200, res)
}

func AdminLetterHandler(c *gin.Context) {
	var adminLetter logic.AdminLetter
	rawClaims, exists := c.Get("claims")
	if !exists {
		err := errors.New("认证信息缺失")
		zap.L().Warn("ResultHandler: claims 不存在", zap.Error(err))
		c.JSON(400, ErrorResponse(err))
		return
	}
	// 类型断言
	claims, ok := rawClaims.(*auth.JwtClaims)
	if !ok {
		err := errors.New("claims 类型错误")
		zap.L().Error("ResultHandler: claims 类型异常", zap.Error(err))
		c.JSON(400, ErrorResponse(err))
		return
	}
	res := adminLetter.Letter(claims.UserID)
	if res.Error != "" {
		c.JSON(400, ErrorResponse(errors.New(res.Error)))
		return
	}
	data, _ := json.Marshal(res)
	fmt.Println(string(data))
	c.JSON(200, res)
}
func AdminReplyHandler(c *gin.Context) {
	var reply logic.AdminReply
	err := c.ShouldBindJSON(&reply)
	if err != nil {
		zap.L().Info("logic/admin.go AdminReplyHandler failed shouldBindJSON error : ", zap.Error(err))
		log.Println("logic/admin.go AdminReplyHandler failed shouldBindJSON error : ", err)
		c.JSON(400, ErrorResponse(errors.New("ShouldBindJSON error")))
		return
	}
	res := reply.Reply()
	if res.Error != "" {
		c.JSON(400, ErrorResponse(errors.New(res.Error)))
		return
	}
	c.JSON(200, res)
}

func UploadHandler(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "未提供文件"})
		return
	}
	var upload logic.Upload
	res := upload.UploadQuestion(file)
	if res.Error != "" {
		c.JSON(400, ErrorResponse(errors.New(res.Error)))
		return
	}
	c.JSON(http.StatusOK, res)
}
