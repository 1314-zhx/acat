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
	db *gorm.DB // 只保留 db，不存 ctx
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
	err := dao.db.WithContext(ctx).Select("id", "name", "password", "phone").
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

	// 7. 创建预约记录
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

	// 8. 提交事务
	return tx.Commit().Error
}
