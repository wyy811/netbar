package iface

import "netbar-management/internal/domain"

type UserRepo interface {
	Create(user *domain.User) error
	GetByID(id uint) (*domain.User, error)
	GetByUsername(username string) (*domain.User, error)
	GetAll() ([]*domain.User, error)
	UpdateLoginTime(id uint) error
	Update(user *domain.User) error
	Delete(id uint) error
}
