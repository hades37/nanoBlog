package dao

import (
	"time"

	"gorm.io/gorm"
)

// Admin 管理员模型
type Admin struct {
	ID        uint           `gorm:"primarykey"`                       // 主键ID
	Username  string         `gorm:"type:varchar(32);unique;not null"` // 用户名
	Password  string         `gorm:"type:varchar(128);not null"`       // 密码
	Nickname  string         `gorm:"type:varchar(32)"`                 // 昵称
	Avatar    string         `gorm:"type:varchar(255)"`                // 头像
	Email     string         `gorm:"type:varchar(64)"`                 // 邮箱
	Phone     string         `gorm:"type:varchar(20)"`                 // 手机号
	Status    uint8          `gorm:"type:tinyint;default:1"`           // 状态 1:启用 2:禁用
	LastLogin time.Time      `gorm:"type:datetime"`                    // 最后登录时间
	CreatedAt time.Time      `gorm:"type:datetime;not null"`           // 创建时间
	UpdatedAt time.Time      `gorm:"type:datetime;not null"`           // 更新时间
	DeletedAt gorm.DeletedAt `gorm:"index"`                            // 软删除
}

// TableName 设置表名
func (Admin) TableName() string {
	return "admin"
}

// ... existing Admin struct and TableName ...
// CreateAdmin 创建管理员
func CreateAdmin(db *gorm.DB, admin *Admin) error {
	return db.Create(admin).Error
}

// GetAdminByID 根据ID获取管理员信息
func GetAdminByID(db *gorm.DB, id uint) (*Admin, error) {
	var admin Admin
	err := db.Where("id = ?", id).First(&admin).Error
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

// GetAdminByUsername 根据用户名获取管理员信息
func GetAdminByUsername(db *gorm.DB, username string) (*Admin, error) {
	var admin Admin
	err := db.Where("username = ?", username).First(&admin).Error
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

// UpdateAdmin 更新管理员信息
func UpdateAdmin(db *gorm.DB, admin *Admin) error {
	return db.Save(admin).Error
}

// DeleteAdmin 删除管理员（软删除）
func DeleteAdmin(db *gorm.DB, id uint) error {
	return db.Delete(&Admin{}, id).Error
}

// ListAdmins 获取管理员列表
func ListAdmins(db *gorm.DB, page, pageSize int) ([]Admin, int64, error) {
	var admins []Admin
	var total int64

	// 获取总数
	err := db.Model(&Admin{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	offset := (page - 1) * pageSize
	err = db.Offset(offset).Limit(pageSize).Find(&admins).Error
	if err != nil {
		return nil, 0, err
	}

	return admins, total, nil
}

// UpdateAdminStatus 更新管理员状态
func UpdateAdminStatus(db *gorm.DB, id uint, status uint8) error {
	return db.Model(&Admin{}).Where("id = ?", id).Update("status", status).Error
}

// UpdateAdminPassword 更新管理员密码
func UpdateAdminPassword(db *gorm.DB, id uint, newPassword string) error {
	return db.Model(&Admin{}).Where("id = ?", id).Update("password", newPassword).Error
}

// UpdateLastLogin 更新最后登录时间
func UpdateLastLogin(db *gorm.DB, id uint) error {
	return db.Model(&Admin{}).Where("id = ?", id).Update("last_login", time.Now()).Error
}
