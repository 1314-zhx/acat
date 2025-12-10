package logic

import (
	"acat/dao"
	"acat/dao/db"
	"acat/model"
	"acat/model/code"
	"acat/serializer"
	"acat/util"
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

//--------->管理员登录<-----------//
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

//---------->管理员设置面试时间表<-------------//
type SetSlot struct {
	StartTime util.DateTimeLocal `json:"start_time"`
	EndTime   util.DateTimeLocal `json:"end_time"`
	Round     int                `json:"round"`
	MaxNum    int                `json:"max_num"`
}

const (
	MinRound  = 1
	MaxRound  = 2
	MinMaxNum = 1
	MaxMaxNum = 50
)

func (s *SetSlot) SetSchedule() serializer.Response {
	co := code.Success
	fmt.Println("时间设置错误")
	if s.StartTime.Time().After(s.EndTime.Time()) {
		co = code.InvalidParam
		return serializer.Response{
			Status: co,
			Msg:    code.GetMsg(co),
			Error:  "无效参数，时间设置有误",
		}
	}
	fmt.Println("轮次错误")
	if s.Round > MaxRound || s.Round < MinRound {
		co = code.InvalidParam
		return serializer.Response{
			Status: co,
			Msg:    code.GetMsg(co),
			Error:  "无效参数，面试轮次不对",
		}
	}
	fmt.Println("人数错误")
	if s.MaxNum > MaxMaxNum || s.MaxNum < MinMaxNum {
		co = code.InvalidParam
		return serializer.Response{
			Status: co,
			Msg:    code.GetMsg(co),
			Error:  "无效参数，最大人数不对",
		}
	}
	slotDao := dao.NewSlotDao(db.DB)
	// 构造成结构体减少传参
	slot := &model.InterviewSlot{
		StartTime: s.StartTime.Time(),
		EndTime:   s.EndTime.Time(),
		Round:     s.Round,
		Num:       0,
		MaxNum:    s.MaxNum,
	}
	fmt.Println("时间段错误")
	err := slotDao.Create(slot)
	if err != nil {
		co = code.Error
		zap.L().Info("logic/admin.go slotDao.Create failed error : ", zap.Error(err))
		log.Println("logic/admin.go slotDao.Create failed error : ", err)
		return serializer.Response{
			Status: co,
			Msg:    code.GetMsg(co),
			Error:  "创建时间段错误",
		}
	}
	fmt.Println("over")
	return serializer.Response{
		Status: co,
		Msg:    code.GetMsg(co),
	}
}

//------------>管理员设置面试结果<---------------//
type SetResult struct {
	Round  int  `json:"round"`
	SlotId uint `json:"slot_id"`
}

func (s *SetResult) SetUserResult(ctx context.Context) serializer.Response {
	co := code.Success
	var err error
	// 判断数据是否合法
	if s.Round < MinRound || s.Round > MaxRound {
		co = code.InvalidParam
		return serializer.Response{
			Status: co,
			Msg:    code.GetMsg(co),
			Error:  "无效参数",
		}
	}
	// 判断有没有该SlotId
	var slot *model.InterviewSlot
	slotDao := dao.NewSlotDao(db.DB)
	slot, err = slotDao.GetSlotById(s.SlotId)
	if err != nil {
		co = code.Error
		zap.L().Info("logic/admin.go SetUserResult failed GetSlotById error : ", zap.Error(err))
		log.Println("logic/admin.go SetUserResult failed GetSlotById error : ", err)
		return serializer.Response{
			Status: co,
			Msg:    code.GetMsg(co),
			Error:  "获取面试表错误",
		}
	}
	if slot == nil {
		co = code.Error
		zap.L().Info("logic/admin.go SetUserResult not exists slot")
		log.Println("没找的slot")
		return serializer.Response{
			Status: co,
			Msg:    code.GetMsg(co),
			Error:  "没有该面试表",
		}
	}
	// 根据该SlotId返回该场面试的所有用户id
	var userIds []uint
	userIds, err = slotDao.GetUserIdsBySlotId(slot.ID)
	if err != nil {
		co = code.Error
		zap.L().Info("logic/admin.go SetUserResult failed get userIds error : ", zap.Error(err))
		log.Println("logic/admin.go SetUserResult failed get userIds error : ", err)
		return serializer.Response{
			Status: co,
			Msg:    code.GetMsg(co),
			Error:  "获取用户Ids失败",
		}
	}
	// 根据userId包装所有用户的信息
	userDao := dao.NewUserDao(db.DB)
	var users []*model.UserResponse
	for _, userId := range userIds {
		user, err := userDao.GetUserByID(userId, ctx)
		if err != nil || user == nil {
			zap.L().Warn("User not found or DB error", zap.Uint("user_id", userId), zap.Error(err))
			continue
		}
		userRes := &model.UserResponse{
			ID:     userId,
			Name:   user.Name,
			StuId:  user.StuId,
			Phone:  user.Phone,
			Gender: user.Gender,
		}
		users = append(users, userRes)
	}

	return serializer.Response{
		Status: co,
		Data:   users,
		Msg:    code.GetMsg(co),
	}
}

type SetPass struct {
	UserId uint `json:"user_id"`
	SlotId uint `json:"slot_id"`
	Round  int  `json:"round"`
	Pass   int  `json:"pass"` // 0 是没过，1 是过
}

// 需要修改UserModel,如果过了还需要加一个InterviewResult表
func (s *SetPass) SetUserPass(adminId uint) serializer.Response {
	co := code.Success
	userDao := dao.NewUserDao(db.DB)
	if s.Round == 2 {
		user, err := userDao.GetUserByID(s.UserId, context.Background())
		if err != nil {
			co = code.Error
			zap.L().Info("logic/admin.go SetUserPass 查找用户失败")
			return serializer.Response{
				Status: co,
				Msg:    code.GetMsg(co),
				Error:  "查找用户失败",
			}
		}
		if user.FirstPass == 2 {
			co = code.Error
			return serializer.Response{
				Status: co,
				Msg:    code.GetMsg(co),
				Error:  "用户未通过一面无法设置二面结果",
			}
		}
	}
	err := userDao.UpdatePass(s.UserId, s.Round, s.Pass, adminId)
	if err != nil {
		co = code.Error
		zap.L().Info("更新面试结果失败", zap.Error(err))
		log.Println("更新面试结果失败", err)
		return serializer.Response{
			Status: co,
			Msg:    code.GetMsg(co),
			Error:  "更新面试结果失败",
		}
	}
	return serializer.Response{
		Status: co,
		Msg:    code.GetMsg(co),
	}
}

//---------->管理员公布面试结果<-------------//
type PassUser struct {
	UserId    uint   `json:"user_id"`
	Name      string `json:"name"`
	Round     int    `json:"round"`
	Email     string `json:"email"`
	Customize bool   `json:"customize"` // 自定义开关，true为自定义发件信息，false为默认发件信息
	Content   string `json:"content"`   // 自定义内容
}

func (p *PassUser) Publish() serializer.Response {
	co := code.Success
	fmt.Println(p)
	if p.Customize == true {
		err := PublicCustomEmail(p.Email, p.Content)
		if err != nil {
			co = code.Error
			return serializer.Response{
				Status: co,
				Msg:    code.GetMsg(co),
				Error:  "发送面试结果错误",
			}
		}
	} else {
		err := PublicEmail(p.Email, p.Round, p.Name)
		if err != nil {
			co = code.Error
			return serializer.Response{
				Status: co,
				Msg:    code.GetMsg(co),
				Error:  "发送面试结果错误",
			}
		}
	}
	return serializer.Response{
		Status: co,
		Msg:    code.GetMsg(co),
	}
}

type PublishRound struct {
	Round int `json:"round"`
}

func (p *PublishRound) GetPassUser() serializer.Response {
	co := code.Success
	userDao := dao.NewUserDao(db.DB)
	users, err := userDao.GetUsersByRound(p.Round)
	if err != nil {
		co = code.Error
		zap.L().Info("logic/admin.go GetUsersByRound failed error : ", zap.Error(err))
		log.Println("logic/admin.go GetUsersByRound failed error : ", err)
		return serializer.Response{
			Status: co,
			Msg:    code.GetMsg(co),
			Error:  "获取通过学生信息错误",
		}
	}
	return serializer.Response{
		Status: co,
		Msg:    code.GetMsg(co),
		Data:   users,
	}
}

//----------->管理员查看邮箱并回信<------------//
type AdminLetter struct {
}

func (a *AdminLetter) Letter(aid uint) serializer.Response {
	co := code.Success
	letterDao := dao.NewLetterDao(db.DB)
	letters, err := letterDao.GetLetters(aid)
	if err != nil {
		co = code.Error
		zap.L().Info("logic/admin.go letter get error : ", zap.Error(err))
		log.Println("logic/admin.go letter get error : ", err)
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

type AdminReply struct {
	LetterId int64  `json:"letter_id"`
	AdminId  uint   `json:"admin_id"`
	UserId   uint   `json:"user_id"`
	Content  string `json:"content"`
	Title    string `json:"title"`
	IsRead   bool   `json:"is_read"`
}

func (a *AdminReply) Reply() serializer.Response {
	fmt.Println("logic", a)
	co := code.Success
	letterDao := dao.NewLetterDao(db.DB)
	if a.Content == "" {
		err := letterDao.SetIsRead(a.LetterId, a.IsRead)
		if err != nil {
			zap.L().Info("logic/admin.go Reply failed error : ", zap.Error(err))
			log.Println("logic/admin.go Reply failed error : ", err)
			return serializer.Response{
				Status: co,
				Msg:    code.GetMsg(co),
				Error:  "管理员标记已读失败",
			}
		}
	} else {
		err := letterDao.Reply(a.AdminId, a.UserId, a.Title, a.Content)
		if err != nil {
			zap.L().Info("logic/admin.go Reply failed error : ", zap.Error(err))
			log.Println("logic/admin.go Reply failed error : ", err)
			return serializer.Response{
				Status: co,
				Msg:    code.GetMsg(co),
				Error:  "管理员回信失败",
			}
		}
	}
	return serializer.Response{
		Status: co,
		Msg:    code.GetMsg(co),
	}
}

//------------>管理员上传电子版试题<------------//
const uploadDir = "./uploads"

var allowedExts = map[string]bool{
	".pdf":  true,
	".docx": true,
	".doc":  true,
	".txt":  true,
	".md":   true,
}

type Upload struct{}

// UploadQuestion 接收 *multipart.FileHeader，不依赖 gin.Context
func (u *Upload) UploadQuestion(file *multipart.FileHeader) serializer.Response {
	co := code.Success
	err := os.MkdirAll(uploadDir, 0755)
	// 确保上传目录存在
	if err != nil {
		co = code.Error
		zap.L().Error("创建上传目录失败", zap.Error(err))
		return serializer.Response{
			Status: co,
			Msg:    code.GetMsg(co),
			Error:  "服务器内部错误：无法创建上传目录",
		}
	}
	// 删除原有文件
	oldFiles, err := os.ReadDir(uploadDir)
	if err == nil { // 如果目录存在且可读
		for _, f := range oldFiles {
			if !f.IsDir() {
				os.Remove(filepath.Join(uploadDir, f.Name()))
			}
		}
	}
	// 获取安全的文件名并校验扩展名
	filename := filepath.Base(file.Filename)
	ext := strings.ToLower(filepath.Ext(filename))

	if !allowedExts[ext] {
		co = code.Error
		zap.L().Warn("上传了不支持的文件类型", zap.String("ext", ext), zap.String("filename", filename))
		return serializer.Response{
			Status: co,
			Msg:    code.GetMsg(co),
			Error:  "不支持的文件类型，请上传 PDF、Word 或 TXT 文档",
		}
	}

	// 构建保存路径
	savePath := filepath.Join(uploadDir, filename)

	// 打开上传的文件流
	src, err := file.Open()
	if err != nil {
		co = code.Error
		zap.L().Error("无法打开上传的文件", zap.Error(err))
		return serializer.Response{
			Status: co,
			Msg:    code.GetMsg(co),
			Error:  "文件读取失败",
		}
	}
	defer src.Close()

	// 创建目标文件
	dst, err := os.Create(savePath)
	if err != nil {
		co = code.Error
		zap.L().Error("无法创建目标文件", zap.Error(err), zap.String("path", savePath))
		return serializer.Response{
			Status: co,
			Msg:    code.GetMsg(co),
			Error:  "文件保存失败",
		}
	}
	defer dst.Close()
	// 复制内容
	_, err = io.Copy(dst, src)
	if err != nil {
		co = code.Error
		zap.L().Error("文件写入失败", zap.Error(err))
		return serializer.Response{
			Status: co,
			Msg:    code.GetMsg(co),
			Error:  "文件写入失败",
		}
	}
	return serializer.Response{
		Status: co,
		Msg:    code.GetMsg(co),
		Data:   map[string]string{"filename": filename},
	}
}
