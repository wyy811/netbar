// ============================================================
// 包名：model
// 功能：GORM 数据库模型定义
// 路径：backend/internal/model/gorm_models.go
// 说明：根据表结构文档生成，与数据库表一一对应
// ============================================================

package model

import (
	"time"
)

// ============================================================
// 1. users - 员工/管理员表
// ============================================================

// GormUser 员工/管理员表映射
type GormUser struct {
	ID           uint       `gorm:"primaryKey;autoIncrement;comment:用户ID"`
	Username     string     `gorm:"uniqueIndex:idx_users_username;size:50;not null;comment:登录名"`
	PasswordHash string     `gorm:"size:255;not null;comment:bcrypt加密密码"`
	RealName     string     `gorm:"size:50;comment:真实姓名"`
	Role         int8       `gorm:"default:1;index:idx_users_role;comment:1=普通员工,2=管理员,3=超级管理员"`
	Status       int8       `gorm:"default:1;comment:1=启用,0=禁用"`
	LastLoginAt  *time.Time `gorm:"comment:最后登录时间"`
	CreatedAt    time.Time  `gorm:"autoCreateTime;comment:创建时间"`
	UpdatedAt    time.Time  `gorm:"autoUpdateTime;comment:更新时间"`
}

// TableName 指定表名
func (GormUser) TableName() string {
	return "users"
}

// ============================================================
// 2. members - 会员表
// ============================================================

// GormMember 会员表映射
type GormMember struct {
	ID            uint      `gorm:"primaryKey;autoIncrement;comment:会员ID"`
	CardNumber    string    `gorm:"uniqueIndex:idx_members_card_number;size:20;not null;comment:会员卡号"`
	Name          string    `gorm:"size:50;not null;comment:会员姓名"`
	Phone         string    `gorm:"index:idx_members_phone;size:15;comment:手机号"`
	IDCard        string    `gorm:"size:18;comment:身份证号"`
	Balance       float64   `gorm:"type:decimal(10,2);default:0.00;comment:账户余额（元）"`
	TotalSpent    float64   `gorm:"type:decimal(10,2);default:0.00;comment:累计消费（元）"`
	DiscountLevel int8      `gorm:"default:0;comment:折扣等级 0=无,1=9折,2=8折"`
	Status        int8      `gorm:"default:1;index:idx_members_status;comment:1=正常,0=冻结"`
	RegisteredAt  time.Time `gorm:"autoCreateTime;comment:注册时间"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime;comment:更新时间"`
}

func (GormMember) TableName() string {
	return "members"
}

// ============================================================
// 3. computers - 电脑机位表
// ============================================================

// GormComputer 电脑机位表映射
type GormComputer struct {
	ID               uint       `gorm:"primaryKey;autoIncrement;comment:机位ID"`
	MachineNumber    string     `gorm:"uniqueIndex:idx_computers_machine_number;size:10;not null;comment:机位编号（如A-01）"`
	Area             string     `gorm:"index:idx_computers_area;size:20;comment:区域（电竞区/普通区/包间）"`
	IPAddress        string     `gorm:"size:15;comment:内网IP"`
	Status           int8       `gorm:"default:0;index:idx_computers_status;comment:0=空闲,1=使用中,2=维护中,3=预约中"`
	HourlyRate       float64    `gorm:"type:decimal(8,2);not null;comment:该机位每小时单价"`
	CurrentSessionID *uint      `gorm:"comment:当前上机记录ID（冗余，加速查询）"`
	LastOnlineAt     *time.Time `gorm:"comment:最后上线时间"`
	CreatedAt        time.Time  `gorm:"autoCreateTime;comment:创建时间"`
	UpdatedAt        time.Time  `gorm:"autoUpdateTime;comment:更新时间"`
}

func (GormComputer) TableName() string {
	return "computers"
}

// ============================================================
// 4. sessions - 上机记录表（核心业务表）
// ============================================================

// GormSession 上机记录表映射
type GormSession struct {
	ID              uint       `gorm:"primaryKey;autoIncrement;comment:记录ID"`
	ComputerID      uint       `gorm:"index:idx_sessions_computer_id;not null;comment:关联机位"`
	MemberID        *uint      `gorm:"index:idx_sessions_member_id;comment:会员ID（非会员为NULL）"`
	UserID          uint       `gorm:"index:idx_sessions_user_id;not null;comment:操作员工ID"`
	StartTime       time.Time  `gorm:"index:idx_sessions_start_time;not null;comment:上机开始时间"`
	EndTime         *time.Time `gorm:"comment:下机时间"`
	DurationMinutes *int       `gorm:"comment:总时长（分钟）"`
	TotalAmount     float64    `gorm:"type:decimal(10,2);comment:总费用"`
	DiscountAmount  float64    `gorm:"type:decimal(10,2);default:0.00;comment:优惠/折扣金额"`
	PaidAmount      float64    `gorm:"type:decimal(10,2);comment:实付金额"`
	PaymentMethod   *int8      `gorm:"comment:1=现金,2=会员余额,3=支付宝,4=微信"`
	Status          int8       `gorm:"default:0;index:idx_sessions_status;comment:0=进行中,1=已结束,2=异常中断"`
	Note            string     `gorm:"size:255;comment:备注"`
	CreatedAt       time.Time  `gorm:"autoCreateTime;comment:创建时间"`
	UpdatedAt       time.Time  `gorm:"autoUpdateTime;comment:更新时间"`
}

func (GormSession) TableName() string {
	return "sessions"
}

// ============================================================
// 5. rate_rules - 计费规则表
// ============================================================

// GormRateRule 计费规则表映射
type GormRateRule struct {
	ID          uint      `gorm:"primaryKey;autoIncrement;comment:规则ID"`
	RuleName    string    `gorm:"size:50;not null;comment:规则名称（如"普通时段"）"`
	DayOfWeek   *int8     `gorm:"index:idx_rate_rules_day_time;comment:1-7（星期几），NULL表示每天"`
	StartTime   string    `gorm:"size:8;not null;comment:时段开始（HH:MM:SS）"`
	EndTime     string    `gorm:"size:8;not null;comment:时段结束（HH:MM:SS）"`
	HourlyRate  float64   `gorm:"type:decimal(8,2);not null;comment:该时段单价"`
	IsOvernight int8      `gorm:"default:0;comment:是否通宵场（跨天）0=否,1=是"`
	Priority    int8      `gorm:"default:1;comment:优先级（数值越大越优先匹配）"`
	Enabled     int8      `gorm:"default:1;comment:是否启用 1=启用,0=禁用"`
	CreatedAt   time.Time `gorm:"autoCreateTime;comment:创建时间"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime;comment:更新时间"`
}

func (GormRateRule) TableName() string {
	return "rate_rules"
}

// ============================================================
// 6. recharges - 充值记录表
// ============================================================

// GormRecharge 充值记录表映射
type GormRecharge struct {
	ID            uint      `gorm:"primaryKey;autoIncrement;comment:充值ID"`
	MemberID      uint      `gorm:"index:idx_recharges_member_id;not null;comment:会员ID"`
	UserID        uint      `gorm:"index:idx_recharges_user_id;not null;comment:操作员工ID"`
	Amount        float64   `gorm:"type:decimal(10,2);not null;comment:充值金额"`
	BonusAmount   float64   `gorm:"type:decimal(10,2);default:0.00;comment:赠送金额"`
	PaymentMethod int8      `gorm:"not null;comment:1=现金,2=支付宝,3=微信"`
	BeforeBalance float64   `gorm:"type:decimal(10,2);not null;comment:充值前余额"`
	AfterBalance  float64   `gorm:"type:decimal(10,2);not null;comment:充值后余额"`
	Remark        string    `gorm:"size:255;comment:备注"`
	CreatedAt     time.Time `gorm:"autoCreateTime;index:idx_recharges_created_at;comment:创建时间"`
}

func (GormRecharge) TableName() string {
	return "recharges"
}

// ============================================================
// 7. audit_logs - 操作日志表
// ============================================================

// GormAuditLog 操作日志表映射
type GormAuditLog struct {
	ID         uint      `gorm:"primaryKey;autoIncrement;comment:日志ID"`
	UserID     uint      `gorm:"index:idx_audit_logs_user_id;not null;comment:操作人"`
	Action     string    `gorm:"index:idx_audit_logs_action;size:50;not null;comment:操作类型（LOGIN/LOGOUT/CREATE_SESSION等）"`
	TargetType string    `gorm:"size:30;comment:目标类型（computer/member/session）"`
	TargetID   *uint     `gorm:"comment:目标ID"`
	Detail     string    `gorm:"type:json;comment:详细数据（变更前后）"`
	IPAddress  string    `gorm:"size:45;comment:操作IP"`
	UserAgent  string    `gorm:"size:255;comment:客户端信息"`
	CreatedAt  time.Time `gorm:"autoCreateTime;index:idx_audit_logs_created_at;comment:创建时间"`
}

func (GormAuditLog) TableName() string {
	return "audit_logs"
}
