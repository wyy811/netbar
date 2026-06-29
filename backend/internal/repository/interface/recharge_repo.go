package repository

import "netbar-management/internal/domain"

type RechargeRepo interface {
	Create(recharge *domain.Recharge) error
	GetByID(id uint) (*domain.Recharge, error)
	GetAll() ([]*domain.Recharge, error)
	Update(recharge *domain.Recharge) error
	Delete(id uint) error
}
