// ============================================================
// 包名：database
// 功能：MySQL 数据库连接管理（基于 GORM）
// 技术栈：gorm.io/gorm + gorm.io/driver/mysql
// ============================================================

package database

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ============================================================
// 1. 全局变量
// ============================================================

var (
	// DB 全局数据库实例（私有，通过 GetDB() 访问）
	db *gorm.DB
)

// ============================================================
// 2. 配置结构体
// ============================================================

// Config 数据库配置
type Config struct {
	Host         string        // 数据库主机地址
	Port         int           // 端口
	Username     string        // 用户名
	Password     string        // 密码
	Database     string        // 数据库名
	Charset      string        // 字符集（默认 utf8mb4）
	MaxIdleConns int           // 最大空闲连接数
	MaxOpenConns int           // 最大打开连接数
	MaxLifetime  time.Duration // 连接最大存活时间
	LogLevel     string        // 日志级别（silent/error/warn/info）
}

// ============================================================
// 3. 初始化连接
// ============================================================

// InitMySQL 初始化 MySQL 连接
// 参数：cfg 数据库配置
// 返回：*gorm.DB 数据库实例，error 错误信息
func InitMySQL(cfg Config) (*gorm.DB, error) {
	// 1. 构建 DSN（Data Source Name）连接字符串
	// 格式：用户名:密码@tcp(主机:端口)/数据库名?charset=utf8mb4&parseTime=True&loc=Local
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
		cfg.Charset,
	)

	// 2. 配置 GORM 日志级别
	var logLevel logger.LogLevel
	switch cfg.LogLevel {
	case "silent":
		logLevel = logger.Silent
	case "error":
		logLevel = logger.Error
	case "warn":
		logLevel = logger.Warn
	case "info":
		logLevel = logger.Info
	default:
		logLevel = logger.Info // 默认 Info 级别
	}

	// 3. 打开数据库连接
	// gorm.Open 会建立连接池，而不是单次连接
	gormDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel), // 设置日志级别
		// 其他可选配置：
		// SkipDefaultTransaction: true,   // 跳过默认事务（提高性能）
		// PrepareStmt: true,              // 预编译 SQL（提高性能）
	})
	if err != nil {
		return nil, fmt.Errorf("连接数据库失败: %w", err)
	}

	// 4. 获取底层的 *sql.DB 对象，配置连接池参数
	sqlDB, err := gormDB.DB()
	if err != nil {
		return nil, fmt.Errorf("获取底层 sql.DB 失败: %w", err)
	}

	// 5. 设置连接池参数
	// MaxIdleConns: 空闲连接池中最大连接数，默认 2
	// 设置太小会导致频繁创建和关闭连接，设置太大会浪费资源
	if cfg.MaxIdleConns > 0 {
		sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	} else {
		sqlDB.SetMaxIdleConns(10) // 默认 10
	}

	// MaxOpenConns: 数据库最大打开连接数，默认 0（无限制）
	// 需要根据数据库的最大连接数来设置，避免超过数据库限制
	if cfg.MaxOpenConns > 0 {
		sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	} else {
		sqlDB.SetMaxOpenConns(100) // 默认 100
	}

	// ConnMaxLifetime: 连接最大存活时间，默认 0（永久）
	// MySQL 有 wait_timeout 设置（默认 8 小时），连接存活时间应小于该值
	if cfg.MaxLifetime > 0 {
		sqlDB.SetConnMaxLifetime(cfg.MaxLifetime)
	} else {
		sqlDB.SetConnMaxLifetime(time.Hour) // 默认 1 小时
	}

	// 6. 保存到全局变量
	db = gormDB

	// 7. 测试连接是否可用
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("数据库 Ping 失败: %w", err)
	}

	log.Println("MySQL 连接初始化成功")
	return gormDB, nil
}

// ============================================================
// 4. 获取数据库实例
// ============================================================

// GetDB 获取全局数据库实例
// 用于在 Repository 层获取 DB 连接
// 返回：*gorm.DB 数据库实例
func GetDB() *gorm.DB {
	if db == nil {
		log.Fatal("数据库未初始化，请先调用 InitMySQL()")
	}
	return db
}

// ============================================================
// 5. 关闭数据库连接
// ============================================================

// Close 关闭数据库连接
// 用于程序优雅退出时释放资源
// 返回：error 关闭过程中的错误
func Close() error {
	if db == nil {
		return nil
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("获取底层 sql.DB 失败: %w", err)
	}

	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("关闭数据库连接失败: %w", err)
	}

	log.Println("MySQL 连接已关闭")
	return nil
}

// ============================================================
// 6. 健康检查
// ============================================================

// Ping 检查数据库连接是否正常
// 用于健康检查接口（/health）
// 返回：error 连接异常信息
func Ping() error {
	if db == nil {
		return fmt.Errorf("数据库未初始化")
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("获取底层 sql.DB 失败: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("数据库 Ping 失败: %w", err)
	}

	return nil
}

// ============================================================
// 7. 辅助方法（可选）
// ============================================================

// IsConnected 检查数据库是否已连接
// 返回：bool true=已连接 false=未连接
func IsConnected() bool {
	return db != nil
}

// Stats 获取数据库连接池统计信息
// 用于监控和调试
// 返回：连接池统计信息
func Stats() (map[string]interface{}, error) {
	if db == nil {
		return nil, fmt.Errorf("数据库未初始化")
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	stats := sqlDB.Stats()
	return map[string]interface{}{
		"max_open_connections": stats.MaxOpenConnections,
		"open_connections":     stats.OpenConnections,
		"in_use":               stats.InUse,
		"idle":                 stats.Idle,
		"wait_count":           stats.WaitCount,
		"wait_duration":        stats.WaitDuration.String(),
		"max_idle_closed":      stats.MaxIdleClosed,
		"max_lifetime_closed":  stats.MaxLifetimeClosed,
	}, nil
}