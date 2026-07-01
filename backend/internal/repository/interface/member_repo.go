package iface

import "netbar-management/internal/domain"

type MemberRepo interface {
	Create(member *domain.Member) error
	GetByID(id uint) (*domain.Member, error)
	GetByName(name string) (*domain.Member, error)
	GetByCardID(cardID string) (*domain.Member, error)
	GetAll() ([]*domain.Member, error)
	Update(member *domain.Member) error
	Delete(id uint) error
}
