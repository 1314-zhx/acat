package model

import (
	"time"
)

// Message 用户与管理员之间的通信消息
type Message struct {
	ID        int64     `json:"id" gorm:"primaryKey;autoIncrement"`      // 主键
	SendID    uint      `json:"send_id" gorm:"not null"`                 // 发送者ID
	ReceiveID uint      `json:"receive_id" gorm:"not null"`              // 接收者ID
	Title     string    `json:"title" gorm:"type:varchar(100);not null"` // 标题
	Content   string    `json:"content" gorm:"type:text;not null"`       // 消息正文
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`        // 自动记录创建时间
	IsRead    bool      `json:"is_read" gorm:"default:false"`            // 是否已读
	Type      int       `json:"type" gorm:"not null"`                    // 消息类型（如 0 用户→管理员, 1 管理员→用户）
}
