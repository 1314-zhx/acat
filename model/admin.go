package model

import (
	"gorm.io/gorm"
	"time"
)

// AdminModel 管理员模型（对应数据库 admin_model 表）
type AdminModel struct {
	ID        uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string         `json:"name" gorm:"not null"`
	Phone     string         `json:"phone,omitempty" gorm:"not null"`
	Password  string         `json:"-" gorm:"not null"`
	Level     int            `json:"level" gorm:"not null;default:1"` // 1=普通, 2=超级, 3=只读
	Email     string         `json:"email,omitempty" gorm:"not null"`
	CreatedAt time.Time      `json:"create_time" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"update_time" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"-"`
}

// AdminLogin 管理员登录请求模型
type AdminLogin struct {
	Aid      int    `json:"aid,omitempty"`               // 登录成功后可选返回
	Name     string `json:"name" binding:"required"`     // 姓名
	Password string `json:"password" binding:"required"` // 密码
}
