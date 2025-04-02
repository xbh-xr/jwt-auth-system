package repository

import (
	"authentication/internal/config"
	"authentication/internal/model"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// InitDB 初始化数据库连接
func InitDB(cfg config.DBConfig) (*gorm.DB, error) {
	// 构建DSN
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName)

	// 连接数据库
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("连接数据库失败: %w", err)
	}

	// 自动迁移模型
	err = db.AutoMigrate(
		&model.User{},
		&model.Role{},
		&model.Permission{},
	)
	if err != nil {
		return nil, fmt.Errorf("迁移数据库模型失败: %w", err)
	}

	// 初始化基础数据
	if err := initBaseData(db); err != nil {
		return nil, fmt.Errorf("初始化基础数据失败: %w", err)
	}

	return db, nil
}

// initBaseData 初始化基础数据
func initBaseData(db *gorm.DB) error {
	// 检查是否已有权限数据
	var count int64
	db.Model(&model.Permission{}).Count(&count)
	if count > 0 {
		// 已有数据，不需要初始化
		return nil
	}

	// 创建基础权限
	permissions := []model.Permission{
		{Code: "user:list", Name: "用户列表", Description: "查看用户列表"},
		{Code: "user:read", Name: "查看用户", Description: "查看用户详情"},
		{Code: "user:create", Name: "创建用户", Description: "创建新用户"},
		{Code: "user:update", Name: "更新用户", Description: "更新用户信息"},
		{Code: "user:delete", Name: "删除用户", Description: "删除用户"},

		{Code: "role:list", Name: "角色列表", Description: "查看角色列表"},
		{Code: "role:read", Name: "查看角色", Description: "查看角色详情"},
		{Code: "role:create", Name: "创建角色", Description: "创建新角色"},
		{Code: "role:update", Name: "更新角色", Description: "更新角色信息"},
		{Code: "role:delete", Name: "删除角色", Description: "删除角色"},
		{Code: "role:assign", Name: "分配权限", Description: "为角色分配权限"},

		{Code: "permission:list", Name: "权限列表", Description: "查看权限列表"},
	}

	// 创建基础角色
	adminRole := model.Role{
		Name:        "admin",
		Description: "系统管理员",
		Permissions: permissions,
	}

	userRole := model.Role{
		Name:        "user",
		Description: "普通用户",
		Permissions: []model.Permission{
			permissions[0], // user:list
			permissions[1], // user:read
		},
	}

	// 创建管理员用户
	adminUser := model.User{
		Username: "admin",
		Email:    "admin@example.com",
		Password: "password", // 会在BeforeSave钩子中自动加密
		FullName: "系统管理员",
		Active:   true,
		Roles:    []model.Role{adminRole},
	}

	// 使用事务保证数据一致性
	return db.Transaction(func(tx *gorm.DB) error {
		// 创建权限
		for _, perm := range permissions {
			if err := tx.Create(&perm).Error; err != nil {
				return err
			}
		}

		// 创建角色
		if err := tx.Create(&adminRole).Error; err != nil {
			return err
		}

		if err := tx.Create(&userRole).Error; err != nil {
			return err
		}

		// 创建管理员用户
		if err := tx.Create(&adminUser).Error; err != nil {
			return err
		}

		return nil
	})
}
