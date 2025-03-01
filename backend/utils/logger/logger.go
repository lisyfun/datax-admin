package logger

import (
	"datax-admin/config"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

var (
	// 不同级别的日志记录器
	InfoLogger  *log.Logger
	WarnLogger  *log.Logger
	ErrorLogger *log.Logger
	DebugLogger *log.Logger
)

// Config 日志配置
type Config struct {
	LogPath    string // 日志文件路径
	MaxSize    int    // 单个日志文件最大大小（MB）
	MaxBackups int    // 最大保留的旧日志文件数
	MaxAge     int    // 日志文件保留的最大天数
	Compress   bool   // 是否压缩旧日志文件
}

// Init 初始化日志系统
func Init() error {
	logConfig := config.GlobalConfig.Logger

	// 确保日志目录存在
	if err := os.MkdirAll(logConfig.LogPath, 0755); err != nil {
		return fmt.Errorf("创建日志目录失败: %v", err)
	}

	// 创建日志文件
	logFile, err := os.OpenFile(filepath.Join(logConfig.LogPath, "app.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("打开日志文件失败: %v", err)
	}

	// 同时输出到文件和控制台
	multiWriter := io.MultiWriter(os.Stdout, logFile)

	// 初始化不同级别的日志记录器，不设置前缀，在消息中包含级别
	InfoLogger = log.New(multiWriter, "", 0)
	WarnLogger = log.New(multiWriter, "", 0)
	ErrorLogger = log.New(multiWriter, "", 0)
	DebugLogger = log.New(multiWriter, "", 0)

	return nil
}

// formatLog 格式化日志消息
func formatLog(level, msg string) string {
	// 获取调用者的文件名和行号
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "unknown"
		line = 0
	}
	// 只保留文件的最后一部分
	file = filepath.Base(file)

	// 格式化日志消息，保持与 Gin 日志格式一致
	return fmt.Sprintf("%s %s | %+5s | %s:%d | %s ",
		"[GIN]",
		time.Now().Format("2006/01/02 - 15:04:05"),
		level,
		file,
		line,
		msg,
	)
}

// Info 记录信息级别的日志
func Info(format string, v ...any) {
	msg := fmt.Sprintf(format, v...)
	InfoLogger.Println(formatLog("INFO", msg))
}

// Warn 记录警告级别的日志
func Warn(format string, v ...any) {
	msg := fmt.Sprintf(format, v...)
	WarnLogger.Println(formatLog("WARN", msg))
}

// Error 记录错误级别的日志
func Error(format string, v ...any) {
	msg := fmt.Sprintf(format, v...)
	ErrorLogger.Println(formatLog("ERROR", msg))
}

// Debug 记录调试级别的日志
func Debug(format string, v ...any) {
	msg := fmt.Sprintf(format, v...)
	DebugLogger.Println(formatLog("DEBUG", msg))
}

// Fatal 记录致命错误并退出程序
func Fatal(format string, v ...any) {
	msg := fmt.Sprintf(format, v...)
	ErrorLogger.Println(formatLog("FATAL", msg))
	os.Exit(1)
}
