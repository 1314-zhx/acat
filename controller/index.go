package controller

import (
	"acat/logic"
	"acat/model/code"
	"acat/serializer"
	"context"
	"github.com/gin-gonic/gin"
	"time"
)

func IndexHandler(c *gin.Context) {
	//TODO:
	c.HTML(200, "index.html", nil)
}
func CheckHandler(c *gin.Context) {
	co := code.Success
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()
	schedule, err := logic.Check(ctx)
	if err != nil {
		co = code.Error
		c.JSON(400, serializer.Response{
			Status: co,
			Data:   nil,
			Msg:    code.GetMsg(co),
			Error:  err.Error(),
		})
		return
	}
	c.JSON(200, serializer.Response{
		Status: co,
		Data:   schedule,
		Msg:    code.GetMsg(co),
		Error:  "",
	})
}
func LetterHandler(c *gin.Context) {
	//TODO:
}
func Norouter(c *gin.Context) {
	c.HTML(200, "norouter.html", nil)
}
