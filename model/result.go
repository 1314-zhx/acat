package model

import "time"

// InterviewResult 面试结果记录
type InterviewResult struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	UserID    uint      `gorm:"not null;index"` // 关联用户（注意：GORM 外键默认是 UserID）
	Round     int       `gorm:"not null;index"`
	Status    int       `gorm:"not null;default:0"` // 0=未面试, 1=通过, 2=不通过, 3=待定
	Comment   string    `gorm:"type:text"`
	AdminID   uint      `gorm:"index"` // 评分管理员ID
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

// UserResultToUser 用于向前端返回脱敏后的用户结果，给学生展示的
type UserResultToUser struct {
	Name  string `json:"name"`
	StuId string `json:"stu_id"`
	Pass  bool   `json:"pass"`
	Round int    `json:"round"`
}

// UserResultToAdmin 用于给前端展示脱敏后的用户结果，给管理员展示
type UserResultToAdmin struct {
	Name  string `json:"name"`
	StuId string `json:"stu_id"`
	Pass  bool   `json:"pass"`
	Round int    `json:"round"`
	Phone string `json:"phone"`
	Email string `json:"email"`
}
