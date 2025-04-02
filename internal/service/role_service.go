package service

import (
	"authentication/internal/model"
	"authentication/internal/repository"
	"errors"
	"fmt"
)

// RoleService 角色服务接口
type RoleService interface {
	Create(role *model.Role) error
	GetByID(id uint) (*model.Role, error)
	GetByName(name string) (*model.Role, error)
	Update(role *model.Role) error
	Delete(id uint) error
	List(page, pageSize int) ([]model.Role, int64, error)
	AssignPermissions(roleID uint, permissionIDs []uint) error
}

// roleService 角色服务实现
type roleService struct {
	roleRepo       repository.RoleRepository
	permissionRepo repository.PermissionRepository
}

// NewRoleService 创建角色服务实例
func NewRoleService(roleRepo repository.RoleRepository, permissionRepo repository.PermissionRepository) RoleService {
	return &roleService{
		roleRepo:       roleRepo,
		permissionRepo: permissionRepo,
	}
}

// Create 创建角色
func (s *roleService) Create(role *model.Role) error {
	// 检查角色名是否已存在
	_, err := s.roleRepo.GetByName(role.Name)
	if err == nil {
		return errors.New("角色名已存在")
	}

	// 创建角色
	return s.roleRepo.Create(role)
}

// GetByID 根据ID获取角色
func (s *roleService) GetByID(id uint) (*model.Role, error) {
	return s.roleRepo.GetByID(id)
}

// GetByName 根据名称获取角色
func (s *roleService) GetByName(name string) (*model.Role, error) {
	return s.roleRepo.GetByName(name)
}

// Update 更新角色
func (s *roleService) Update(role *model.Role) error {
	// 检查角色是否存在
	existingRole, err := s.roleRepo.GetByID(role.ID)
	if err != nil {
		return fmt.Errorf("角色不存在: %w", err)
	}

	// 如果角色名已更改，检查新名称是否已存在
	if existingRole.Name != role.Name {
		_, err := s.roleRepo.GetByName(role.Name)
		if err == nil {
			return errors.New("角色名已存在")
		}
	}

	// 更新角色
	return s.roleRepo.Update(role)
}

// Delete 删除角色
func (s *roleService) Delete(id uint) error {
	// 检查角色是否存在
	_, err := s.roleRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("角色不存在: %w", err)
	}

	// 删除角色
	return s.roleRepo.Delete(id)
}

// List 获取角色列表
func (s *roleService) List(page, pageSize int) ([]model.Role, int64, error) {
	return s.roleRepo.List(page, pageSize)
}

// AssignPermissions 分配权限到角色
func (s *roleService) AssignPermissions(roleID uint, permissionIDs []uint) error {
	// 检查角色是否存在
	_, err := s.roleRepo.GetByID(roleID)
	if err != nil {
		return fmt.Errorf("角色不存在: %w", err)
	}

	// 检查所有权限是否存在
	for _, permID := range permissionIDs {
		_, err := s.permissionRepo.GetByID(permID)
		if err != nil {
			return fmt.Errorf("权限ID %d 不存在: %w", permID, err)
		}
	}

	// 分配权限
	return s.roleRepo.AssignPermissions(roleID, permissionIDs)
}
