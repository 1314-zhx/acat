/*
controller层用于接收前端数据
*/
package controller

import (
	"acat/logic"
	"acat/model/code"
	"acat/serializer"
	"acat/setting"
	"acat/util/auth"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"time"
)

func RegisterHandler(c *gin.Context) {
	var userRegister logic.UserRes
	err := c.ShouldBindJSON(&userRegister)
	if err != nil {
		zap.L().Info("controller/user.go RegisterHandler() failed shouldBindJSON() error : ", zap.Error(err))
		log.Println("controller/user.go RegisterHandler() failed shouldBindJSON() error : ", err)
		c.JSON(200, ErrorResponse(err))
		return
	}
	res := userRegister.Register()
	if res.Error != "" {
		c.JSON(400, ErrorResponse(errors.New(res.Error)))
		return
	}
	c.JSON(200, res)
}
func LoginHandler(c *gin.Context) {
	var user logic.UserLog
	ctx, cancel := context.WithTimeout(c.Request.Context(), 4*time.Second)
	defer cancel()
	err := c.ShouldBindJSON(&user)
	if err != nil {
		zap.L().Info("controller/user.go LoginHandler() failed shouldBindJSON() error : ", zap.Error(err))
		log.Println("controller/user.go LoginHandler() failed shouldBindJSON() error : ", err)
		c.JSON(200, ErrorResponse(err))
		return
	}

	res, userLogin := user.Login(ctx)
	if res.Error != "" {
		c.JSON(400, ErrorResponse(errors.New(res.Error)))
		return
	}
	if res.Status == code.Success {
		token, err := auth.GenerateToken(userLogin.UId, userLogin.Name, 0)
		if err != nil {
			c.JSON(500, serializer.Response{
				Status: 500,
				Msg:    "登录失败",
			})
			return
		}
		// 生产环境将secure设置为true，dev模式设置为false
		secure := setting.Conf.WebMode == "release"
		// Token 是7天，cookie也要是7天 path 为"/"全站通用
		c.SetCookie("token", token, 3600*24*7, "/user", "", secure, true)
	}
	c.JSON(http.StatusOK, res)
}

// 查询面试结果
func ResultHandler(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 4*time.Second)
	defer cancel()
	var userChe logic.UserChe
	var err error
	err = c.ShouldBindJSON(&userChe)
	if err != nil {
		zap.L().Info("controller/user.go ResultHandler() failed shouldBindJSON() error : ", zap.Error(err))
		log.Println("controller/user.go ResultHandler() failed shouldBindJSON() error : ", err)
		c.JSON(200, ErrorResponse(err))
		return
	}
	rawClaims, exists := c.Get("claims")
	if !exists {
		err := errors.New("认证信息缺失")
		zap.L().Warn("ResultHandler: claims 不存在", zap.Error(err))
		c.JSON(200, ErrorResponse(err))
		return
	}

	claims, ok := rawClaims.(*auth.JwtClaims)
	if !ok {
		err := errors.New("claims 类型错误")
		zap.L().Error("ResultHandler: claims 类型异常", zap.Error(err))
		c.JSON(200, ErrorResponse(err))
		return
	}
	res := userChe.Result(claims.UserID, ctx)
	if res.Error != "" {
		c.JSON(400, ErrorResponse(errors.New(res.Error)))
		return
	}
	c.JSON(200, res)
}
func PostHandler(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 4*time.Second)
	defer cancel()
	var userPost *logic.UserPost
	err := c.ShouldBindJSON(&userPost)
	if err != nil {
		zap.L().Info("controller/user.go PostHandler() failed shouldBindJSON() error : ", zap.Error(err))
		log.Println("controller/user.go PostHandler() failed shouldBindJSON() error : ", err)
		c.JSON(400, ErrorResponse(err))
		return
	}
	rawClaims, exists := c.Get("claims")
	if !exists {
		err := errors.New("认证信息缺失")
		zap.L().Warn("ResultHandler: claims 不存在", zap.Error(err))
		c.JSON(400, ErrorResponse(err))
		return
	}

	claims, ok := rawClaims.(*auth.JwtClaims)
	if !ok {
		err := errors.New("claims 类型错误")
		zap.L().Error("ResultHandler: claims 类型异常", zap.Error(err))
		c.JSON(400, ErrorResponse(err))
		return
	}
	res := userPost.Post(claims.UserID, ctx)
	if res.Error != "" {
		c.JSON(400, ErrorResponse(errors.New(res.Error)))
		return
	}
	c.JSON(200, res)
}
func LetterHandler(c *gin.Context) {
	var letter logic.Letter
	err := c.ShouldBindJSON(&letter)
	if err != nil {
		zap.L().Info("controller/user.go LetterHandler() failed shouldBindJSON() error : ", zap.Error(err))
		log.Println("controller/user.go LetterHandler() failed shouldBindJSON() error : ", err)
		c.JSON(400, ErrorResponse(err))
		return
	}
	rawClaims, exists := c.Get("claims")
	if !exists {
		err := errors.New("认证信息缺失")
		zap.L().Warn("ResultHandler: claims 不存在", zap.Error(err))
		c.JSON(400, ErrorResponse(err))
		return
	}

	claims, ok := rawClaims.(*auth.JwtClaims)
	if !ok {
		err := errors.New("claims 类型错误")
		zap.L().Error("ResultHandler: claims 类型异常", zap.Error(err))
		c.JSON(400, ErrorResponse(err))
		return
	}
	res := letter.Letter(claims.UserID)
	if res.Error != "" {
		c.JSON(400, ErrorResponse(errors.New(res.Error)))
		return
	}
	c.JSON(200, res)
}
func ShowAdminHandler(c *gin.Context) {
	res := logic.ShowAdmin()
	if res.Error != "" {
		c.JSON(400, ErrorResponse(errors.New(res.Error)))
		return
	}
	c.JSON(200, res)
}

func UpdateHandler(c *gin.Context) {
	var userUpdate logic.UserUpdate
	err := c.ShouldBindJSON(&userUpdate)
	fmt.Println("updateHandler", userUpdate)
	if err != nil {
		zap.L().Info("controller/user.go UpdateHandler() failed shouldBindJSON() error : ", zap.Error(err))
		log.Println("controller/user.go UpdateHandler() failed shouldBindJSON() error : ", err)
		c.JSON(400, ErrorResponse(err))
		return
	}
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
	res := userUpdate.Update(claims.UserID)
	if res.Error != "" {
		c.JSON(400, ErrorResponse(errors.New(res.Error)))
		return
	}
	c.JSON(200, res)
}
func ShowSlotHandler(c *gin.Context) {
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
	res := logic.ShowSlot(claims.UserID)
	if res.Error != "" {
		c.JSON(400, ErrorResponse(errors.New(res.Error)))
		return
	}
	c.JSON(200, res)
}
func ForgetHandler(c *gin.Context) {
	var forget logic.Forget
	err := c.ShouldBindJSON(&forget)
	if err != nil {
		zap.L().Info("controller/user.go ForgetHandler() failed shouldBindJSON() error : ", zap.Error(err))
		log.Println("controller/user.go ForgetHandler() failed shouldBindJSON() error : ", err)
		c.JSON(400, ErrorResponse(err))
		return
	}
	res := forget.Forget()
	if res.Error != "" {
		c.JSON(400, ErrorResponse(errors.New(res.Error)))
		return
	}
	c.JSON(200, res)
}
func ReSetPasswordHandler(c *gin.Context) {
	var reset logic.ResetPassword
	err := c.ShouldBindJSON(&reset)
	fmt.Println("reset ", reset)
	if err != nil {
		zap.L().Info("controller/user.go ReSetPasswordHandler() failed shouldBindJSON() error : ", zap.Error(err))
		log.Println("controller/user.go ReSetPasswordHandler() failed shouldBindJSON() error : ", err)
		c.JSON(400, ErrorResponse(err))
		return
	}
	res := reset.ResetPassword(c.Request.Context())
	if res.Error != "" {
		c.JSON(400, ErrorResponse(errors.New(res.Error)))
		return
	}
	c.JSON(200, res)
}
func LoginOutHandler(c *gin.Context) {
	secure := setting.Conf.WebMode == "release"
	// 后端发送过期cookie，将浏览器有效的cookie替换掉
	c.SetCookie("token", "", -1, "/user", "", secure, true)
	c.Redirect(http.StatusFound, "/user/center")
}
func CheckReplyHandler(c *gin.Context) {
	var checkReply logic.CheckReply
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
	res := checkReply.Reply(claims.UserID)
	if res.Error != "" {
		c.JSON(400, ErrorResponse(errors.New(res.Error)))
		return
	}
	c.JSON(200, res)
}
func DownloadHandler(c *gin.Context) {
	var dl logic.Download
	res := dl.DownloadQuestion("")
	date, _ := json.Marshal(res)
	fmt.Println("cs : ", string(date))
	if res.Error != "" {
		if res.Error == "文件不存在" || res.Error == "不支持的文件类型" {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "file not available"})
		} else {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": res.Error})
		}
		return
	}

	data, ok := res.Data.(map[string]string)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	filePath := data["file_path"]
	downloadName := data["filename"]

	c.Header("Cache-Control", "no-store")
	c.Header("Content-Disposition", "attachment; filename="+url.PathEscape(filepath.Base(downloadName)))
	c.Header("Content-Type", "application/octet-stream")
	c.File(filePath)

}

// 渲染
func ShowLoginHandler(c *gin.Context) {
	c.HTML(200, "login.html", nil)
}
