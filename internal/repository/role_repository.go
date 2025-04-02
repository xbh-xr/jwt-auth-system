package repository

import (
	"authentication/internal/model"
	"gorm.io/gorm"
)

// RoleRepository 角色存储库接口
type RoleRepository interface {
	Create(role *model.Role) error
	GetByID(id uint) (*model.Role, error)
	GetByName(name string) (*model.Role, error)
	Update(role *model.Role) error
	Delete(id uint) error
	List(page, pageSize int) ([]model.Role, int64, error)
	AssignPermissions(roleID uint, permissionIDs []uint) error
}

// roleRepository 角色存储库实现
type roleRepository struct {
	db *gorm.DB
}

// NewRoleRepository 创建角色存储库实例
func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{db: db}
}

// Create 创建角色
func (r *roleRepository) Create(role *model.Role) error {
	return r.db.Create(role).Error
}

// GetByID 根据ID获取角色
func (r *roleRepository) GetByID(id uint) (*model.Role, error) {
	var role model.Role
	err := r.db.Preload("Permissions").First(&role, id).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// GetByName 根据名称获取角色
func (r *roleRepository) GetByName(name string) (*model.Role, error) {
	var role model.Role
	err := r.db.Preload("Permissions").Where("name = ?", name).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// Update 更新角色
func (r *roleRepository) Update(role *model.Role) error {
	return r.db.Save(role).Error
}

// Delete 删除角色
func (r *roleRepository) Delete(id uint) error {
	return r.db.Delete(&model.Role{}, id).Error
}

// List 获取角色列表
func (r *roleRepository) List(page, pageSize int) ([]model.Role, int64, error) {
	var roles []model.Role
	var total int64

	// 计算总数
	if err := r.db.Model(&model.Role{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := r.db.Preload("Permissions").Offset(offset).Limit(pageSize).Find(&roles).Error
	if err != nil {
		return nil, 0, err
	}

	return roles, total, nil
}

// AssignPermissions 分配权限到角色
func (r *roleRepository) AssignPermissions(roleID uint, permissionIDs []uint) error {
	// 开始事务
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 获取角色
		role := model.Role{}
		if err := tx.First(&role, roleID).Error; err != nil {
			return err
		}

		// 清除现有权限关联
		if err := tx.Model(&role).Association("Permissions").Clear(); err != nil {
			return err
		}

		// 如果没有新的权限，直接返回
		if len(permissionIDs) == 0 {
			return nil
		}

		// 添加新的权限关联
		var permissions []model.Permission
		if err := tx.Find(&permissions, permissionIDs).Error; err != nil {
			return err
		}

		return tx.Model(&role).Association("Permissions").Append(permissions)
	})
}
