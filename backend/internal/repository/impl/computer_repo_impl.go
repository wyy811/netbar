package impl

import (
	"errors"
	"netbar-management/internal/domain"
	iface "netbar-management/internal/repository/interface"

	"gorm.io/gorm"
)

type computerRepo struct {
	db *gorm.DB
}

func NewComputerRepo(db *gorm.DB) iface.ComputerRepo {
	return &computerRepo{db: db}
}

func (r *computerRepo) Create(computer *domain.Computer) error {
	result := r.db.Create(computer)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("创建电脑信息失败")
	}
	return nil
}

func (r *computerRepo) GetByID(id uint) (*domain.Computer, error) {
	var computer domain.Computer
	result := r.db.First(&computer, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &computer, nil
}

func (r *computerRepo) GetAll() ([]*domain.Computer, error) {
	var computers []*domain.Computer
	result := r.db.Order("created_at DESC").Find(&computers)
	if result.Error != nil {
		return nil, result.Error
	}
	return computers, nil
}

func (r *computerRepo) Update(computer *domain.Computer) error {
	result := r.db.Save(computer)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("更新电脑信息失败")
	}
	return nil
}

// Delete 删除用户（硬删除）
func (r *computerRepo) Delete(id uint) error {
	result := r.db.Delete(&domain.Computer{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("此电脑不存在")
	}
	return nil
}
