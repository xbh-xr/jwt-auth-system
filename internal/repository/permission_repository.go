package repository

import (
	"authentication/internal/model"
	"gorm.io/gorm"
)

// PermissionRepository 权限存储库接口
type PermissionRepository interface {
	Create(permission *model.Permission) error
	GetByID(id uint) (*model.Permission, error)
	GetByCode(code string) (*model.Permission, error)
	Update(permission *model.Permission) error
	Delete(id uint) error
	List(page, pageSize int) ([]model.Permission, int64, error)
}

// permissionRepository 权限存储库实现
type permissionRepository struct {
	db *gorm.DB
}

// NewPermissionRepository 创建权限存储库实例
func NewPermissionRepository(db *gorm.DB) PermissionRepository {
	return &permissionRepository{db: db}
}

// Create 创建权限
func (r *permissionRepository) Create(permission *model.Permission) error {
	return r.db.Create(permission).Error
}

// GetByID 根据ID获取权限
func (r *permissionRepository) GetByID(id uint) (*model.Permission, error) {
	var permission model.Permission
	err := r.db.First(&permission, id).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

// GetByCode 根据代码获取权限
func (r *permissionRepository) GetByCode(code string) (*model.Permission, error) {
	var permission model.Permission
	err := r.db.Where("code = ?", code).First(&permission).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

// Update 更新权限
func (r *permissionRepository) Update(permission *model.Permission) error {
	return r.db.Save(permission).Error
}

// Delete 删除权限
func (r *permissionRepository) Delete(id uint) error {
	return r.db.Delete(&model.Permission{}, id).Error
}

// List 获取权限列表
func (r *permissionRepository) List(page, pageSize int) ([]model.Permission, int64, error) {
	var permissions []model.Permission
	var total int64

	// 计算总数
	if err := r.db.Model(&model.Permission{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := r.db.Offset(offset).Limit(pageSize).Find(&permissions).Error
	if err != nil {
		return nil, 0, err
	}

	return permissions, total, nil
}
