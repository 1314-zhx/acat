package model

import "time"

// Message 用户与管理员之间的通信消息
type Message struct {
	ID        int64     `json:"id" xorm:"pk autoincr"`              // 主键
	SendId    int       `json:"send_id" xorm:"not null"`            // 发送者ID
	ReceiveId int       `json:"receive_id" xorm:"not null"`         // 接收者ID
	Title     string    `json:"title" xorm:"varchar(100) not null"` // 标题
	Content   string    `json:"content" xorm:"text not null"`       // 消息正文（关键！）
	CreatedAt time.Time `json:"created_at" xorm:"created"`          // 自动记录创建时间
	IsRead    bool      `json:"is_read" xorm:"default false"`       // 是否已读
}
