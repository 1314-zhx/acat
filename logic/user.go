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
	"gorm.io/gorm"
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
	fmt.Println(user.FirstPass, user)
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
		co = code.Error
		zap.L().Info("logic/user.go failed post error : ", zap.Error(err))
		log.Println("logic/user.go failed post error : ", err)
	}
	return serializer.Response{
		Status: co,
		Msg:    code.GetMsg(co),
	}
}

//-------------->用户给管理员发消息<-----------------//
type Letter struct {
	ReceiveID uint   `json:"receive_id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
}

func (l *Letter) Letter(uid uint) serializer.Response {
	co := code.Success
	if l.Title == "" {
		co = code.Error
		return serializer.Response{
			Status: co,
			Msg:    code.GetMsg(co),
			Error:  errors.New("无标题").Error(),
		}
	}
	if len(l.Content) > 50 {
		co = code.Error
		return serializer.Response{
			Status: co,
			Msg:    code.GetMsg(co),
			Error:  errors.New("正文超过50字").Error(),
		}
	}
	letterDao := dao.NewLetterDao(db.DB)
	err := letterDao.Letter(l.Title, l.Content, l.ReceiveID, uid)
	if err != nil {
		zap.L().Info("logic/letter.go failed dao error : ", zap.Error(err))
		log.Println("logic/letter.go failed dao error : ", err)
		co = code.Error
		return serializer.Response{
			Status: co,
			Msg:    code.GetMsg(co),
			Error:  errors.New("数据库错误").Error(),
		}
	}
	return serializer.Response{
		Status: co,
		Msg:    code.GetMsg(co),
	}
}

// 展示管理员列表
func ShowAdmin() serializer.Response {
	co := code.Success
	// 因为都是用户相关操作，使用都使用用户的数据库实例模型
	showAdminDao := dao.NewUserDao(db.DB)
	admins, err := showAdminDao.GetAdmins()
	if err != nil {
		co = code.Error
		return serializer.Response{
			Status: co,
			Data:   nil, // 失败，不返回任何值
			Msg:    code.GetMsg(co),
			Error:  errors.New("数据库错误").Error(),
		}
	}
	// 脱敏
	var adminViews []model.AdminView
	for _, admin := range admins {
		var adminView model.AdminView
		adminView.Name = admin.Name
		adminView.Phone = admin.Phone
		adminView.Direction = admin.Direction
		adminView.Aid = admin.ID
		adminViews = append(adminViews, adminView)
	}
	return serializer.Response{
		Status: co,
		Data:   adminViews,
		Msg:    code.GetMsg(co),
	}
}

//--------------->用户更新<-----------------//
type UserUpdate struct {
	IsDelete int `json:"is_delete"` // 1 表示删除该面试预约 0 表示不删除
	UserPost                        // 继承UserPost，其实更新也相当于一次新的报名
}

func (u *UserUpdate) Update(uid uint) serializer.Response {
	co := code.Success
	if u.IsDelete == 1 {
		userDao := dao.NewUserDao(db.DB)
		err := userDao.Delete(uid)
		if err != nil {
			co = code.Error
			zap.L().Info("logic/user.go Update failed error : ", zap.Error(err))
			log.Println("logic/user.go Update failed error : ", err)
			return serializer.Response{
				Status: co,
				Msg:    code.GetMsg(co),
				Error:  errors.New("数据库错误").Error(),
			}
		}
		return serializer.Response{
			Status: co,
			Msg:    code.GetMsg(co),
		}
	} else {
		userDao := dao.NewUserDao(db.DB)
		err := userDao.Update(uid, u.SlotId, u.Name, u.Direction)
		if err != nil {
			co = code.Error
			zap.L().Info("logic/user.go Update failed error : ", zap.Error(err))
			log.Println("logic/user.go Update failed error : ", err)
			return serializer.Response{
				Status: co,
				Msg:    code.GetMsg(co),
				Error:  errors.New("数据库错误").Error(),
			}
		}
		return serializer.Response{
			Status: co,
			Msg:    code.GetMsg(co),
		}
	}
}
func ShowSlot(uid uint) serializer.Response {
	co := code.Success
	var err error
	var slotId uint
	var slot *model.InterviewSlot
	slotDao := dao.NewSlotDao(db.DB)
	userDao := dao.NewUserDao(db.DB)
	slotId, err = userDao.GetUserCurrentValidSlotID(uid)
	if err != nil {
		co = code.Error
		return serializer.Response{
			Status: co,
			Msg:    "获取面试时段失败",
			Error:  err.Error(),
		}
	}
	slot, err = slotDao.GetSlotById(slotId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return serializer.Response{
				Status: 200, // 注意：业务上“无数据”仍返回 200
				Data:   nil, // Data 为 nil 表示无预约
				Msg:    "暂无面试安排",
			}
		}
		// 数据库错误
		return serializer.Response{
			Status: 500,
			Data:   nil,
			Msg:    "数据库查询失败",
			Error:  err.Error(),
		}
	}
	fmt.Println("logic : ", slot)
	return serializer.Response{
		Status: co,
		Data:   slot,
		Msg:    code.GetMsg(co),
		Error:  "",
	}
}
