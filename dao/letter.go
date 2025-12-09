package dao

import (
	"acat/model"
	"context"
	"fmt"
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
func (u *UserDao) GetUsersByIDs(ids []uint, ctx context.Context) ([]model.UserModel, error) {
	var users []model.UserModel
	err := u.db.WithContext(ctx).Where("id IN ?", ids).Find(&users).Error
	return users, err
}
func (dao *LetterDao) GetAdminLetters(aid uint) ([]model.MessageResponse, error) {
	var letters []model.Message
	err := dao.db.Select("id, send_id, receive_id, title, content,is_read").
		Where("receive_id = ?  AND type = 0", aid).
		Find(&letters).Error
	if err != nil {
		return nil, err
	}

	if len(letters) == 0 {
		return []model.MessageResponse{}, nil
	}

	// 1. 收集所有 SendID ,去重,保证每个用户只会被查一次，即使有多封信
	sendIDs := make(map[uint]bool)
	for _, letter := range letters {
		sendIDs[letter.SendID] = true
	}

	// 2. 转为 slice 用于 IN 查询
	ids := make([]uint, 0, len(sendIDs))
	for id := range sendIDs {
		ids = append(ids, id)
	}

	// 3. 批量查询用户
	userDao := NewUserDao(dao.db)
	users, err := userDao.GetUsersByIDs(ids, context.Background())
	if err != nil {
		return nil, err
	}

	// 4. 构建 ID -> User 映射
	userMap := make(map[uint]model.UserModel, len(users))
	for _, u := range users {
		userMap[u.ID] = u
	}

	// 5. 组装响应
	lettersResponse := make([]model.MessageResponse, len(letters))
	for i, letter := range letters {
		user := userMap[letter.SendID]
		lettersResponse[i] = model.MessageResponse{
			ID:        letter.ID,
			SendID:    letter.SendID,
			SendName:  user.Name, // 注意：如果用户被删除，user.Name 会是空字符串
			ReceiveID: letter.ReceiveID,
			Title:     letter.Title,
			Content:   letter.Content,
			IsRead:    letter.IsRead,
		}
	}
	return lettersResponse, nil
}
func (dao *LetterDao) SetIsRead(lid int64, isRead bool) error {
	return dao.db.Model(&model.Message{}).
		Where("id = ?", lid).
		Update("is_read", isRead).Error
}
func (dao *LetterDao) Reply(aid uint, uid uint, title string, content string) error {
	message := model.Message{
		SendID:    aid,
		ReceiveID: uid,
		Title:     title,
		Content:   content,
		IsRead:    false,
		Type:      1,
	}
	fmt.Println(message)
	err := dao.db.Create(&message).Error
	return err
}
