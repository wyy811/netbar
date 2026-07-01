package iface

import "netbar-management/internal/domain"

type AuditLogRepo interface {
	Create(auditLog *domain.AuditLog) error
	GetByID(id uint) (*domain.AuditLog, error)
	GetAll() ([]*domain.AuditLog, error)
	Update(auditLog *domain.AuditLog) error
	Delete(id uint) error
}
