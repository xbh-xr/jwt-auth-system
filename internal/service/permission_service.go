package service

import (
	"authentication/internal/model"
	"authentication/internal/repository"
	"errors"
	"fmt"
)

// PermissionService 权限服务接口
type PermissionService interface {
	Create(permission *model.Permission) error
	GetByID(id uint) (*model.Permission, error)
	GetByCode(code string) (*model.Permission, error)
	Update(permission *model.Permission) error
	Delete(id uint) error
	List(page, pageSize int) ([]model.Permission, int64, error)
}

// permissionService 权限服务实现
type permissionService struct {
	permissionRepo repository.PermissionRepository
}

// NewPermissionService 创建权限服务实例
func NewPermissionService(permissionRepo repository.PermissionRepository) PermissionService {
	return &permissionService{
		permissionRepo: permissionRepo,
	}
}

// Create 创建权限
func (s *permissionService) Create(permission *model.Permission) error {
	// 检查权限代码是否已存在
	_, err := s.permissionRepo.GetByCode(permission.Code)
	if err == nil {
		return errors.New("权限代码已存在")
	}

	// 创建权限
	return s.permissionRepo.Create(permission)
}

// GetByID 根据ID获取权限
func (s *permissionService) GetByID(id uint) (*model.Permission, error) {
	return s.permissionRepo.GetByID(id)
}

// GetByCode 根据代码获取权限
func (s *permissionService) GetByCode(code string) (*model.Permission, error) {
	return s.permissionRepo.GetByCode(code)
}

// Update 更新权限
func (s *permissionService) Update(permission *model.Permission) error {
	// 检查权限是否存在
	existingPermission, err := s.permissionRepo.GetByID(permission.ID)
	if err != nil {
		return fmt.Errorf("权限不存在: %w", err)
	}

	// 如果权限代码已更改，检查新代码是否已存在
	if existingPermission.Code != permission.Code {
		_, err := s.permissionRepo.GetByCode(permission.Code)
		if err == nil {
			return errors.New("权限代码已存在")
		}
	}

	// 更新权限
	return s.permissionRepo.Update(permission)
}

// Delete 删除权限
func (s *permissionService) Delete(id uint) error {
	// 检查权限是否存在
	_, err := s.permissionRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("权限不存在: %w", err)
	}

	// 删除权限
	return s.permissionRepo.Delete(id)
}

// List 获取权限列表
func (s *permissionService) List(page, pageSize int) ([]model.Permission, int64, error) {
	return s.permissionRepo.List(page, pageSize)
}
