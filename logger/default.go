package logger

import (
	"context"
	"runtime/debug"

	"go.uber.org/zap"
)

var (
	// logEntry default logger entry.
	logEntry Logger

	// DefaultLogDir default log dir.
	DefaultLogDir = "./logs"

	// DefaultLogFile default log file.
	DefaultLogFile = "go-app.log"
)

// Default 默认zap logger对象
func Default(opts ...Option) {
	options := []Option{
		WithLogDir(DefaultLogDir),       // 日志目录
		WithLogFilename(DefaultLogFile), // 日志文件名，默认zap.log
		WithJsonFormat(true),            // json格式化
		WithCallerSkip(2),               // 如果基于这个Logger包，再包装一次，这个skip = 2,以此类推
		WithEnableColor(false),          // 日志是否染色，默认不染色
		WithLogLevel(zap.DebugLevel),    // 设置日志打印最低级别,如果不设置默认为info级别
		WithMaxAge(3),                   // 最大保存3天
		WithMaxSize(200),                // 每个日志文件最大20MB
		WithCompress(false),             // 日志不压缩
		WithStdout(false),               // 日志不输出到终端，可用opts方式改变stdout的值
		WriteToFile(true),               // 默认开启日志写入文件中
	}

	if len(opts) > 0 {
		options = append(options, opts...)
	}

	logEntry = New(options...)
}

// Debug debug级别日志
func Debug(ctx context.Context, msg string, fields ...interface{}) {
	logEntry.Debug(ctx, msg, fields...)
}

// Info info级别日志
func Info(ctx context.Context, msg string, fields ...interface{}) {
	logEntry.Info(ctx, msg, fields...)
}

// Error 错误类型的日志
func Error(ctx context.Context, msg string, fields ...interface{}) {
	logEntry.Error(ctx, msg, fields...)
}

// Warn 警告类型的日志
func Warn(ctx context.Context, msg string, fields ...interface{}) {
	logEntry.Warn(ctx, msg, fields...)
}

// DPanic 调试模式下的panic，程序不退出，继续运行
func DPanic(ctx context.Context, msg string, fields ...interface{}) {
	logEntry.DPanic(ctx, msg, fields...)
}

// Recover 用来捕获程序运行出现的panic信息，并记录到日志中
// 这个panic信息，将采用 DPanic 方法进行记录
func Recover(ctx context.Context, msg string, fields ...interface{}) {
	if err := recover(); err != nil {
		if len(fields) == 0 {
			fields = make([]interface{}, 0, 2)
		}

		fields = append(fields, Fullstack.String(), string(debug.Stack()))
		logEntry.DPanic(ctx, msg, fields...)
	}
}

// Panic 抛出panic的时候，先记录日志，然后执行panic,退出当前goroutine
// 如果没有捕获，就会退出当前程序,建议程序做defer捕获处理
func Panic(ctx context.Context, msg string, fields ...interface{}) {
	logEntry.Panic(ctx, msg, fields...)
}

// Fatal 抛出致命错误，然后退出程序
func Fatal(ctx context.Context, msg string, fields ...interface{}) {
	logEntry.Fatal(ctx, msg, fields...)
}
