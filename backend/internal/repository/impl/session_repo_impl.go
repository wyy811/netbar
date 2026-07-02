package impl

import (
	"errors"
	"netbar-management/internal/domain"
	iface "netbar-management/internal/repository/interface"

	"gorm.io/gorm"
)

type sessionRepo struct {
	db *gorm.DB
}

func NewSessionRepo(db *gorm.DB) iface.SessionRepo {
	return &sessionRepo{db: db}
}

func (r *sessionRepo) Create(session *domain.Session) error {
	result := r.db.Create(session)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("创建上下机记录失败")
	}
	return nil
}

func (r *sessionRepo) GetByID(id uint) (*domain.Session, error) {
	var session domain.Session
	result := r.db.Where("id = ?", id).First(&session)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("上下机记录不存在")
	}
	return &session, nil
}

func (r *sessionRepo) GetAll() ([]*domain.Session, error) {
	var sessions []*domain.Session
	result := r.db.Find(&sessions)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("没有上下机记录")
	}
	return sessions, nil
}

func (r *sessionRepo) GetByMemberID(memberID int) ([]*domain.Session, error) {
	var sessions []*domain.Session
	result := r.db.Where("member_id = ?", memberID).Find(&sessions)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("没有上下机记录")
	}
	return sessions, nil
}

func (r *sessionRepo) GetByComputerID(computerID int) ([]*domain.Session, error) {
	var sessions []*domain.Session
	result := r.db.Where("computer_id = ?", computerID).Find(&sessions)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("没有上下机记录")
	}
	return sessions, nil
}

func (r *sessionRepo) Update(session *domain.Session) error {
	result := r.db.Save(session)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("更新上下机记录失败")
	}
	return nil
}

func (r *sessionRepo) Delete(id uint) error {
	result := r.db.Delete(&domain.Session{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("上下机记录不存在")
	}
	return nil
}
