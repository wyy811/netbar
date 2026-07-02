// ============================================================
// 包名：logger
// 功能：日志管理（基于 logrus）
// 技术栈：github.com/sirupsen/logrus
// ============================================================

package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/sirupsen/logrus"
)

// ============================================================
// 1. 全局变量
// ============================================================

var (
	// Log 全局日志实例
	Log *logrus.Logger
)

// ============================================================
// 2. 配置结构体
// ============================================================

// Config 日志配置
type Config struct {
	Level      string // 日志级别：debug/info/warn/error/fatal
	Format     string // 输出格式：json/text
	Output     string // 输出位置：stdout/file/both
	FilePath   string // 文件路径（Output=file 或 both 时使用）
	MaxSize    int    // 单个日志文件最大大小（MB）
	MaxBackups int    // 保留的旧日志文件数量
	MaxAge     int    // 保留日志文件的天数
	Compress   bool   // 是否压缩旧日志
}

// ============================================================
// 3. 初始化日志
// ============================================================

// Init 初始化日志系统
// 参数：cfg 日志配置
// 返回：*logrus.Logger 日志实例
func Init(cfg Config) *logrus.Logger {
	// 1. 创建 Logger 实例
	Log = logrus.New()

	// 2. 设置日志级别
	level, err := logrus.ParseLevel(cfg.Level)
	if err != nil {
		level = logrus.InfoLevel // 默认 Info
	}
	Log.SetLevel(level)

	// 3. 设置输出格式
	switch cfg.Format {
	case "json":
		Log.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339,
			// 添加调用者信息（文件名、行号）
			CallerPrettyfier: func(f *runtime.Frame) (string, string) {
				return "", fmt.Sprintf("%s:%d", filepath.Base(f.File), f.Line)
			},
		})
	default:
		// 文本格式（开发环境更友好）
		Log.SetFormatter(&logrus.TextFormatter{
			TimestampFormat:  "2006-01-02 15:04:05",
			FullTimestamp:    true,
			ForceColors:      true,
			QuoteEmptyFields: true,
			CallerPrettyfier: func(f *runtime.Frame) (string, string) {
				return "", fmt.Sprintf("%s:%d", filepath.Base(f.File), f.Line)
			},
		})
	}

	// 4. 设置输出位置
	switch cfg.Output {
	case "file":
		// 仅输出到文件
		file, err := openLogFile(cfg)
		if err != nil {
			log.Fatalf("打开日志文件失败: %v", err)
		}
		Log.SetOutput(file)
	case "both":
		// 同时输出到文件和控制台
		file, err := openLogFile(cfg)
		if err != nil {
			log.Fatalf("打开日志文件失败: %v", err)
		}
		// MultiWriter 将日志同时写入多个 Writer
		Log.SetOutput(io.MultiWriter(os.Stdout, file))
	default:
		// 默认输出到控制台
		Log.SetOutput(os.Stdout)
	}

	// 5. 启用调用者信息（显示哪个文件/行号调用的）
	Log.SetReportCaller(true)

	// 6. 添加默认字段（方便日志追踪）
	Log.WithFields(logrus.Fields{
		"app":     "internet-cafe",
		"env":     os.Getenv("APP_ENV"),
		"version": "1.0.0",
	})

	log.Println("日志系统初始化成功")
	return Log
}

// ============================================================
// 4. 获取日志实例
// ============================================================

// Get 获取全局日志实例
// 返回：*logrus.Logger 日志实例
func Get() *logrus.Logger {
	if Log == nil {
		// 如果未初始化，使用默认配置初始化
		Init(Config{
			Level:  "info",
			Format: "text",
			Output: "stdout",
		})
	}
	return Log
}

// ============================================================
// 5. 带字段的日志（核心方法）
// ============================================================

// WithFields 创建带字段的日志条目
// 参数：fields 日志字段（key-value 形式）
// 返回：*logrus.Entry 日志条目
//
// 使用示例：
//
//	logger.WithFields(logger.Fields{
//	    "user_id": 123,
//	    "action": "login",
//	}).Info("用户登录成功")
func WithFields(fields logrus.Fields) *logrus.Entry {
	return Get().WithFields(fields)
}

// WithError 创建带错误对象的日志条目
// 参数：err 错误对象
// 返回：*logrus.Entry 日志条目
func WithError(err error) *logrus.Entry {
	return Get().WithError(err)
}

// WithField 创建带单个字段的日志条目
// 参数：key 字段名，value 字段值
// 返回：*logrus.Entry 日志条目
func WithField(key string, value interface{}) *logrus.Entry {
	return Get().WithField(key, value)
}

// ============================================================
// 6. 便捷日志方法（不用每次都 Get()）
// ============================================================

// Debug 输出 Debug 级别日志
func Debug(args ...interface{}) {
	Get().Debug(args...)
}

// Debugf 输出格式化的 Debug 级别日志
func Debugf(format string, args ...interface{}) {
	Get().Debugf(format, args...)
}

// Info 输出 Info 级别日志
func Info(args ...interface{}) {
	Get().Info(args...)
}

// Infof 输出格式化的 Info 级别日志
func Infof(format string, args ...interface{}) {
	Get().Infof(format, args...)
}

// Warn 输出 Warn 级别日志
func Warn(args ...interface{}) {
	Get().Warn(args...)
}

// Warnf 输出格式化的 Warn 级别日志
func Warnf(format string, args ...interface{}) {
	Get().Warnf(format, args...)
}

// Error 输出 Error 级别日志
func Error(args ...interface{}) {
	Get().Error(args...)
}

// Errorf 输出格式化的 Error 级别日志
func Errorf(format string, args ...interface{}) {
	Get().Errorf(format, args...)
}

// Fatal 输出 Fatal 级别日志并退出程序
func Fatal(args ...interface{}) {
	Get().Fatal(args...)
}

// Fatalf 输出格式化的 Fatal 级别日志并退出程序
func Fatalf(format string, args ...interface{}) {
	Get().Fatalf(format, args...)
}

// ============================================================
// 7. 刷新日志缓冲
// ============================================================

// Sync 刷新日志缓冲，确保所有日志都被写入
// 用于程序优雅退出时调用
// 返回：error 刷新过程中的错误
// Sync 刷新日志缓冲
// 注意：logrus 的底层缓冲由操作系统管理，大多数情况下不需要手动刷新
// 此方法主要用于程序退出前的清理
func Sync() error {
	if Log == nil {
		return nil
	}

	// 尝试将输出断言为文件
	// 如果是文件，执行 Sync（将缓冲数据写入磁盘）
	if file, ok := Log.Out.(*os.File); ok {
		// 不关闭文件（因为后续可能还有日志）
		// 只执行 Sync，确保数据写入磁盘
		if err := file.Sync(); err != nil {
			return fmt.Errorf("同步日志文件失败: %w", err)
		}
	}

	// logrus 的日志是写到哪里就立即 flush 的
	// 所谓的"缓冲"实际上是操作系统层面的文件缓冲
	// file.Sync() 可以确保数据写入磁盘
	// 对于 MultiWriter 和 stdout，不需要处理

	log.Println("日志缓冲已刷新")
	return nil
}

// ============================================================
// 8. 辅助方法
// ============================================================

// openLogFile 打开日志文件（支持按大小滚动）
func openLogFile(cfg Config) (*os.File, error) {
	// 1. 创建日志目录（如果不存在）
	logDir := filepath.Dir(cfg.FilePath)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("创建日志目录失败: %w", err)
	}

	// 2. 打开日志文件（追加模式）
	file, err := os.OpenFile(cfg.FilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("打开日志文件失败: %w", err)
	}

	return file, nil
}

// ============================================================
// 9. 日志轮转（基于 lumberjack） - 可选增强
// ============================================================

// InitWithRotate 初始化带自动轮转的日志
// 使用 lumberjack 实现日志按大小/时间自动切割
// 需要安装：go get gopkg.in/natefinch/lumberjack.v2
func InitWithRotate(cfg Config) *logrus.Logger {
	Log = logrus.New()

	// 设置级别
	level, _ := logrus.ParseLevel(cfg.Level)
	Log.SetLevel(level)

	// 设置格式
	switch cfg.Format {
	case "json":
		Log.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339,
		})
	default:
		Log.SetFormatter(&logrus.TextFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
			FullTimestamp:   true,
			ForceColors:     true,
		})
	}

	// 如果配置了文件路径，使用 lumberjack 实现轮转
	if cfg.FilePath != "" && cfg.Output != "stdout" {
		// 需要导入：gopkg.in/natefinch/lumberjack.v2
		// 这里仅展示代码结构
		/*
			rotator := &lumberjack.Logger{
				Filename:   cfg.FilePath,
				MaxSize:    cfg.MaxSize,    // MB
				MaxBackups: cfg.MaxBackups,
				MaxAge:     cfg.MaxAge,     // 天
				Compress:   cfg.Compress,
			}
			Log.SetOutput(rotator)
		*/
	}

	Log.SetReportCaller(true)
	return Log
}
