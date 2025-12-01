package dao

import (
	"acat/model"
	"context"
	"fmt"
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
			slot.StartTime.Unix(),
			slot.EndTime.Unix(),
			int64(slot.ID),
			int64(slot.Round),
			int64(slot.Num),
			int64(slot.MaxNum),
		}
	}
	fmt.Println("hhh ", schedule)
	return schedule, nil
}
