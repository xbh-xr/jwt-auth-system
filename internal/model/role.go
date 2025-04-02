package model

import (
	"time"
)

// Role 角色模型
type Role struct {
	ID          uint         `json:"id" gorm:"primaryKey"`
	Name        string       `json:"name" gorm:"size:50;uniqueIndex;not null"`
	Description string       `json:"description" gorm:"size:200"`
	Permissions []Permission `json:"permissions" gorm:"many2many:role_permissions;"`
	Users       []User       `json:"users,omitempty" gorm:"many2many:user_roles;"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

// HasPermission 检查角色是否拥有指定权限
func (r *Role) HasPermission(permissionCode string) bool {
	for _, permission := range r.Permissions {
		if permission.Code == permissionCode {
			return true
		}
	}
	return false
}
