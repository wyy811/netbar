package iface

import "netbar-management/internal/domain"

type SessionRepo interface {
	Create(session *domain.Session) error
	GetByID(id uint) (*domain.Session, error)
	GetAll() ([]*domain.Session, error)
	GetByMemberID(memberID int) ([]*domain.Session, error)
	GetByComputerID(computerID int) ([]*domain.Session, error)
	Update(session *domain.Session) error
	Delete(id uint) error
}
