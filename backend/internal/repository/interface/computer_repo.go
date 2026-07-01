package iface

import "netbar-management/internal/domain"

type ComputerRepo interface {
	Create(computer *domain.Computer) error
	GetByID(id uint) (*domain.Computer, error)
	GetAll() ([]*domain.Computer, error)
	Update(computer *domain.Computer) error
	Delete(id uint) error
}
