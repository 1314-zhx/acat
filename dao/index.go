package dao

import (
	"acat/model"
	"context"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

type IndexDao struct {
	db *gorm.DB
}

func NewIndexDao(db *gorm.DB) *IndexDao {
	return &IndexDao{db: db}
}

func (d *IndexDao) Check(ctx context.Context) ([][]int64, error) {
	var slots []model.InterviewSlot

	now := time.Now()
	err := d.db.WithContext(ctx).
		Select("id", "round", "start_time", "end_time", "num", "max_num").
		Where("start_time > ?", now). // 只查未开始的
		Order("start_time ASC").
		Find(&slots).Error

	if err != nil {
		zap.L().Error("dao/index.go Check() failed", zap.Error(err))
		return nil, err
	}

	schedule := make([][]int64, len(slots))
	for i, slot := range slots {
		schedule[i] = []int64{
			slot.StartTime.Unix(), // [0] start
			slot.EndTime.Unix(),   // [1] end
			int64(slot.ID),        // [2] id
			int64(slot.Round),     // [3] round
			int64(slot.Num),       // [4] num
			int64(slot.MaxNum),    // [5] max_num
		}
	}

	return schedule, nil
}
