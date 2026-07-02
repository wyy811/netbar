// ============================================================
// 包名：service
// 功能：认证服务（登录/登出/刷新Token/获取用户信息）
// 依赖：UserRepo（用户数据访问）、JWT Manager（令牌管理）
// ============================================================

package service

import (
	"errors"

	"netbar-management/internal/domain"
	"netbar-management/internal/pkg/crypto"
	"netbar-management/internal/pkg/jwt"
	iface "netbar-management/internal/repository/interface"
)

// ============================================================
// 1. 结构体定义
// ============================================================

// AuthService 认证服务
type AuthService struct {
	userRepo   iface.UserRepo
	jwtManager *jwt.Manager
}

// NewAuthService 创建认证服务实例（依赖注入）
func NewAuthService(userRepo iface.UserRepo, jwtManager *jwt.Manager) *AuthService {
	return &AuthService{
		userRepo:   userRepo,
		jwtManager: jwtManager,
	}
}

// ============================================================
// 2. 请求/响应 DTO
// ============================================================

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	AccessToken  string       `json:"accessToken"`
	RefreshToken string       `json:"refreshToken"`
	User         *domain.User `json:"user"`
}

// RefreshTokenRequest 刷新令牌请求
type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

// RefreshTokenResponse 刷新令牌响应
type RefreshTokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

// LogoutRequest 登出请求
type LogoutRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required,min=6"`
}

// ============================================================
// 3. 核心方法：登录
// ============================================================

// Login 用户登录
// 流程：查询用户 → 验证状态 → 验证密码 → 更新登录时间 → 生成Token
func (s *AuthService) Login(req LoginRequest) (*LoginResponse, error) {
	// 1. 通过用户名查询用户
	user, err := s.userRepo.GetByUsername(req.Username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("用户名或密码错误")
	}

	// 2. 检查用户状态（是否被禁用）
	if user.Status == 0 {
		return nil, errors.New("账号已被禁用，请联系管理员")
	}

	// 3. 验证密码（user.Password 存储的是加密后的哈希值）
	if !crypto.CheckPasswordHash(req.Password, user.Password) {
		return nil, errors.New("用户名或密码错误")
	}

	// 4. 更新最后登录时间（异步执行，不影响主流程）
	_ = s.userRepo.UpdateLoginTime(user.ID)

	// 5. 生成 Access Token 和 Refresh Token
	accessToken, err := s.jwtManager.GenerateAccessToken(user.ID, user.Username, user.Role)
	if err != nil {
		return nil, errors.New("生成访问令牌失败")
	}

	refreshToken, err := s.jwtManager.GenerateRefreshToken(user.ID, user.Username, user.Role)
	if err != nil {
		return nil, errors.New("生成刷新令牌失败")
	}

	// 6. 返回响应（user.Password 因为有 json:"-"，不会被序列化返回给前端）
	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
}

// ============================================================
// 4. 核心方法：刷新 Token
// ============================================================

// RefreshToken 刷新 Access Token
// 使用 Refresh Token 换取新的 Access Token（支持 Token 轮换）
func (s *AuthService) RefreshToken(req RefreshTokenRequest) (*RefreshTokenResponse, error) {
	// 调用 JWT Manager 刷新令牌
	// 内部会验证 Refresh Token 有效性，并生成新的 Access Token 和 Refresh Token
	newAccess, newRefresh, err := s.jwtManager.RefreshTokens(req.RefreshToken)
	if err != nil {
		return nil, err
	}

	return &RefreshTokenResponse{
		AccessToken:  newAccess,
		RefreshToken: newRefresh,
	}, nil
}

// ============================================================
// 5. 核心方法：登出
// ============================================================

// Logout 用户登出
// 将 Access Token 和 Refresh Token 加入黑名单
func (s *AuthService) Logout(accessToken string, req LogoutRequest) error {
	// 将两个 Token 加入黑名单
	return s.jwtManager.Logout(accessToken, req.RefreshToken)
}

// ============================================================
// 6. 辅助方法：获取用户信息
// ============================================================

// GetUserByID 根据 ID 获取用户信息
func (s *AuthService) GetUserByID(id uint) (*domain.User, error) {
	return s.userRepo.GetByID(id)
}

// GetUserByUsername 根据用户名获取用户信息
func (s *AuthService) GetUserByUsername(username string) (*domain.User, error) {
	return s.userRepo.GetByUsername(username)
}

// GetAllUsers 获取所有用户列表
func (s *AuthService) GetAllUsers() ([]*domain.User, error) {
	return s.userRepo.GetAll()
}

// ============================================================
// 7. 安全相关方法
// ============================================================

// ChangePassword 修改密码
func (s *AuthService) ChangePassword(userID uint, req ChangePasswordRequest) error {
	// 1. 获取用户信息
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("用户不存在")
	}

	// 2. 验证旧密码
	if !crypto.CheckPasswordHash(req.OldPassword, user.Password) {
		return errors.New("原密码错误")
	}

	// 3. 加密新密码
	hashedPassword, err := crypto.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	// 4. 更新密码
	user.Password = hashedPassword
	return s.userRepo.Update(user)
}

// ResetPassword 重置密码（管理员操作）
func (s *AuthService) ResetPassword(userID uint, newPassword string) error {
	// 1. 检查用户是否存在
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("用户不存在")
	}

	// 2. 加密新密码
	hashedPassword, err := crypto.HashPassword(newPassword)
	if err != nil {
		return err
	}

	// 3. 更新密码
	user.Password = hashedPassword
	return s.userRepo.Update(user)
}

// UpdateUserStatus 更新用户状态（启用/禁用）
func (s *AuthService) UpdateUserStatus(userID uint, status int) error {
	// 1. 检查用户是否存在
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("用户不存在")
	}

	// 2. 更新状态
	user.Status = status
	return s.userRepo.Update(user)
}

// ============================================================
// 8. 权限检查辅助方法
// ============================================================

// HasPermission 检查用户是否有指定权限
// role: 1=员工, 2=管理员, 3=超级管理员
// minRole: 所需的最低权限
func (s *AuthService) HasPermission(userID uint, minRole int) (bool, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return false, err
	}
	if user == nil {
		return false, errors.New("用户不存在")
	}

	return user.Role >= minRole, nil
}

// IsAdmin 检查用户是否是管理员
func (s *AuthService) IsAdmin(userID uint) (bool, error) {
	return s.HasPermission(userID, 2)
}

// IsSuperAdmin 检查用户是否是超级管理员
func (s *AuthService) IsSuperAdmin(userID uint) (bool, error) {
	return s.HasPermission(userID, 3)
}

// ============================================================
// 9. Token 验证（给中间件使用）
// ============================================================

// ValidateToken 验证 Token 是否有效
func (s *AuthService) ValidateToken(tokenString string) (*jwt.CustomClaims, error) {
	return s.jwtManager.ParseToken(tokenString)
}

// GetUserFromToken 从 Token 中提取用户信息
func (s *AuthService) GetUserFromToken(tokenString string) (*domain.User, error) {
	// 1. 解析 Token
	claims, err := s.jwtManager.ParseToken(tokenString)
	if err != nil {
		return nil, err
	}

	// 2. 通过用户 ID 获取完整用户信息
	return s.userRepo.GetByID(claims.UserID)
}
