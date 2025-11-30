/*
logic层用于数据校验
*/
package logic

import (
	"acat/dao"
	"acat/dao/db"
	"acat/model"
	"acat/model/code"
	"acat/redislock"
	"acat/serializer"
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"log"
	"regexp"
)

//---------------->用户登录<-------------------//

// 用于接收前端登录数据
type UserLog struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
	VCode    string `json:"v_code"`
}

func (l *UserLog) Login(ctx context.Context) (serializer.Response, *model.UserLogin) {
	co := code.Success
	// 参数校验
	if l.Phone == "" || l.Password == "" {
		co = code.MissMustInfo
		return serializer.Response{
			Status: co,
			Msg:    code.GetMsg(co),
		}, nil
	}
	if len(l.Phone) != 11 {
		co = code.InvalidPhoneForm
		return serializer.Response{
			Status: co,
			Msg:    code.GetMsg(co),
		}, nil
	}

	userDao := dao.NewUserDao(db.DB)
	// 1 先查用户（只用手机号）
	userModel, err := userDao.GetUserByPhone(l.Phone, ctx)
	if err != nil {
		co = code.PhoneOrPasswordError
		zap.L().Info("用户不存在或数据库查询失败", zap.String("phone", l.Phone), zap.Error(err))
		return serializer.Response{
			Status: co,
			Msg:    "手机号或密码错误",
		}, nil
	}
	// 2 比对密码（用查到的 userModel.Password）
	if err := bcrypt.CompareHashAndPassword([]byte(userModel.Password), []byte(l.Password)); err != nil {
		co = code.PhoneOrPasswordError
		zap.L().Info("密码错误", zap.String("phone", l.Phone))
		return serializer.Response{
			Status: co,
			Msg:    "手机号或密码错误",
		}, nil
	}

	userLogin := &model.UserLogin{
		UId:   userModel.ID,
		Name:  userModel.Name,
		Phone: userModel.Phone,
	}
	fmt.Println("logic验证成功")
	return serializer.Response{
		Status: co,
		Msg:    code.GetMsg(co),
	}, userLogin
}

//---------------->用户注册<-------------------//

// 用于接收前端注册数据
type UserRes struct {
	Name       string `json:"name"`
	Phone      string `json:"phone"`
	Password   string `json:"password"`
	RePassword string `json:"re_password"`
	Email      string `json:"email"` // 非必填
	StuId      string `json:"stu_id"`
	Gender     int    `json:"gender"`    // 1为男2为女
	Direction  int    `json:"direction"` // 方向选择，0不确定，1为Go，2为Java，3为前端，4为后端
}

func (r *UserRes) Register() serializer.Response {
	co := code.Success
	// 检验必填参数
	if r.Phone == "" || r.Password == "" || r.RePassword == "" || r.StuId == "" || r.Name == "" || r.Gender == 0 {
		co = code.MissMustInfo
		zap.L().Info("logic/user.go Register() failed miss must info : ", zap.Error(errors.New(code.GetMsg(co))))
		log.Println("logic/user.go Register() failed miss must info : ", errors.New(code.GetMsg(co)))
		return serializer.Response{
			Status: co,
			Msg:    code.GetMsg(co),
			Error:  errors.New("缺少必填字段").Error(),
		}
	}
	// 检验参数格式
	if r.Password != r.RePassword {
		co = code.PasswordUnequalRePassword
		zap.L().Info("logic/user.go Register() failed  : ", zap.Error(errors.New(code.GetMsg(co))))
		log.Println("logic/user.go Register() failed  : ", errors.New(code.GetMsg(co)))
		return serializer.Response{
			Status: co,
			Msg:    code.GetMsg(co),
			Error:  errors.New("参数错误").Error(),
		}
	}
	if len(r.Phone) != 11 {
		co = code.InvalidPhoneForm
		zap.L().Info("logic/user.go Register() failed phone form error : ", zap.Error(errors.New(code.GetMsg(co))))
		log.Println("logic/user.go Register() failed phone form error : ", errors.New(code.GetMsg(co)))
		return serializer.Response{
			Status: co,
			Msg:    code.GetMsg(co),
			Error:  errors.New("格式不正确").Error(),
		}
	}
	validEmail := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !validEmail.MatchString(r.Email) {
		co = code.EmailFormError
		zap.L().Info("logic/user.go Register() failed email form error : ", zap.Error(errors.New(code.GetMsg(co))))
		log.Println("logic/user.go Register() failed email form error : ", errors.New(code.GetMsg(co)))
		return serializer.Response{
			Status: co,
			Msg:    code.GetMsg(co),
			Error:  errors.New("格式不正确").Error(),
		}
	}
	// 传给dao进行存储
	userDao := dao.NewUserDao(db.DB)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(r.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("bcrypt.GenerateFromPassword error", err)
	}
	userRegister := &model.UserRegister{
		Name:      r.Name,
		StuId:     r.StuId,
		Password:  string(hashedPassword),
		Phone:     r.Phone,
		Email:     r.Email,
		Gender:    r.Gender,
		Direction: r.Direction,
	}
	err = userDao.Create(userRegister)
	return serializer.Response{
		Status: co,
		Msg:    code.GetMsg(co),
	}
}

//----------->用户查询面试结果<--------------//

type UserChe struct {
	Round int `json:"round"`
}

func (ch *UserChe) Result(uid uint, ctx context.Context) serializer.Response {
	co := code.Success
	userDao := dao.NewUserDao(db.DB)
	pass, err := userDao.Result(uid, ch.Round, ctx)
	if pass == 0 && err == nil {
		co = code.Error
		return serializer.Response{
			Status: co,
			Data:   "学生未参加面试",
			Msg:    code.GetMsg(co),
			Error:  "学生未参加面试",
		}
	}
	if err != nil && pass == 0 {
		co = code.Error
		return serializer.Response{
			Status: co,
			Data:   "数据库错误",
			Msg:    code.GetMsg(co),
			Error:  "数据库错误",
		}
	}
	if pass == 2 {
		return serializer.Response{
			Status: co,
			Data:   "学生未通过面试",
			Msg:    code.GetMsg(co),
		}
	}

	return serializer.Response{
		Status: co,
		Data:   "学生通过面试",
		Msg:    code.GetMsg(co),
	}
}

//------------>用户报名<--------------//

type UserPost struct {
	Name      string `json:"name"`
	Direction int    `json:"direction"`
	SlotId    uint   `json:"slot_id"` // 由前端返回
}

func (p *UserPost) Post(uid uint, ctx context.Context) serializer.Response {
	co := code.Success
	var err error
	var user *model.UserModel
	userDao := dao.NewUserDao(db.DB)
	user, err = userDao.GetUserByID(uid, ctx)
	if err != nil {
		zap.L().Info("logic/user.go Post get user by id error : ", zap.Error(err))
		log.Println("logic/user.go Post get user by id error : ", err)
	}
	slotDao := dao.NewSlotDao(db.DB)
	var slot *model.InterviewSlot
	slot, err = slotDao.GetSlotById(p.SlotId)
	if err != nil {
		zap.L().Info("logic/user.go Post get slot by id error : ", zap.Error(err))
		log.Println("logic/user.go Post get slot by id error : ", err)
	}
	// 用户报二面，但一面没过
	if slot.Round == 2 && user.FirstPass != 1 {
		co = code.FirstViewNotPass
		return serializer.Response{
			Status: co,
			Msg:    code.GetMsg(co),
			Error:  "一面未通过，无法参加二面",
		}
	}
	userDaoWithRdb := dao.NewUserDaoWithRdb(db.DB, redislock.GetRDB())
	err = userDaoWithRdb.Post(p.SlotId, user.ID, ctx)
	if err != nil {
		zap.L().Info("logic/user.go failed post error : ", zap.Error(err))
		log.Println("logic/user.go failed post error : ", err)
	}
	return serializer.Response{
		Status: co,
		Msg:    code.GetMsg(co),
	}
}
