package dao

import (
	"acat/model"
	"acat/util/encryption"
	"gorm.io/gorm"
)

type AdminDao struct {
	db *gorm.DB
}

func NewAdminDao(db *gorm.DB) *AdminDao {
	return &AdminDao{db: db}
}

// 统一返回错误 nil 不暴露底层错误原因
func (dao *AdminDao) Login(phone, password string) (bool, error) {
	var admin model.AdminModel
	err := dao.db.Where("phone = ?", phone).First(&admin).Error
	if err != nil {
		return false, nil
	}
	if !encryption.CheckPassword(password, admin.Password) {
		return false, nil
	}
	return true, nil
}
func (dao *AdminDao) GetAdmin(phone string) (*model.AdminView, error) {
	var admin model.AdminModel
	err := dao.db.Where("phone = ?", phone).First(&admin).Error
	// 防止空指针，虽然错误上交给logic处理，但不能返回零值空指针，而应返回nil空指针
	if err != nil {
		return nil, err
	}
	var adminView model.AdminView
	adminView.Aid = admin.ID
	adminView.Name = admin.Name
	adminView.Phone = admin.Phone
	adminView.Direction = admin.Direction
	return &adminView, nil
}
