// Package logger for log interface
package logger

import (
	"context"
)

// Logger interface.
type Logger interface {
	// Debug debug级别日志
	Debug(ctx context.Context, msg string, fields ...interface{})

	// Info info级别日志
	Info(ctx context.Context, msg string, fields ...interface{})

	// Error 错误类型的日志
	Error(ctx context.Context, msg string, fields ...interface{})

	// Warn 警告类型的日志
	Warn(ctx context.Context, msg string, fields ...interface{})

	// DPanic 调试模式下的panic，程序不退出，继续运行
	DPanic(ctx context.Context, msg string, fields ...interface{})

	// Recover 用来捕获程序运行出现的panic信息，并记录到日志中
	// 这个panic信息，将采用 DPanic 方法进行记录
	Recover(ctx context.Context, msg string, fields ...interface{})

	// Panic 抛出panic的时候，先记录日志，然后执行panic,退出当前goroutine
	// 如果没有捕获，就会退出当前程序,建议程序做defer捕获处理
	Panic(ctx context.Context, msg string, fields ...interface{})

	// Fatal 抛出致命错误，然后退出程序
	Fatal(ctx context.Context, msg string, fields ...interface{})
}
