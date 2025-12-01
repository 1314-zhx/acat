package model

import (
	"gorm.io/gorm"
	"time"
)

// InterviewSlot 表示一个面试时间槽，对应数据库中的 interview_slot 表。
// 每个时间槽属于某一轮次，并具有开始/结束时间。
type InterviewSlot struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Round     int       `json:"round" gorm:"not null;index"`
	StartTime time.Time `json:"start_time" gorm:"not null;index"`
	EndTime   time.Time `json:"end_time" gorm:"not null"`
	Num       int       `json:"num" gorm:"not null;default:0"`
	MaxNum    int       `json:"max_num" gorm:"not null;default:50"`
}

// InterviewAssignment 是用户与面试时间槽之间的多对多关联表（中间表）。
// 通过 UserId 和 SlotId 建立唯一约束，确保同一用户在同一轮次不会被重复分配。
// 对应interview_assignment表
type InterviewAssignment struct {
	ID        uint           `gorm:"primaryKey;autoIncrement"`
	UserID    uint           `gorm:"not null;uniqueIndex:slot_user"`
	SlotID    uint           `gorm:"not null;uniqueIndex:slot_user"`
	Round     int            `gorm:"not null;index"`
	Direction int            `gorm:"not null;default:0"` // 0不确定，1为Go，2为Java，3为前端，4为后端
	DeletedAt gorm.DeletedAt // 软删除，DeletedAt不能改
}

// ScheduleResponse 是返回给前端的面试日程响应结构。
// 用于在 API 响应中嵌套展示某轮次参与面试的所有候选人。
type ScheduleResponse struct {
	Round int                 `json:"round"`
	Slots []TimeSlotWithUsers `json:"slots"`
}

// TimeSlotWithUsers 表示一个具体的时间段及其关联的用户列表。
// 用于在 API 响应中嵌套展示某个时间段内参与面试的所有候选人。
type TimeSlotWithUsers struct {
	StartTime time.Time      `json:"start_time"` // 时间段开始时间戳
	EndTime   time.Time      `json:"end_time"`   // 时间段结束时间戳
	Users     []UserResponse `json:"users"`      // 该时间段内所有被安排的用户
}
