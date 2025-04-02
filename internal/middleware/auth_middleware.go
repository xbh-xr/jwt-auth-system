package middleware

import (
	"authentication/internal/config"
	"authentication/internal/service"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware 认证中间件
type AuthMiddleware struct {
	jwtConfig config.JWTConfig
}

// NewAuthMiddleware 创建认证中间件实例
func NewAuthMiddleware(jwtConfig config.JWTConfig) *AuthMiddleware {
	return &AuthMiddleware{
		jwtConfig: jwtConfig,
	}
}

// AuthRequired 需要认证的中间件
func (m *AuthMiddleware) AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取令牌
		tokenString, err := extractTokenFromHeader(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		// 验证令牌
		authService := c.MustGet("authService").(service.AuthService)
		claims, err := authService.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的令牌"})
			c.Abort()
			return
		}

		// 将用户信息存储在上下文中
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("permissions", claims.Permissions)

		c.Next()
	}
}

// HasPermission 检查是否有指定权限的中间件
func (m *AuthMiddleware) HasPermission(permissionCode string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取用户权限
		permissions, exists := c.Get("permissions")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "未找到权限信息"})
			c.Abort()
			return
		}

		// 检查是否有指定权限
		hasPermission := false
		for _, p := range permissions.([]string) {
			if p == permissionCode {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{"error": "没有权限执行此操作"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// HasRole 检查是否有指定角色的中间件
func (m *AuthMiddleware) HasRole(roleName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取用户ID
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "未找到用户信息"})
			c.Abort()
			return
		}

		// 获取用户
		authService := c.MustGet("authService").(service.AuthService)
		user, err := authService.GetUserByID(userID.(uint))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户信息失败"})
			c.Abort()
			return
		}

		// 检查是否有指定角色
		if !user.HasRole(roleName) {
			c.JSON(http.StatusForbidden, gin.H{"error": "没有权限执行此操作"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// extractTokenFromHeader 从请求头提取令牌
func extractTokenFromHeader(c *gin.Context) (string, error) {
	auth := c.GetHeader("Authorization")
	if auth == "" {
		return "", errors.New("未提供认证信息")
	}

	parts := strings.SplitN(auth, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		return "", errors.New("认证格式无效")
	}

	return parts[1], nil
}
