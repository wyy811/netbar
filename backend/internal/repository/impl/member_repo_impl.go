package impl

import (
	"errors"
	"netbar-management/internal/domain"
	iface "netbar-management/internal/repository/interface" // ✅ 导入接口包并起别名

	"gorm.io/gorm"
)

type memberRepo struct {
	db *gorm.DB
}

// NewMemberRepo 创建 MemberRepo 实例(返回接口类型)
func NewMemberRepo(db *gorm.DB) iface.MemberRepo {
	return &memberRepo{db: db}
}

// Create 创建用户
func (r *memberRepo) Create(member *domain.Member) error {
	result := r.db.Create(member)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("创建用户失败")
	}
	return nil
}

// GetByID 根据ID查询用户
func (r *memberRepo) GetByID(id uint) (*domain.Member, error) {
	var member domain.Member
	result := r.db.First(&member, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &member, nil
}

// GetByName 根据用户名查询用户
func (r *memberRepo) GetByName(name string) (*domain.Member, error) {
	var member domain.Member
	result := r.db.Where("username=?", name).First(&member)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &member, nil
}

func (r *memberRepo) GetByCardID(cardID string) (*domain.Member, error) {
	var member domain.Member
	result := r.db.Where("cardid=?", cardID).First(&member)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &member, nil
}

// GetAll 获取所有用户（支持分页）
func (r *memberRepo) GetAll() ([]*domain.Member, error) {
	var members []*domain.Member
	result := r.db.Order("created_at DESC").Find(&members)
	if result.Error != nil {
		return nil, result.Error
	}
	return members, nil
}

// / Update 更新用户完整信息
func (r *memberRepo) Update(member *domain.Member) error {
	result := r.db.Save(member)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("更新用户失败")
	}
	return nil
}

// Delete 删除用户（硬删除）
func (r *memberRepo) Delete(id uint) error {
	result := r.db.Delete(&domain.Member{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("用户不存在")
	}
	return nil
}
