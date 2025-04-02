package service

import (
	"authentication/internal/config"
	"authentication/internal/model"
	"authentication/internal/repository"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"time"

	"gorm.io/gorm"
)

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// RefreshTokenRequest 刷新令牌请求
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// AuthService 认证服务接口
type AuthService interface {
	Register(req RegisterRequest) error
	Login(req LoginRequest) (*model.TokenPair, error)
	ValidateToken(token string) (*model.TokenClaims, error)
	RefreshToken(req RefreshTokenRequest) (*model.TokenPair, error)
	GetUserByID(id uint) (*model.User, error)
}

// authService 认证服务实现
type authService struct {
	userRepo  repository.UserRepository
	jwtConfig config.JWTConfig
}

// NewAuthService 创建认证服务实例
func NewAuthService(userRepo repository.UserRepository, jwtConfig config.JWTConfig) AuthService {
	return &authService{
		userRepo:  userRepo,
		jwtConfig: jwtConfig,
	}
}

// Register 用户注册
func (s *authService) Register(req RegisterRequest) error {
	// 检查用户名是否已存在
	_, err := s.userRepo.GetByUsername(req.Username)
	if err == nil {
		return errors.New("用户名已存在")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("检查用户名失败: %w", err)
	}

	// 检查邮箱是否已存在
	_, err = s.userRepo.GetByEmail(req.Email)
	if err == nil {
		return errors.New("邮箱已存在")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("检查邮箱失败: %w", err)
	}

	// 创建用户
	user := model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password, // 会在BeforeSave钩子中自动加密
		FullName: req.FullName,
		Active:   true,
	}

	// 保存用户
	if err := s.userRepo.Create(&user); err != nil {
		return fmt.Errorf("创建用户失败: %w", err)
	}

	return nil
}

// Login 用户登录
func (s *authService) Login(req LoginRequest) (*model.TokenPair, error) {
	// 获取用户
	user, err := s.userRepo.GetByUsername(req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户名或密码错误")
		}
		return nil, fmt.Errorf("获取用户失败: %w", err)
	}

	// 检查用户是否激活
	if !user.Active {
		return nil, errors.New("用户已被禁用")
	}

	// 验证密码
	if !user.CheckPassword(req.Password) {
		return nil, errors.New("用户名或密码错误")
	}

	// 获取用户权限
	var permissions []string
	for _, role := range user.Roles {
		for _, perm := range role.Permissions {
			permissions = append(permissions, perm.Code)
		}
	}

	// 生成令牌对
	tokenPair, err := s.generateTokenPair(user.ID, user.Username, permissions)
	if err != nil {
		return nil, fmt.Errorf("生成令牌失败: %w", err)
	}

	return tokenPair, nil
}

// ValidateToken 验证令牌
func (s *authService) ValidateToken(tokenString string) (*model.TokenClaims, error) {
	// 解析令牌
	claims, err := s.parseToken(tokenString)
	if err != nil {
		return nil, err
	}

	// 检查令牌类型
	if claims.TokenType != "access" {
		return nil, errors.New("无效的令牌类型")
	}

	return claims, nil
}

// RefreshToken 刷新令牌
func (s *authService) RefreshToken(req RefreshTokenRequest) (*model.TokenPair, error) {
	// 验证刷新令牌
	claims, err := s.parseToken(req.RefreshToken)
	if err != nil {
		return nil, err
	}

	// 检查令牌类型
	if claims.TokenType != "refresh" {
		return nil, errors.New("无效的令牌类型")
	}

	// 获取用户
	user, err := s.userRepo.GetByID(claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("获取用户失败: %w", err)
	}

	// 检查用户是否激活
	if !user.Active {
		return nil, errors.New("用户已被禁用")
	}

	// 获取用户权限
	var permissions []string
	for _, role := range user.Roles {
		for _, perm := range role.Permissions {
			permissions = append(permissions, perm.Code)
		}
	}

	// 生成新令牌对
	tokenPair, err := s.generateTokenPair(user.ID, user.Username, permissions)
	if err != nil {
		return nil, fmt.Errorf("生成令牌失败: %w", err)
	}

	return tokenPair, nil
}

// GetUserByID 根据ID获取用户
func (s *authService) GetUserByID(id uint) (*model.User, error) {
	return s.userRepo.GetByID(id)
}

// generateTokenPair 生成访问令牌和刷新令牌对
func (s *authService) generateTokenPair(userID uint, username string, permissions []string) (*model.TokenPair, error) {
	// 创建访问令牌
	accessTokenClaims := model.TokenClaims{
		UserID:      userID,
		Username:    username,
		Permissions: permissions,
		TokenType:   "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(s.jwtConfig.AccessExpire) * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    s.jwtConfig.Issuer,
			Subject:   fmt.Sprintf("%d", userID),
		},
	}

	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims).SignedString([]byte(s.jwtConfig.Secret))
	if err != nil {
		return nil, fmt.Errorf("生成访问令牌失败: %w", err)
	}

	// 创建刷新令牌
	refreshTokenClaims := model.TokenClaims{
		UserID:      userID,
		Username:    username,
		Permissions: permissions,
		TokenType:   "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(s.jwtConfig.RefreshExpire) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    s.jwtConfig.Issuer,
			Subject:   fmt.Sprintf("%d", userID),
		},
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims).SignedString([]byte(s.jwtConfig.Secret))
	if err != nil {
		return nil, fmt.Errorf("生成刷新令牌失败: %w", err)
	}

	return &model.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    s.jwtConfig.AccessExpire * 60, // 转换为秒
	}, nil
}

// parseToken 解析令牌
func (s *authService) parseToken(tokenString string) (*model.TokenClaims, error) {
	// 解析令牌
	token, err := jwt.ParseWithClaims(tokenString, &model.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("无效的签名方法")
		}
		return []byte(s.jwtConfig.Secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("解析令牌失败: %w", err)
	}

	if !token.Valid {
		return nil, errors.New("无效的令牌")
	}

	claims, ok := token.Claims.(*model.TokenClaims)
	if !ok {
		return nil, errors.New("无效的令牌声明")
	}

	return claims, nil
}
