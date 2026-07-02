// ============================================================
// 包名：config
// 功能：配置结构体定义 + 配置加载逻辑
// ============================================================

package config

import (
	"time"

	"github.com/spf13/viper"
)

// ============================================================
// 1. 顶层配置结构体
// ============================================================

// Config 应用总配置
type Config struct {
	Server         ServerConfig    `mapstructure:"server"`
	Database       DatabaseConfig  `mapstructure:"database"`
	JWT            JWTConfig       `mapstructure:"jwt"`
	Log            LogConfig       `mapstructure:"log"`
	Billing        BillingConfig   `mapstructure:"billing"`
	RateLimit      RateLimitConfig `mapstructure:"rate_limit"`
	CORS           CORSConfig      `mapstructure:"cors"`
	Initialization InitConfig      `mapstructure:"initialization"`
}

// ============================================================
// 2. 各模块配置结构体定义
// ============================================================

// ServerConfig 服务端配置
type ServerConfig struct {
	Port         int           `mapstructure:"port"`          // 服务端口
	Mode         string        `mapstructure:"mode"`          // 运行模式: debug/release/test
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`  // 读取超时
	WriteTimeout time.Duration `mapstructure:"write_timeout"` // 写入超时
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host            string        `mapstructure:"host"`              // 主机地址
	Port            int           `mapstructure:"port"`              // 端口
	Username        string        `mapstructure:"username"`          // 用户名
	Password        string        `mapstructure:"password"`          // 密码
	DBName          string        `mapstructure:"dbname"`            // 数据库名称
	Charset         string        `mapstructure:"charset"`           // 字符集
	ParseTime       bool          `mapstructure:"parse_time"`        // 是否解析时间
	Loc             string        `mapstructure:"loc"`               // 时区
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`    // 最大空闲连接数
	MaxOpenConns    int           `mapstructure:"max_open_conns"`    // 最大打开连接数
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"` // 连接最大存活时间
}

// JWTConfig JWT 认证配置
type JWTConfig struct {
	SecretKey     string        `mapstructure:"secret_key"`     // 签名密钥
	AccessExpire  time.Duration `mapstructure:"access_expire"`  // Access Token 过期时间
	RefreshExpire time.Duration `mapstructure:"refresh_expire"` // Refresh Token 过期时间
	Issuer        string        `mapstructure:"issuer"`         // 签发者
}

// LogConfig 日志配置
type LogConfig struct {
	Level  string `mapstructure:"level"`  // 日志级别: debug/info/warn/error
	Output string `mapstructure:"output"` // 输出位置: stdout/stderr/文件路径
	Format string `mapstructure:"format"` // 格式: text/json
}

// BillingConfig 计费引擎配置
type BillingConfig struct {
	MinPreDeduct    float64 `mapstructure:"min_pre_deduct"`   // 最低预扣金额
	Precision       int     `mapstructure:"precision"`        // 计费精度（小数位数）
	EnableOvernight bool    `mapstructure:"enable_overnight"` // 是否启用通宵场自动跨天
}

// RateLimitConfig 限流配置
type RateLimitConfig struct {
	RequestsPerMinute int `mapstructure:"requests_per_minute"` // 每分钟最大请求数
	Burst             int `mapstructure:"burst"`               // 突发流量额外请求数
}

// CORSConfig 跨域配置
type CORSConfig struct {
	AllowOrigins     []string `mapstructure:"allow_origins"`     // 允许的源
	AllowMethods     []string `mapstructure:"allow_methods"`     // 允许的方法
	AllowHeaders     []string `mapstructure:"allow_headers"`     // 允许的请求头
	AllowCredentials bool     `mapstructure:"allow_credentials"` // 是否允许携带凭证
	MaxAge           int      `mapstructure:"max_age"`           // 预检请求缓存时间（秒）
}

// InitConfig 系统初始化配置
type InitConfig struct {
	CreateAdmin   bool   `mapstructure:"create_admin"`    // 是否创建管理员
	AdminUsername string `mapstructure:"admin_username"`  // 管理员用户名
	AdminPassword string `mapstructure:"admin_password"`  // 管理员密码
	SeedRateRules bool   `mapstructure:"seed_rate_rules"` // 是否导入预设计费规则
}

// ============================================================
// 3. 配置加载函数
// ============================================================

// Load 加载配置文件
// 参数：configPath 配置文件路径（如 "configs/config.yaml"）
// 返回：Config 指针和错误
func Load(configPath string) (*Config, error) {
	// 设置配置文件
	viper.SetConfigFile(configPath)

	// 开启环境变量支持（可以覆盖 YAML 中的值）
	viper.AutomaticEnv()

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	// 解析到结构体
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// ============================================================
// 4. 辅助函数：获取配置默认值（可选）
// ============================================================

// DefaultConfig 返回默认配置（用于测试或快速启动）
func DefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port:         8080,
			Mode:         "debug",
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
		},
		Database: DatabaseConfig{
			Host:            "localhost",
			Port:            3306,
			Username:        "root",
			Password:        "wyy050811",
			DBName:          "internet_cafe",
			Charset:         "utf8mb4",
			ParseTime:       true,
			Loc:             "Local",
			MaxIdleConns:    10,
			MaxOpenConns:    100,
			ConnMaxLifetime: 3600 * time.Second,
		},
		JWT: JWTConfig{
			SecretKey:     "default-secret-key",
			AccessExpire:  24 * time.Hour,
			RefreshExpire: 7 * 24 * time.Hour,
			Issuer:        "internet-cafe-system",
		},
		Log: LogConfig{
			Level:  "debug",
			Output: "stdout",
			Format: "text",
		},
		Billing: BillingConfig{
			MinPreDeduct:    5.00,
			Precision:       2,
			EnableOvernight: true,
		},
		RateLimit: RateLimitConfig{
			RequestsPerMinute: 100,
			Burst:             20,
		},
		CORS: CORSConfig{
			AllowOrigins: []string{
				"http://localhost:5173",
				"http://localhost:3000",
			},
			AllowMethods: []string{
				"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS",
			},
			AllowHeaders: []string{
				"Authorization", "Content-Type", "X-Requested-With",
			},
			AllowCredentials: true,
			MaxAge:           86400,
		},
		Initialization: InitConfig{
			CreateAdmin:   true,
			AdminUsername: "admin",
			AdminPassword: "123456",
			SeedRateRules: true,
		},
	}
}
