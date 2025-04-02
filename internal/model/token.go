package model

import (
	"github.com/golang-jwt/jwt/v4"
)

// TokenClaims JWT令牌的声明
type TokenClaims struct {
	UserID      uint     `json:"user_id"`
	Username    string   `json:"username"`
	Permissions []string `json:"permissions"`
	TokenType   string   `json:"token_type"` // "access" 或 "refresh"
	jwt.RegisteredClaims
}

// TokenPair 包含访问令牌和刷新令牌
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"` // 访问令牌过期时间（秒）
}
