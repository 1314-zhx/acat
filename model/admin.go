package model

import (
	"gorm.io/gorm"
	"time"
)

// AdminModel 管理员模型（对应数据库 admin_model 表）
type AdminModel struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	Name      string    `gorm:"not null"`
	Phone     string    `gorm:"not null"`
	Password  string    `gorm:"not null"`
	Level     int       `gorm:"not null;default:1"` // 1=普通, 2=超级, 3=只读
	Email     string    `gorm:"not null"`
	Direction int       `gorm:"not null"` // 1 go 2 java 3 前端
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt
}

// AdminLogin 管理员登录请求模型
type AdminLogin struct {
	Aid      int    `json:"aid,omitempty"`               // 登录成功后可选返回
	Name     string `json:"name" binding:"required"`     // 姓名
	Password string `json:"password" binding:"required"` // 密码
}

// 想用户展示的管理员脱敏后的模型
type AdminView struct {
	Aid       uint   `json:"aid"`
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	Direction int    `json:"direction"`
}
