package service

import (
	"authentication/internal/model"
	"authentication/internal/repository"
)

// UserService 用户服务接口
type UserService interface {
	GetByID(id uint) (*model.User, error)
	List(page, pageSize int) ([]model.User, int64, error)
	Update(user *model.User) error
	Delete(id uint) error
}

// userService 用户服务实现
type userService struct {
	userRepo repository.UserRepository
}

// NewUserService 创建用户服务实例
func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

// GetByID 根据ID获取用户
func (s *userService) GetByID(id uint) (*model.User, error) {
	return s.userRepo.GetByID(id)
}

// List 获取用户列表
func (s *userService) List(page, pageSize int) ([]model.User, int64, error) {
	return s.userRepo.List(page, pageSize)
}

// Update 更新用户
func (s *userService) Update(user *model.User) error {
	return s.userRepo.Update(user)
}

// Delete 删除用户
func (s *userService) Delete(id uint) error {
	return s.userRepo.Delete(id)
}
