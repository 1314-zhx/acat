package dao

import (
	"acat/model"
	"fmt"
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
	fmt.Println("i am here", slot)
	return &slot, err
}
func (dao *SlotDao) Create(slot *model.InterviewSlot) error {
	fmt.Println("dao")
	return dao.db.Create(slot).Error
}
func (dao *SlotDao) GetUserIdsBySlotId(slotId uint) ([]uint, error) {
	var slotAssignments []model.InterviewAssignment
	err := dao.db.Select("user_id").Where("slot_id = ?", slotId).Find(&slotAssignments).Error
	var userIds []uint
	for _, slot := range slotAssignments {
		userIds = append(userIds, slot.UserID)
	}
	return userIds, err
}
