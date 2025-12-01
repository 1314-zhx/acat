package dao

import (
	"acat/model"
	"gorm.io/gorm"
)

type LetterDao struct {
	db *gorm.DB
}

func NewLetterDao(db *gorm.DB) *LetterDao {
	return &LetterDao{db: db}
}
func (dao *LetterDao) Letter(title string, content string, receiveId uint, uid uint) error {
	message := model.Message{
		SendID:    uid,
		ReceiveID: receiveId,
		Title:     title,
		Content:   content,
		Type:      0, // 用户给管理员
	}
	return dao.db.Create(&message).Error
}
