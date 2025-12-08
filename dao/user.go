package dao

import (
	"acat/model"
	"acat/redislock"
	"acat/util"
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"log"
	"strings"
	"time"
)

type UserDao struct {
	db *gorm.DB
}
type UserDaoWithRdb struct {
	db  *gorm.DB
	rdb *redis.Client
}

func NewUserDao(db *gorm.DB) *UserDao {
	return &UserDao{db: db}
}
func NewUserDaoWithRdb(db *gorm.DB, rdb *redis.Client) *UserDaoWithRdb {
	return &UserDaoWithRdb{db: db, rdb: rdb}
}

// GetUserByPhone 根据手机号获取用户
func (dao *UserDao) GetUserByPhone(phone string, ctx context.Context) (*model.UserModel, error) {
	var user model.UserModel
	err := dao.db.WithContext(ctx).Where("phone = ?", phone).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (dao *UserDao) NotExistsUserByPhone(phone string) bool {
	var count int64
	err := dao.db.Model(&model.UserModel{}).
		Where("phone = ?", phone).
		Count(&count).Error
	if err != nil {
		zap.L().Error("dao/user.go NotExistsUserByPhone failed", zap.Error(err))
		log.Printf("检查手机号是否存在失败: %v", err)
		return false
	}
	return count == 0
}
func (dao *UserDao) GetUserByID(id uint, ctx context.Context) (*model.UserModel, error) {
	var user model.UserModel
	err := dao.db.WithContext(ctx).Select("id", "name", "password", "phone", "first_pass").
		Where("id = ?", id).
		First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}
func (dao *UserDao) Create(register *model.UserRegister) error {
	user := model.UserModel{
		StuId:     register.StuId,
		Name:      register.Name,
		Password:  register.Password,
		Phone:     register.Phone,
		Email:     register.Email,
		Direction: register.Direction,
		Gender:    register.Gender,
	}
	return dao.db.Create(&user).Error
}
func (dao *UserDao) GetUserByEmail(email string, ctx context.Context) (*model.UserModel, error) {
	var user model.UserModel
	err := dao.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}
func (dao *UserDao) ResetPassword(ctx context.Context, userID uint, hashedPassword string) error {
	err := dao.db.WithContext(ctx).
		Model(&model.UserModel{}).
		Where("id = ?", userID).
		Update("password", hashedPassword).Error
	return err
}
func (dao *UserDao) Result(uid uint, round int, ctx context.Context) (int, error) {
	var userResult model.InterviewResult
	err := dao.db.WithContext(ctx).Where("id = ? and round = ?", uid, round).First(&userResult).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, nil
		}
		zap.L().Error("Query interview result failed", zap.Error(err), zap.Uint("uid", uid))
		return 0, errors.New("数据库异常")
	}
	return userResult.Status, nil
}

// Post 用于提交用户报名申请，分布式锁加事务保证同一时间只有一个用户操作，并且错误后会回滚。
func (dao *UserDaoWithRdb) Post(sid uint, uid uint, ctx context.Context) error {
	lockKey := fmt.Sprintf("redislock:slot:%d", sid)
	redisLock := redislock.NewRedisLock(dao.rdb, lockKey, 10*time.Second)

	// 1. 获取分布式锁
	locked, err := redisLock.Lock(ctx)
	if err != nil {
		return fmt.Errorf("加锁失败: %w", err)
	}
	if !locked {
		return errors.New("当前时段繁忙，请稍后重试")
	}
	defer func() {
		_ = redisLock.Unlock(ctx)
	}()

	// 2. 开启数据库事务
	tx := dao.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 3. 查询 slot
	var slot model.InterviewSlot
	if err := tx.Where("id = ?", sid).First(&slot).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("面试时段不存在")
		}
		return err
	}

	// 4. 检查报名是否已截止
	if util.NowUnix() > slot.StartTime.Unix() {
		tx.Rollback()
		return errors.New("报名时间已经截止")
	}

	// 5. 先判断是否已满
	if slot.Num >= slot.MaxNum {
		tx.Rollback()
		return errors.New("该时段已报满")
	}

	// 6. 更新名额
	slot.Num++
	if err := tx.Save(&slot).Error; err != nil {
		tx.Rollback()
		return err
	}
	// 7. 检查用户有几次面试
	var existingCount int64
	if err := tx.Model(&model.InterviewAssignment{}).
		Where("user_id = ?", uid).
		Count(&existingCount).Error; err != nil {
		tx.Rollback()
		return err
	}
	if existingCount > 0 {
		tx.Rollback()
		return errors.New("你已预约过面试，请勿重复提交")
	}
	// 8. 创建预约记录
	assignment := model.InterviewAssignment{
		UserID: uid,
		SlotID: sid,
		Round:  slot.Round,
	}
	if err := tx.Create(&assignment).Error; err != nil {
		tx.Rollback()
		// 处理重复预约
		if strings.Contains(err.Error(), "slot_user") ||
			strings.Contains(err.Error(), "Duplicate entry") ||
			strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return errors.New("你已预约过该时段，请勿重复提交")
		}
		return err
	}

	// 9. 提交事务
	return tx.Commit().Error
}
func (dao *UserDao) GetAdmins() ([]model.AdminModel, error) {
	var admins []model.AdminModel
	err := dao.db.Select("name", "id", "phone", "direction").Find(&admins).Error
	return admins, err
}
func (dao *UserDao) Delete(uid uint) error {
	// 开启事务
	tx := dao.db.Begin()
	// 防止panic导致连接丢失
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var ia model.InterviewAssignment
	if err := tx.Where("user_id = ?", uid).First(&ia).Error; err != nil {
		tx.Rollback()
		return errors.New("未报名，无法进行该操作")
	}
	// 物理删除
	if err := tx.Unscoped().Delete(&ia).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 减少 num，防止减成负数
	res := tx.Model(&model.InterviewSlot{}).
		Where("id = ? AND num > 0", ia.SlotID).
		Update("num", gorm.Expr("num - 1"))
	if res.Error != nil {
		tx.Rollback()
		return res.Error
	}
	if res.RowsAffected == 0 {
		tx.Rollback()
		return errors.New("slot 人数已为0，无法再减少")
	}

	return tx.Commit().Error
}
func (dao *UserDao) Update(uid uint, newSlotID uint, name string, direction int) error {

	if uid == 0 || newSlotID == 0 {
		return errors.New("无效的用户ID或时段ID")
	}
	tx := dao.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var ia model.InterviewAssignment
	if err := tx.Where("user_id = ?", uid).First(&ia).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("未找到您的预约记录: %w", err)
	}

	oldSlotID := ia.SlotID

	// 1. 旧 slot 减 1
	resDec := tx.Model(&model.InterviewSlot{}).
		Where("id = ? AND num > 0", oldSlotID).
		Update("num", gorm.Expr("num - 1"))
	if resDec.Error != nil {
		tx.Rollback()
		return fmt.Errorf("释放原时段名额失败: %w", resDec.Error)
	}

	if resDec.RowsAffected == 0 {
		tx.Rollback()
		return errors.New("原面试时段人数异常，无法释放名额")
	}
	// 2. 新 slot 加 1
	resInc := tx.Model(&model.InterviewSlot{}).
		Where("id = ? AND num < max_num", newSlotID).
		Update("num", gorm.Expr("num + 1"))
	if resInc.Error != nil {
		tx.Model(&model.InterviewSlot{}).Where("id = ?", oldSlotID).Update("num", gorm.Expr("num + 1"))
		tx.Rollback()
		return fmt.Errorf("目标时段已满或不可用: %w", resInc.Error)
	}

	if resInc.RowsAffected == 0 {
		// 同样需要回滚旧 slot
		tx.Model(&model.InterviewSlot{}).Where("id = ?", oldSlotID).Update("num", gorm.Expr("num + 1"))
		tx.Rollback()
		return errors.New("目标面试时段已满，无法更换")
	}

	// 3. 更新 assignment 的 slot_id
	if err := tx.Model(&ia).Updates(model.InterviewAssignment{
		SlotID:    newSlotID,
		Direction: direction,
	}).Error; err != nil {
		tx.Model(&model.InterviewSlot{}).Where("id = ?", oldSlotID).Update("num", gorm.Expr("num + 1"))
		tx.Model(&model.InterviewSlot{}).Where("id = ?", newSlotID).Update("num", gorm.Expr("num - 1"))
		tx.Rollback()
		return fmt.Errorf("更新预约记录失败: %w", err)
	}

	return tx.Commit().Error
}

// GetUserCurrentValidSlotID 返回用户当前唯一有效的面试时段 ID

func (dao *UserDao) GetUserCurrentValidSlotID(uid uint) (uint, error) {
	var assignments []model.InterviewAssignment

	// 获取用户所有面试预约
	if err := dao.db.Where("user_id = ?", uid).Find(&assignments).Error; err != nil {
		return 0, err
	}

	if len(assignments) == 0 {
		return 0, gorm.ErrRecordNotFound
	}

	now := util.NowUnix()
	slotDao := NewSlotDao(dao.db)
	var validSlotID uint
	validCount := 0

	for _, a := range assignments {
		slot, err := slotDao.GetSlotById(a.SlotID)
		if err != nil || slot == nil {
			continue // 跳过无效关联
		}
		// 判断是否未截至
		if now < slot.EndTime.Unix() {
			validSlotID = slot.ID
			validCount++
		}
	}
	if validCount == 0 {
		return 0, errors.New("当前无有效的面试预约")
	}
	fmt.Println("dao : ", validSlotID)
	return validSlotID, nil
}
func (dao *UserDao) UpdatePass(uid uint, round int, pass int, adminId uint) error {
	var err error

	// 开启事务
	err = dao.db.Transaction(func(tx *gorm.DB) error {
		// 表示事务内部错误
		var innerErr error
		field := ""
		if round == 1 {
			field = "first_pass"
		} else if round == 2 {
			field = "second_pass"
		}

		//  更新用户状态
		result := tx.Model(&model.UserModel{}).
			Where("id = ?", uid).
			Update(field, pass)

		if result.Error != nil {
			innerErr = result.Error
			return innerErr
		}
		if result.RowsAffected == 0 {
			innerErr = gorm.ErrRecordNotFound
			return innerErr
		}

		// 如果通过，创建面试结果记录
		if pass == 1 {
			interviewResult := model.InterviewResult{
				UserID:  uid,
				Round:   round,
				Status:  1,
				AdminID: adminId,
			}
			innerErr = tx.Create(&interviewResult).Error
			if innerErr != nil {
				return innerErr
			}
		}

		return nil // 事务成功
	})

	return err
}
func (dao *UserDao) GetUsersByRound(round int) ([]model.UserResponse, error) {
	var users []model.UserModel
	var err error
	var usersRes []model.UserResponse
	if round == 1 {
		err = dao.db.Where("first_pass = ?", 1).Find(&users).Error
	} else {
		err = dao.db.Where("second_pass = ?", 1).Find(&users).Error
	}
	if err != nil {
		return nil, err
	}
	for _, user := range users {
		userRes := model.UserResponse{
			Name:  user.Name,
			Email: user.Email,
			ID:    user.ID,
		}
		usersRes = append(usersRes, userRes)
	}
	return usersRes, nil
}
