// repository/impl/user_repo_impl.go
package impl

import (
	"errors"
	"time"

	"netbar-management/internal/domain"
	iface "netbar-management/internal/repository/interface" // ✅ 导入接口包并起别名

	"gorm.io/gorm"
)

// userRepo 是 UserRepo 接口的具体实现
type userRepo struct {
	db *gorm.DB
}

// NewUserRepo 创建 UserRepo 实例（返回接口类型）
func NewUserRepo(db *gorm.DB) iface.UserRepo { // ✅ 返回接口类型
	return &userRepo{db: db}
}

// Create 创建用户
func (r *userRepo) Create(user *domain.User) error {
	result := r.db.Create(user)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("创建用户失败")
	}
	return nil
}

// GetByID 根据 ID 查询用户0
func (r *userRepo) GetByID(id uint) (*domain.User, error) {
	var user domain.User
	result := r.db.First(&user, id) // ✅ 简写方式
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // ✅ 记录不存在返回 nil, nil
		}
		return nil, result.Error
	}
	return &user, nil
}

// GetByUsername 根据用户名查询用户
func (r *userRepo) GetByUsername(username string) (*domain.User, error) {
	var user domain.User
	result := r.db.Where("username = ?", username).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &user, nil
}

// GetAll 获取所有用户（支持分页）
func (r *userRepo) GetAll() ([]*domain.User, error) {
	var users []*domain.User
	result := r.db.Order("created_at DESC").Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

// UpdateLoginTime 更新用户最后登录时间
func (r *userRepo) UpdateLoginTime(id uint) error {
	result := r.db.Model(&domain.User{}).
		Where("id = ?", id).
		Update("last_login_at", time.Now())

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("用户不存在")
	}
	return nil
}

// Update 更新用户完整信息
func (r *userRepo) Update(user *domain.User) error {
	result := r.db.Save(user)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("更新用户失败")
	}
	return nil
}

// Delete 删除用户（硬删除）
func (r *userRepo) Delete(id uint) error {
	result := r.db.Delete(&domain.User{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("用户不存在")
	}
	return nil
}
