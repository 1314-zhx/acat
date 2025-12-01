package model

import (
	"gorm.io/gorm"
	"time"
)

// UserModel 用户通用模型 对应user_models表
type UserModel struct {
	ID         uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	StuId      string         `json:"stu_id" gorm:"type:varchar(32);uniqueIndex;not null"` // 修改这里
	Name       string         `json:"name" gorm:"type:varchar(64);not null"`
	Password   string         `json:"-" gorm:"type:varchar(255);not null"`                // bcrypt hash 约 60 字符
	Phone      string         `json:"phone" gorm:"type:varchar(32);uniqueIndex;not null"` //  修改这里
	Email      string         `json:"email" gorm:"type:varchar(128);not null"`
	FirstPass  int            `json:"first_pass" gorm:"not null;default:0"`
	SecondPass int            `json:"second_pass" gorm:"not null;default:0"`
	Direction  int            `json:"direction" gorm:"not null"`
	CreatedAt  time.Time      `json:"create_time" gorm:"autoCreateTime"`
	UpdatedAt  time.Time      `json:"update_time" gorm:"autoUpdateTime"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"softDelete"`
	Gender     int            `json:"gender" gorm:"default:1"`
}

// UserLogin 用户登录模型，用于整合UserLog和UserModel，UserLog太小信息不全，UserModel太大了并且直接映射数据库字段，需要一个中量且安全的结构
type UserLogin struct {
	UId      uint   `json:"uid,omitempty"`    // 用户ID
	Name     string `json:"name,omitempty"`   // 名字
	Password string `json:"password"`         // 密码
	Phone    string `json:"phone"`            // 手机号
	VCode    string `json:"v_code,omitempty"` // 验证码
}

// UserRegister 用户注册模型 用于整合UserRes和UserModel
type UserRegister struct {
	UId        uint      `json:"uid,omitempty"`
	Name       string    `json:"name" binding:"required"`           // 姓名
	StuId      string    `json:"stu_id" binding:"required"`         // 学号
	Password   string    `json:"password" binding:"required,min=6"` // 密码
	Phone      string    `json:"phone" binding:"required"`          // 手机号
	Email      string    `json:"email" binding:"required,email"`    // QQ邮箱
	VCode      string    `json:"v_code" binding:"required"`         // 短信验证码
	Gender     int       `json:"gender" binding:"required"`
	Direction  int       `json:"direction"` // 方向选择，0不确定，1为Go，2为Java，3为前端，4为后端
	CreateTime time.Time `json:"create_time"`
}

// UserResponse 是脱敏后的用户信息摘要，专用于向前端返回。
type UserResponse struct {
	ID     uint   `json:"id"`              // 用户唯一 ID
	Name   string `json:"name"`            // 用户姓名
	StuId  string `json:"stu_id"`          // 学号
	Phone  string `json:"phone,omitempty"` // 手机号
	Gender int    `json:"gender"`          // 性别 1 女，2 男
}
