package impl

import (
	"errors"
	"netbar-management/internal/domain"
	iface "netbar-management/internal/repository/interface"

	"gorm.io/gorm"
)

type auditLogRepo struct {
	db *gorm.DB
}

func NewAuditLogRepo(db *gorm.DB) iface.AuditLogRepo {
	return &auditLogRepo{db: db}
}

func (r *auditLogRepo) Create(auditLog *domain.AuditLog) error {
	result := r.db.Create(auditLog)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("创建审计日志失败")
	}
	return nil
}

func (r *auditLogRepo) GetByID(id uint) (*domain.AuditLog, error) {
	var auditLog domain.AuditLog
	result := r.db.Where("id = ?", id).First(&auditLog)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("审计日志不存在")
	}
	return &auditLog, nil
}

func (r *auditLogRepo) GetAll() ([]*domain.AuditLog, error) {
	var auditLogs []*domain.AuditLog
	result := r.db.Find(&auditLogs)
	if result.Error != nil {
		return nil, result.Error
	}
	return auditLogs, nil
}

func (r *auditLogRepo) Update(auditLog *domain.AuditLog) error {
	result := r.db.Save(auditLog)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("更新审计日志失败")
	}
	return nil
}

func (r *auditLogRepo) Delete(id uint) error {
	result := r.db.Delete(&domain.AuditLog{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("此审计日志不存在")
	}
	return nil
}
