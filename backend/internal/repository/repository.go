package repository

import (
	repoImpl "netbar-management/internal/repository/impl"
	repoInterface "netbar-management/internal/repository/interface"

	"gorm.io/gorm"
)

type Repositories struct {
	UserRepo     repoInterface.UserRepo
	MemberRepo   repoInterface.MemberRepo
	AuditLogRepo repoInterface.AuditLogRepo
	ComputerRepo repoInterface.ComputerRepo
	RateRuleRepo repoInterface.RateRuleRepo
	RechargeRepo repoInterface.RechargeRepo
	SessionRepo  repoInterface.SessionRepo
}

// ✅ 添加构造函数
func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		UserRepo:     repoImpl.NewUserRepo(db),
		MemberRepo:   repoImpl.NewMemberRepo(db),
		AuditLogRepo: repoImpl.NewAuditLogRepo(db),
		ComputerRepo: repoImpl.NewComputerRepo(db),
		RateRuleRepo: repoImpl.NewRateRuleRepo(db),
		RechargeRepo: repoImpl.NewRechargeRepo(db),
		SessionRepo:  repoImpl.NewSessionRepo(db),
	}
}
