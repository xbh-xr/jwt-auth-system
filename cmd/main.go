package main

import (
	"authentication/internal/config"
	"authentication/internal/handler"
	"authentication/internal/middleware"
	"authentication/internal/repository"
	"authentication/internal/service"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 设置Gin模式
	gin.SetMode(cfg.Server.Mode)

	// 初始化数据库连接
	db, err := repository.InitDB(cfg.DB)
	if err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}

	// 初始化存储库
	userRepo := repository.NewUserRepository(db)
	roleRepo := repository.NewRoleRepository(db)
	permissionRepo := repository.NewPermissionRepository(db)

	// 初始化服务
	authService := service.NewAuthService(userRepo, cfg.JWT)
	userService := service.NewUserService(userRepo)
	roleService := service.NewRoleService(roleRepo, permissionRepo)
	permissionService := service.NewPermissionService(permissionRepo)

	// 初始化处理器
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	roleHandler := handler.NewRoleHandler(roleService)
	permissionHandler := handler.NewPermissionHandler(permissionService)

	// 创建路由
	r := gin.Default()

	// 注册中间件
	authMiddleware := middleware.NewAuthMiddleware(cfg.JWT)

	// API路由
	api := r.Group("/api")
	{
		// 认证路由 - 不需要认证
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)

			// 需要认证的路由
			auth.GET("/profile", authMiddleware.AuthRequired(), authHandler.GetProfile)
		}

		// 用户管理 - 需要认证
		users := api.Group("/users", authMiddleware.AuthRequired())
		{
			users.GET("", authMiddleware.HasPermission("user:list"), userHandler.ListUsers)
			users.GET("/:id", authMiddleware.HasPermission("user:read"), userHandler.GetUser)
			users.PUT("/:id", authMiddleware.HasPermission("user:update"), userHandler.UpdateUser)
			users.DELETE("/:id", authMiddleware.HasPermission("user:delete"), userHandler.DeleteUser)
		}

		// 角色管理 - 需要认证
		roles := api.Group("/roles", authMiddleware.AuthRequired())
		{
			roles.GET("", authMiddleware.HasPermission("role:list"), roleHandler.ListRoles)
			roles.POST("", authMiddleware.HasPermission("role:create"), roleHandler.CreateRole)
			roles.GET("/:id", authMiddleware.HasPermission("role:read"), roleHandler.GetRole)
			roles.PUT("/:id", authMiddleware.HasPermission("role:update"), roleHandler.UpdateRole)
			roles.DELETE("/:id", authMiddleware.HasPermission("role:delete"), roleHandler.DeleteRole)
			roles.POST("/:id/permissions", authMiddleware.HasPermission("role:assign"), roleHandler.AssignPermissions)
		}

		// 权限管理 - 需要认证
		permissions := api.Group("/permissions", authMiddleware.AuthRequired())
		{
			permissions.GET("", authMiddleware.HasPermission("permission:list"), permissionHandler.ListPermissions)
		}
	}

	// 启动服务器
	serverAddr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("服务器启动在 %s", serverAddr)
	if err := r.Run(serverAddr); err != nil {
		log.Fatalf("启动服务器失败: %v", err)
	}
}
