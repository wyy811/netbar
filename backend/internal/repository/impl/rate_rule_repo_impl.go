package impl

import (
	"errors"
	"netbar-management/internal/domain"
	iface "netbar-management/internal/repository/interface"

	"gorm.io/gorm"
)

type rateRuleRepo struct {
	db *gorm.DB
}

func NewRateRuleRepo(db *gorm.DB) iface.RateRuleRepo {
	return &rateRuleRepo{db: db}
}

func (r *rateRuleRepo) Create(rateRule *domain.RateRule) error {
	result := r.db.Create(rateRule)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("创建费率规则失败")
	}
	return nil
}

func (r *rateRuleRepo) GetByID(id uint) (*domain.RateRule, error) {
	var rateRule domain.RateRule
	result := r.db.Where("id = ?", id).First(&rateRule)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("费率规则不存在")
	}
	return &rateRule, nil
}

func (r *rateRuleRepo) GetAll() ([]*domain.RateRule, error) {
	var rateRules []*domain.RateRule
	result := r.db.Find(&rateRules)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("没有费率规则")
	}
	return rateRules, nil
}

func (r *rateRuleRepo) Update(rateRule *domain.RateRule) error {
	result := r.db.Save(rateRule)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("更新费率规则失败")
	}
	return nil
}

func (r *rateRuleRepo) Delete(id uint) error {
	result := r.db.Delete(&domain.RateRule{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("删除费率规则失败")
	}
	return nil
}
