package repository

import "netbar-management/internal/domain"

type RateRuleRepo interface {
	Create(rateRule *domain.RateRule) error
	GetByID(id uint) (*domain.RateRule, error)
	GetAll() ([]*domain.RateRule, error)
	Update(rateRule *domain.RateRule) error
	Delete(id uint) error
}
