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
	"acat/util"
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
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
			Error:  "缺少必填字段",
		}, nil
	}
	if len(l.Phone) != 11 {
		co = code.InvalidPhoneForm
		return serializer.Response{
			Status: co,
			Msg:    code.GetMsg(co),
			Error:  "参数格式错误",
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
			Error:  "手机号或密码错误",
		}, nil
	}
	// 2 比对密码（用查到的 userModel.Password）
	if err := bcrypt.CompareHashAndPassword([]byte(userModel.Password), []byte(l.Password)); err != nil {
		co = code.PhoneOrPasswordError
		zap.L().Info("密码错误", zap.String("phone", l.Phone))
		return serializer.Response{
			Status: co,
			Msg:    "手机号或密码错误",
			Error:  "手机号或密码错误",
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
	Email      string `json:"email"`
	Code       string `json:"code"` // 防止用户邮箱填错
	StuId      string `json:"stu_id"`
	Gender     int    `json:"gender"`    // 1为男2为女
	Direction  int    `json:"direction"` // 方向选择，0不确定，1为Go，2为Java，3为前端，4为后端
}

// SendRegisterCode 用于发送注册验证码
func (r *UserRes) SendRegisterCode() serializer.Response {
	// 必填字段（用于发码，不需要密码）
	co := code.Success
	if r.Phone == "" || r.StuId == "" || r.Name == "" || r.Gender == 0 || r.Email == "" {
		co = code.MissMustInfo
		zap.L().Info("logic/user.go SendRegisterCode() failed miss must info : ", zap.Error(errors.New(code.GetMsg(co))))
		log.Println("logic/user.go SendRegisterCode() failed miss must info : ", errors.New(code.GetMsg(co)))
		return serializer.Response{
			Status: co,
			Msg:    code.GetMsg(co),
			Error:  "缺少必填字段",
		}
	}

	// 格式校验
	if len(r.Phone) != 11 || !util.IsPhone(r.Phone) {
		co = code.InvalidPhoneForm
		zap.L().Info("logic/user.go SendRegisterCode() failed phone form error : ", zap.Error(errors.New(code.GetMsg(co))))
		log.Println("logic/user.go SendRegisterCode() failed phone form error : ", errors.New(code.GetMsg(co)))
		return serializer.Response{
			Status: co,
			Msg:    code.GetMsg(co),
			Error:  "格式不正确",
		}
	}
	if !util.IsEmail(r.Email) {
		co = code.EmailFormError
		zap.L().Info("logic/user.go SendRegisterCode() failed email form error : ", zap.Error(errors.New(code.GetMsg(co))))
		log.Println("logic/user.go SendRegisterCode() failed email form error : ", errors.New(code.GetMsg(co)))
		return serializer.Response{
			Status: co,
			Msg:    code.GetMsg(co),
			Error:  "格式不正确",
		}
	}

	userDao := dao.NewUserDao(db.DB)

	// 唯一性校验，每个用户对应一个邮箱，后续的改密码改手机号等操作由邮箱发送验证码
	if !userDao.NotExistsUserByPhone(r.Phone) {
		co = code.Error
		zap.L().Info("logic/user.go SendRegisterCode() failed exists user  error : ", zap.Error(errors.New(code.GetMsg(co))))
		log.Println("logic/user.go SendRegisterCode() failed exists user error : ", errors.New(code.GetMsg(co)))
		return serializer.Response{
			Status: co,
			Msg:    code.GetMsg(co),
			Error:  "该手机号已被注册",
		}
	}
	if !userDao.NotExistsUserByEmail(r.Email) {
		co = code.Error
		zap.L().Info("logic/user.go SendRegisterCode() failed exists user  error : ", zap.Error(errors.New(code.GetMsg(co))))
		log.Println("logic/user.go SendRegisterCode() failed exists user error : ", errors.New(code.GetMsg(co)))
		return serializer.Response{
			Status: co,
			Msg:    code.GetMsg(co),
			Error:  "该邮箱已被注册",
		}
	}
	// 生成并发送验证码
	vcode := GenerateCode()
	if err := SendEmailCode(r.Email, vcode); err != nil {
		co = code.Error
		zap.L().Error("发送注册验证码失败", zap.String("email", r.Email), zap.Error(err))
		return serializer.Response{
			Status: co,
			Msg:    code.GetMsg(co),
			Error:  "发送注册验证码失败",
		}
	}

	// 存入 Redis，5 分钟有效
	key := "register:email:" + r.Email
	ctx := context.Background()
	if err := redislock.GetRDB().Set(ctx, key, vcode, 5*time.Minute).Err(); err != nil {
		co = code.Error
		zap.L().Error("Redis 存储验证码失败", zap.Error(err))
		// 即使存储失败，邮件已发，可选择继续或报错；这里建议报错以保证一致性
		return serializer.Response{
			Status: co,
			Msg:    code.GetMsg(co),
			Error:  "内部错误",
		}
	}
	return serializer.Response{
		Status: code.Success,
		Msg:    "验证码已发送，请查收邮箱",
	}
}

// CompleteRegister 使用验证码完成注册
func (r *UserRes) CompleteRegister() serializer.Response {
	co := code.Success
	// 必填字段（含验证码）
	if r.Phone == "" || r.Password == "" || r.RePassword == "" || r.StuId == "" ||
		r.Name == "" || r.Gender == 0 || r.Email == "" || r.Code == "" {
		co = code.MissMustInfo
		zap.L().Info("logic/user.go CompleteRegister() failed miss must info : ", zap.Error(errors.New(code.GetMsg(co))))
		log.Println("logic/user.go CompleteRegister() failed miss must info : ", errors.New(code.GetMsg(co)))
		return serializer.Response{
			Status: co,
			Msg:    code.GetMsg(co),
			Error:  "缺少必填字段",
		}
	}

	// 密码一致性
	if r.Password != r.RePassword {
		co = code.Error
		return serializer.Response{
			Status: co,
			Msg:    code.GetMsg(co),
			Error:  "两次密码不一致",
		}
	}

	// 格式校验
	if len(r.Phone) != 11 || !util.IsPhone(r.Phone) {
		co = code.Error
		return serializer.Response{
			Status: co,
			Msg:    code.GetMsg(co),
			Error:  "手机号格式不正确",
		}
	}
	if !util.IsEmail(r.Email) {
		co = code.Error
		return serializer.Response{
			Status: co,
			Msg:    code.GetMsg(co),
			Error:  "邮箱格式不正确",
		}
	}

	// 验证码校验
	key := "register:email:" + r.Email
	ctx := context.Background()
	storedCode, err := redislock.GetRDB().Get(ctx, key).Result()
	if err != nil {
		co = code.Error
		if err == redis.Nil {
			return serializer.Response{
				Status: co,
				Msg:    code.GetMsg(co),
				Error:  "验证码已过期或未请求",
			}
		}
		zap.L().Error("Redis 读取验证码失败", zap.Error(err))
		return serializer.Response{
			Status: co,
			Msg:    code.GetMsg(co),
			Error:  "服务异常",
		}
	}

	if storedCode != r.Code {
		co = code.Error
		return serializer.Response{
			Status: co,
			Msg:    code.GetMsg(co),
			Error:  "验证码错误",
		}
	}
	userDao := dao.NewUserDao(db.DB)
	// 最终唯一性校验（防并发注册）
	if !userDao.NotExistsUserByPhone(r.Phone) {
		co = code.Error
		zap.L().Info("logic/user.go Register() failed exists user  error : ", zap.Error(errors.New(code.GetMsg(co))))
		log.Println("logic/user.go Register() failed exists user error : ", errors.New(code.GetMsg(co)))
		return serializer.Response{
			Status: co,
			Msg:    code.GetMsg(co),
			Error:  "该手机号已被注册",
		}
	}
	if !userDao.NotExistsUserByEmail(r.Email) {
		co = code.Error
		zap.L().Info("logic/user.go Register() failed exists user  error : ", zap.Error(errors.New(code.GetMsg(co))))
		log.Println("logic/user.go Register() failed exists user error : ", errors.New(code.GetMsg(co)))
		return serializer.Response{
			Status: co,
			Msg:    code.GetMsg(co),
			Error:  "该邮箱已被注册",
		}
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(r.Password), bcrypt.DefaultCost)
	if err != nil {
		co = code.Error
		zap.L().Error("密码加密失败", zap.Error(err))
		return serializer.Response{
			Status: co,
			Msg:    code.GetMsg(co),
			Error:  "内部错误",
		}
	}

	// 创建用户
	userRegister := &model.UserRegister{
		Name:      r.Name,
		StuId:     r.StuId,
		Password:  string(hashedPassword),
		Phone:     r.Phone,
		Email:     r.Email,
		Gender:    r.Gender,
		Direction: r.Direction,
	}

	if err := userDao.Create(userRegister); err != nil {
		co = code.Error
		zap.L().Error("用户注册写入数据库失败", zap.Error(err))
		return serializer.Response{
			Status: co,
			Data:   nil,
			Msg:    "",
			Error:  "",
		}
	}

	// 清除验证码（防止重复使用）
	_ = redislock.GetRDB().Del(ctx, key)

	return serializer.Response{
		Status: code.Success,
		Msg:    "注册成功",
	}
}

//----------->用户查询面试结果<--------------//

type UserChe struct {
	Round int `json:"round"`
}

func (ch *UserChe) Result(uid uint, ctx context.Context) serializer.Response {
	co := code.Success
	if ch.Round > 2 || ch.Round < 1 {
		co = code.InvalidParam
		zap.L().Warn("logic/user.go Result failed error : ", zap.Error(errors.New("无效参数")))
		log.Println("logic/user.go Result failed error : ", "无效参数")
		return serializer.Response{
			Status: co,
			Msg:    code.GetMsg(co),
			Error:  "无效参数",
		}
	}
	userDao := dao.NewUserDao(db.DB)
	pass, err := userDao.Result(uid, ch.Round, ctx)
	if pass == 0 && err == nil {
		return serializer.Response{
			Status: co,
			Data:   "学生未通过面试",
			Msg:    code.GetMsg(co),
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
		co = code.Error
		zap.L().Info("logic/user.go Post get user by id error : ", zap.Error(err))
		log.Println("logic/user.go Post get user by id error : ", "查询用户失败")
		return serializer.Response{
			Status: co,
			Msg:    code.GetMsg(co),
			Error:  "查询用户失败",
		}
	}
	slotDao := dao.NewSlotDao(db.DB)
	var slot *model.InterviewSlot
	slot, err = slotDao.GetSlotById(p.SlotId)
	if err != nil {
		co = code.Error
		zap.L().Info("logic/user.go Post get slot by id error : ", zap.Error(err))
		log.Println("logic/user.go Post get slot by id error : ", "查询面试表失败")
		return serializer.Response{
			Status: co,
			Msg:    code.GetMsg(co),
			Error:  "查询面试表失败",
		}
	}
	err = userDao.UserOnlyOneSlot(uid)
	if err != nil {
		co = code.Error
		if err.Error() == "用户已有其它面试" {
			zap.L().Info("logic/user.go Post UserOnlyOneSlot error : ", zap.Error(err))
			log.Println("logic/user.go Post UserOnlyOneSlot error : ", "用户已有其它面试")
			return serializer.Response{
				Status: co,
				Msg:    code.GetMsg(co),
				Error:  "用户已有其它面试",
			}
		}
		zap.L().Info("logic/user.go Post UserOnlyOneSlot error : ", zap.Error(err))
		log.Println("logic/user.go Post UserOnlyOneSlot error : ", "查询用户面试关系表失败")
		return serializer.Response{
			Status: co,
			Msg:    code.GetMsg(co),
			Error:  "查询用户面试关系表失败",
		}
	}
	// 用户报二面，但一面没过
	if slot.Round == 2 && user.FirstPass != 1 {
		co = code.FirstViewNotPass
		zap.L().Info("logic/user.go Post round 1 no pass post 2 error : ", zap.Error(errors.New("一面未通过，无法参加二面")))
		log.Println("logic/user.go Post round 1 no pass post 2 error : ", "一面未通过，无法参加二面")
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
	slotDao := dao.NewSlotDao(db.DB)
	slot, err := slotDao.GetSlotById(u.SlotId)
	if err != nil {
		co = code.Error
		if slot == nil {
			return serializer.Response{
				Status: co,
				Msg:    code.GetMsg(co),
				Error:  "没有该面试时段",
			}
		}
		return serializer.Response{
			Status: co,
			Msg:    code.GetMsg(co),
			Error:  "查询面试时段错误",
		}
	}
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
	} else if u.IsDelete == 0 {
		fmt.Println("111")
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
	return serializer.Response{
		Status: code.Error,
		Msg:    "无效更改",
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

//------------>用户忘记密码<--------------//
type Forget struct {
	Param     string `json:"param"`
	Test_mode bool   `json:"test_mode"`
}

func (f *Forget) Forget() serializer.Response {
	// 获取6位验证码
	vCode := GenerateCode()

	// 构造 redis key
	var key string
	if f.Param == "" {
		return serializer.Response{
			Status: code.InvalidParam,
			Msg:    code.GetMsg(code.InvalidParam),
			Error:  "空传递",
		}
	}
	key = "verify:email:" + f.Param
	if util.IsEmail(f.Param) && f.Test_mode == false {
		err := SendEmailCode(f.Param, vCode)
		if err != nil {
			zap.L().Info("logic/user.go SendEmailCode failed error : ", zap.Error(err))
			log.Println("logic/user.go SendEmailCode failed error : ", err)
			return serializer.Response{
				Status: code.Error,
				Msg:    code.GetMsg(code.Error),
				Error:  "发送邮件失败",
			}
		}
	} else if !util.IsEmail(f.Param) {
		return serializer.Response{
			Status: code.InvalidParam,
			Msg:    code.GetMsg(code.InvalidParam),
			Error:  "无效参数",
		}
	}
	ctx := context.Background()
	// 将验证码存入redis5分钟
	if err := redislock.GetRDB().Set(ctx, key, vCode, 5*time.Minute).Err(); err != nil {
		zap.L().Error("Redis 存储验证码失败", zap.Error(err))
	}
	//if f.Test_mode == true {
	//	return serializer.Response{
	//		Status: code.Success,
	//		Data:   vCode,
	//		Msg:    code.GetMsg(code.Success),
	//		Error:  "",
	//	}
	//}
	return serializer.Response{
		Status: code.Success,
		Msg:    code.GetMsg(code.Success),
	}
}

type ResetPassword struct {
	Account     string `json:"account"`
	Code        string `json:"code"`
	NewPassword string `json:"new_password"`
}

func (r *ResetPassword) ResetPassword(ctx context.Context) serializer.Response {
	co := code.Success
	// 1. 参数校验
	fmt.Println(r)
	if r.Account == "" || r.Code == "" || r.NewPassword == "" {
		zap.L().Warn("logic.user.go ResetPassword failed", zap.Error(errors.New("缺少必要参数")))
		log.Println("logic.user.go ResetPassword failed ", "缺少必要参数")
		return serializer.Response{
			Status: code.InvalidParam,
			Msg:    code.GetMsg(code.InvalidParam),
			Error:  "缺少必要参数",
		}
	}
	// 2. 构造 Redis key（和 Forget 里完全一致！）
	var key string
	if util.IsEmail(r.Account) {
		key = "verify:email:" + r.Account
	} else {
		zap.L().Warn("logic.user.go ResetPassword failed", zap.Error(errors.New("无效账号格式")))
		log.Println("logic.user.go ResetPassword failed ", "无效账号格式")
		return serializer.Response{
			Status: code.InvalidParam,
			Msg:    "无效账号格式",
			Error:  "无效账号格式",
		}
	}
	// 3. 从 Redis 读取验证码
	storedCode, err := redislock.GetRDB().Get(ctx, key).Result()
	fmt.Println("storedCode", storedCode)
	if err != nil {
		if err == redis.Nil {
			zap.L().Warn("logic.user.go ResetPassword failed", zap.Error(errors.New("无效验证码")))
			log.Println("logic.user.go ResetPassword failed ", "无效验证码")
			return serializer.Response{
				Status: code.Error,
				Msg:    "验证码已过期或未请求",
				Error:  "无效验证码",
			}
		}
		zap.L().Error("Redis 读取失败", zap.Error(err))
		zap.L().Warn("logic.user.go ResetPassword failed", zap.Error(errors.New("内部错误")))
		log.Println("logic.user.go ResetPassword failed ", "内部错误")
		return serializer.Response{
			Status: code.Error,
			Msg:    "服务异常",
			Error:  "内部错误",
		}
	}
	// 4. 验证码比对
	if storedCode != r.Code {
		zap.L().Warn("logic.user.go ResetPassword failed", zap.Error(errors.New("验证码不匹配")))
		log.Println("logic.user.go ResetPassword failed ", "验证码不匹配")
		return serializer.Response{
			Status: code.Error,
			Msg:    "验证码错误",
			Error:  "验证码不匹配",
		}
	}
	// 5. 查找用户（通过 account）
	userDao := dao.NewUserDao(db.DB)
	var user *model.UserModel
	if util.IsEmail(r.Account) {
		user, err = userDao.GetUserByEmail(r.Account, ctx)
		if err != nil {
			co = code.Error
			zap.L().Info("logic/user.go ResetPassword failed error : ", zap.Error(err))
			return serializer.Response{
				Status: co,
				Data:   nil,
				Msg:    "",
				Error:  "",
			}
		}
	}
	if user == nil {
		co = code.Error
		zap.L().Warn("logic.user.go ResetPassword failed", zap.Error(errors.New("未找到该用户")))
		log.Println("logic.user.go ResetPassword failed ", "未找到该用户")
		return serializer.Response{
			Status: co,
			Data:   nil,
			Msg:    "未找到该用户",
			Error:  "未找到该用户",
		}
	}
	// 6. 加密新密码
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(r.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		zap.L().Warn("logic.user.go ResetPassword failed", zap.Error(errors.New("加密错误")))
		log.Println("logic.user.go ResetPassword failed ", "加密错误")
		return serializer.Response{
			Status: code.Error,
			Msg:    "密码处理失败",
			Error:  "加密错误",
		}
	}

	// 7. 更新数据库
	err = userDao.ResetPassword(ctx, user.ID, string(hashedPwd))
	if err != nil {
		co = code.Error
		zap.L().Warn("logic.user.go ResetPassword failed", zap.Error(errors.New("重置密码失败")))
		log.Println("logic.user.go ResetPassword failed ", "重置密码失败")
		return serializer.Response{
			Status: co,
			Msg:    code.GetMsg(co),
			Error:  "重置密码失败",
		}
	}
	// 8. 删除 Redis 验证码（一次性使用）
	redislock.GetRDB().Del(ctx, key)

	return serializer.Response{
		Status: co,
		Msg:    code.GetMsg(co),
	}
}

//------------>用户查看来信<--------------//
type CheckReply struct{}

func (c *CheckReply) Reply(uid uint) serializer.Response {
	co := code.Success
	letterDao := dao.NewLetterDao(db.DB)
	letters, err := letterDao.GetLetters(uid)
	if err != nil {
		co = code.Error
		zap.L().Info("logic/user.go Reply get error : ", zap.Error(err))
		log.Println("logic/user.go Reply get error : ", err)
		return serializer.Response{
			Status: co,
			Msg:    code.GetMsg(co),
			Error:  "获取信件失败",
		}
	}
	return serializer.Response{
		Status: co,
		Data:   letters,
		Msg:    code.GetMsg(co),
	}
}

//----------->用户下载面试试题<------------//
type Download struct {
}

// DownloadQuestion 接收文件名，返回文件路径或错误
func (d *Download) DownloadQuestion(filename string) serializer.Response {
	// 如果 filename 为空，自动查找 uploads 目录中第一个合法文件
	if filename == "" {
		files, err := os.ReadDir(uploadDir)
		if err != nil || len(files) == 0 {
			return serializer.Response{
				Status: code.Error,
				Msg:    code.GetMsg(code.Error),
				Error:  "文件不存在",
			}
		}

		var foundName string
		for _, f := range files {
			if !f.IsDir() {
				name := f.Name()
				ext := strings.ToLower(filepath.Ext(name))
				if allowedExts[ext] {
					foundName = name
					break
				}
			}
		}

		if foundName == "" {
			return serializer.Response{
				Status: code.Error,
				Msg:    code.GetMsg(code.Error),
				Error:  "文件不存在",
			}
		}
		filename = foundName
	}

	// 防止路径穿越或非法字符
	if strings.ContainsAny(filename, "/\\") || strings.Contains(filename, "..") {
		zap.L().Warn("拒绝非法文件名", zap.String("filename", filename))
		return serializer.Response{
			Status: code.Error,
			Msg:    code.GetMsg(code.Error),
			Error:  "无效的文件名",
		}
	}

	// 校验扩展名
	ext := strings.ToLower(filepath.Ext(filename))
	if !allowedExts[ext] {
		zap.L().Warn("尝试下载不支持的文件类型", zap.String("filename", filename))
		return serializer.Response{
			Status: code.Error,
			Msg:    code.GetMsg(code.Error),
			Error:  "不支持的文件类型",
		}
	}

	// 构建绝对路径,避免相对路径和工作目录问题
	filePath := filepath.Join(uploadDir, filename)
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		zap.L().Error("无法解析文件绝对路径", zap.String("path", filePath), zap.Error(err))
		return serializer.Response{
			Status: code.Error,
			Msg:    code.GetMsg(code.Error),
			Error:  "内部错误",
		}
	}

	// 文件是否真实存在
	_, err = os.Stat(absPath)
	if os.IsNotExist(err) {
		zap.L().Warn("文件不存在（logic 层校验）", zap.String("path", absPath))
		return serializer.Response{
			Status: code.Error,
			Msg:    code.GetMsg(code.Error),
			Error:  "文件不存在",
		}
	}

	// 成功返回
	return serializer.Response{
		Status: code.Success,
		Msg:    code.GetMsg(code.Success),
		Data:   map[string]string{"file_path": absPath, "filename": filename},
	}
}
