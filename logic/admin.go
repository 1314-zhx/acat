package logic

import (
	"acat/dao"
	"acat/dao/db"
	"acat/model/code"
	"acat/serializer"
	"errors"
	"go.uber.org/zap"
	"log"
)

type AdminLog struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

func (a *AdminLog) Login() serializer.Response {
	co := code.Success
	adminDao := dao.NewAdminDao(db.DB)
	is, err := adminDao.Login(a.Phone, a.Password)
	if err != nil {
		co = code.Error
		zap.L().Info("logic/admin.go Login failed error : ", zap.Error(err))
		log.Println("logic/admin.go Login failed error : ", err)
		return serializer.Response{
			Status: co,
			Msg:    code.GetMsg(co),
			Error:  errors.New("数据库错误").Error(),
		}
	}
	if is == false {
		co = code.Error
		return serializer.Response{
			Status: code.Error,
			Msg:    code.GetMsg(co),
			Error:  "",
		}
	}
	return serializer.Response{
		Status: co,
		Msg:    code.GetMsg(co),
	}
}
