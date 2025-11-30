package dao

import (
	"acat/model"
	"gorm.io/gorm"
)

type SlotDao struct {
	db *gorm.DB
}

func NewSlotDao(db *gorm.DB) *SlotDao {
	return &SlotDao{db: db}
}
func (dao *SlotDao) GetSlotById(sid uint) (*model.InterviewSlot, error) {
	var slot model.InterviewSlot
	err := dao.db.Where("id = ? ", sid).First(&slot).Error
	return &slot, err
}
