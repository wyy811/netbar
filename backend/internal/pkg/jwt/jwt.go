// ============================================================
// 包名：jwt
// 功能：JWT（JSON Web Token）的生成、解析、刷新和黑名单管理
// 技术栈：github.com/golang-jwt/jwt/v5
// ============================================================

package jwt

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// ============================================================
// 1. 常量定义
// ============================================================

const (
	// TokenTypeAccess 访问令牌类型
	TokenTypeAccess = "access"
	// TokenTypeRefresh 刷新令牌类型
	TokenTypeRefresh = "refresh"
)

// ============================================================
// 2. 自定义 Claims（载荷结构）
// ============================================================

// CustomClaims 自定义 JWT 载荷
// 嵌入 jwt.RegisteredClaims 包含标准字段（exp, iat, iss 等）
type CustomClaims struct {
	UserID    uint   `json:"user_id"`    // 用户ID
	Username  string `json:"username"`   // 用户名
	Role      int    `json:"role"`       // 角色：1=普通员工, 2=管理员, 3=超级管理员
	TokenType string `json:"token_type"` // 令牌类型：access / refresh

	jwt.RegisteredClaims
}

// ============================================================
// 3. JWT 管理器结构体
// ============================================================

// Manager JWT 管理器
type Manager struct {
	secretKey     []byte        // 签名密钥
	accessExpire  time.Duration // Access Token 过期时间
	refreshExpire time.Duration // Refresh Token 过期时间
	issuer        string        // 签发者
	blacklist     *Blacklist    // 黑名单（内存实现，生产环境建议用 Redis）
}

// Config JWT 配置
type Config struct {
	SecretKey     string        `mapstructure:"secret_key"`
	AccessExpire  time.Duration `mapstructure:"access_expire"`
	RefreshExpire time.Duration `mapstructure:"refresh_expire"`
	Issuer        string        `mapstructure:"issuer"`
}

// ============================================================
// 4. 黑名单实现（内存版）
// ============================================================

// Blacklist JWT 黑名单（用于登出/注销场景）
type Blacklist struct {
	mu    sync.RWMutex
	token map[string]int64 // token -> 过期时间戳（用于清理）
}

// NewBlacklist 创建黑名单实例
func NewBlacklist() *Blacklist {
	return &Blacklist{
		token: make(map[string]int64),
	}
}

// Add 将 Token 加入黑名单
func (b *Blacklist) Add(token string, expireAt time.Time) {
	b.mu.Lock()
	defer b.mu.Unlock()
	// 定期清理过期 Token（每 100 条触发一次清理）
	if len(b.token)%100 == 0 {
		b.cleanup()
	}
	b.token[token] = expireAt.Unix()
}

// Contains 检查 Token 是否在黑名单中
func (b *Blacklist) Contains(token string) bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	expireAt, exists := b.token[token]
	if !exists {
		return false
	}
	// 如果 Token 已过期，自动移除
	if expireAt < time.Now().Unix() {
		// 延迟删除，避免在锁内执行删除操作
		go b.remove(token)
		return false
	}
	return true
}

// remove 移除 Token（内部使用）
func (b *Blacklist) remove(token string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	delete(b.token, token)
}

// cleanup 清理过期的 Token
func (b *Blacklist) cleanup() {
	now := time.Now().Unix()
	for token, expireAt := range b.token {
		if expireAt < now {
			delete(b.token, token)
		}
	}
}

// ============================================================
// 5. 构造函数
// ============================================================

// NewManager 创建 JWT 管理器实例
func NewManager(cfg Config) *Manager {
	return &Manager{
		secretKey:     []byte(cfg.SecretKey),
		accessExpire:  cfg.AccessExpire,
		refreshExpire: cfg.RefreshExpire,
		issuer:        cfg.Issuer,
		blacklist:     NewBlacklist(),
	}
}

// ============================================================
// 6. 核心方法：生成 Token
// ============================================================

// GenerateAccessToken 生成 Access Token
func (m *Manager) GenerateAccessToken(userID uint, username string, role int) (string, error) {
	now := time.Now()
	claims := CustomClaims{
		UserID:    userID,
		Username:  username,
		Role:      role,
		TokenType: TokenTypeAccess,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.accessExpire)),
			NotBefore: jwt.NewNumericDate(now),
			ID:        fmt.Sprintf("%d-%d", userID, now.UnixNano()), // 唯一ID
		},
	}
	return m.generateToken(claims)
}

// GenerateRefreshToken 生成 Refresh Token
func (m *Manager) GenerateRefreshToken(userID uint, username string, role int) (string, error) {
	now := time.Now()
	claims := CustomClaims{
		UserID:    userID,
		Username:  username,
		Role:      role,
		TokenType: TokenTypeRefresh,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.refreshExpire)),
			NotBefore: jwt.NewNumericDate(now),
			ID:        fmt.Sprintf("refresh-%d-%d", userID, now.UnixNano()),
		},
	}
	return m.generateToken(claims)
}

// generateToken 生成 Token（内部方法）
func (m *Manager) generateToken(claims CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secretKey)
}

// ============================================================
// 7. 核心方法：解析和验证 Token
// ============================================================

// ParseToken 解析并验证 Token
// 返回 CustomClaims 和错误
func (m *Manager) ParseToken(tokenString string) (*CustomClaims, error) {
	// 1. 检查黑名单
	if m.blacklist.Contains(tokenString) {
		return nil, errors.New("token has been blacklisted")
	}

	// 2. 解析 Token
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return m.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	// 3. 验证 Token 有效性
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	// 4. 提取 Claims
	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return nil, errors.New("invalid claims type")
	}

	return claims, nil
}

// ValidateToken 仅验证 Token 是否有效（不返回 Claims）
// 用于中间件中快速检查
func (m *Manager) ValidateToken(tokenString string) error {
	_, err := m.ParseToken(tokenString)
	return err
}

// ============================================================
// 8. 核心方法：刷新 Token
// ============================================================

// RefreshTokens 使用 Refresh Token 刷新 Access Token
// 返回新的 Access Token 和新的 Refresh Token（可选）
func (m *Manager) RefreshTokens(refreshToken string) (newAccessToken string, newRefreshToken string, err error) {
	// 1. 解析 Refresh Token
	claims, err := m.ParseToken(refreshToken)
	if err != nil {
		return "", "", fmt.Errorf("invalid refresh token: %w", err)
	}

	// 2. 验证 Token 类型
	if claims.TokenType != TokenTypeRefresh {
		return "", "", errors.New("invalid token type: expected refresh token")
	}

	// 3. 将旧的 Refresh Token 加入黑名单（一次性使用）
	// 获取原 Token 的过期时间
	expireAt := claims.ExpiresAt.Time
	m.blacklist.Add(refreshToken, expireAt)

	// 4. 生成新的 Access Token 和 Refresh Token（可选）
	newAccess, err := m.GenerateAccessToken(claims.UserID, claims.Username, claims.Role)
	if err != nil {
		return "", "", err
	}

	// 5. 生成新的 Refresh Token（实现 Refresh Token 轮换）
	newRefresh, err := m.GenerateRefreshToken(claims.UserID, claims.Username, claims.Role)
	if err != nil {
		return "", "", err
	}

	return newAccess, newRefresh, nil
}

// ============================================================
// 9. 辅助方法
// ============================================================

// GetTokenExpireTime 获取 Token 的过期时间（不验证签名）
func (m *Manager) GetTokenExpireTime(tokenString string) (time.Time, error) {
	parser := jwt.NewParser()
	token, _, err := parser.ParseUnverified(tokenString, &CustomClaims{})
	if err != nil {
		return time.Time{}, err
	}
	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return time.Time{}, errors.New("invalid claims")
	}
	return claims.ExpiresAt.Time, nil
}

// Logout 登出：将 Access Token 和 Refresh Token 加入黑名单
func (m *Manager) Logout(accessToken, refreshToken string) error {
	// 解析 Access Token 获取过期时间
	accessClaims, err := m.ParseToken(accessToken)
	if err != nil {
		// 即使解析失败，也尝试加入黑名单（防止恶意攻击）
		expireAt := time.Now().Add(time.Hour) // 保守估计
		m.blacklist.Add(accessToken, expireAt)
	} else {
		m.blacklist.Add(accessToken, accessClaims.ExpiresAt.Time)
	}

	// 解析 Refresh Token 获取过期时间
	refreshClaims, err := m.ParseToken(refreshToken)
	if err == nil {
		m.blacklist.Add(refreshToken, refreshClaims.ExpiresAt.Time)
	} else {
		expireAt := time.Now().Add(7 * 24 * time.Hour)
		m.blacklist.Add(refreshToken, expireAt)
	}

	return nil
}

// ============================================================
// 10. 错误类型定义（便于上层区分处理）
// ============================================================

var (
	// ErrTokenExpired Token 过期
	ErrTokenExpired = errors.New("token expired")
	// ErrTokenInvalid Token 无效
	ErrTokenInvalid = errors.New("token invalid")
	// ErrTokenBlacklisted Token 已被拉黑
	ErrTokenBlacklisted = errors.New("token blacklisted")
	// ErrTokenTypeInvalid Token 类型错误
	ErrTokenTypeInvalid = errors.New("token type invalid")
)
