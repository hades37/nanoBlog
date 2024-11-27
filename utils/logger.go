package utils

import (
	"fmt"
	"sync"
	"time"
)

// 定义日志级别
const (
	DEBUG = iota
	INFO
	WARN
	ERROR
	FATAL
)

// 日志级别对应的文字
var levelText = map[int]string{
	DEBUG: "DEBUG",
	INFO:  "INFO",
	WARN:  "WARN",
	ERROR: "ERROR",
	FATAL: "FATAL",
}

// 日志级别对应的颜色代码
var levelColor = map[int]string{
	DEBUG: "\033[1;37m", // 白色
	INFO:  "\033[1;32m", // 绿色
	WARN:  "\033[1;33m", // 黄色
	ERROR: "\033[1;31m", // 红色
	FATAL: "\033[1;35m", // 紫色
}

// Logger 日志结构体
type Logger struct {
	level int
}

var (
	instance *Logger
	once     sync.Once
)

// GetLogger 获取日志实例（单例）
func GetLogger() *Logger {
	once.Do(func() {
		instance = &Logger{
			level: DEBUG, // 默认日志级别
		}
	})
	return instance
}

// SetLevel 设置日志级别
func (l *Logger) SetLevel(level int) {
	if level >= DEBUG && level <= FATAL {
		l.level = level
	}
}

// log 通用日志输出方法
func (l *Logger) log(level int, format string, args ...interface{}) {
	// 判断是否需要输出该级别的日志
	if level < l.level {
		return
	}

	// 获取当前时间
	now := time.Now().Format("2006-01-02 15:04:05")

	// 构造日志前缀
	prefix := fmt.Sprintf("%s%s [%s]%s",
		levelColor[level],
		now,
		levelText[level],
		"\033[0m", // 重置颜色
	)

	// 格式化日志内容
	var content string
	if len(args) > 0 {
		content = fmt.Sprintf(format, args...)
	} else {
		content = format
	}

	// 输出日志
	fmt.Printf("%s %s\n", prefix, content)
}

// Debug 输出调试日志
func (l *Logger) Debug(format string, args ...interface{}) {
	l.log(DEBUG, format, args...)
}

// Info 输出信息日志
func (l *Logger) Info(format string, args ...interface{}) {
	l.log(INFO, format, args...)
}

// Warn 输出警告日志
func (l *Logger) Warn(format string, args ...interface{}) {
	l.log(WARN, format, args...)
}

// Error 输出错误日志
func (l *Logger) Error(format string, args ...interface{}) {
	l.log(ERROR, format, args...)
}

// Fatal 输出致命错误日志
func (l *Logger) Fatal(format string, args ...interface{}) {
	l.log(FATAL, format, args...)
}

// 包级别的快捷方法
func Debug(format string, args ...interface{}) {
	GetLogger().Debug(format, args...)
}

func Info(format string, args ...interface{}) {
	GetLogger().Info(format, args...)
}

func Warn(format string, args ...interface{}) {
	GetLogger().Warn(format, args...)
}

func Error(format string, args ...interface{}) {
	GetLogger().Error(format, args...)
}

func Fatal(format string, args ...interface{}) {
	GetLogger().Fatal(format, args...)
}

// SetLogLevel 设置全局日志级别
func SetLogLevel(level int) {
	GetLogger().SetLevel(level)
}
