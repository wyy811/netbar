package impl

import (
	"errors"
	"netbar-management/internal/domain"
	iface "netbar-management/internal/repository/interface"

	"gorm.io/gorm"
)

type rechargeRepo struct {
	db *gorm.DB
}

func NewRechargeRepo(db *gorm.DB) iface.RechargeRepo {
	return &rechargeRepo{db: db}
}

func (r *rechargeRepo) Create(recharge *domain.Recharge) error {
	result := r.db.Create(recharge)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("创建充值记录失败")
	}
	return nil
}

func (r *rechargeRepo) GetByID(id uint) (*domain.Recharge, error) {
	var recharge domain.Recharge
	result := r.db.Where("id = ?", id).First(&recharge)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("充值记录不存在")
	}
	return &recharge, nil
}

func (r *rechargeRepo) GetAll() ([]*domain.Recharge, error) {
	var recharges []*domain.Recharge
	result := r.db.Find(&recharges)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("没有充值记录")
	}
	return recharges, nil
}

func (r *rechargeRepo) Update(recharge *domain.Recharge) error {
	result := r.db.Save(recharge)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("更新充值记录失败")
	}
	return nil
}

func (r *rechargeRepo) Delete(id uint) error {
	result := r.db.Delete(&domain.Recharge{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("充值记录不存在")
	}
	return nil
}
