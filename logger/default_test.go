package logger

import (
	"context"
	"testing"

	"go.uber.org/zap"
)

// TestDefaultLogEntry test default log entry.
func TestDefaultLogEntry(t *testing.T) {
	opts := []Option{
		WithLogDir("./logs"),           // 日志目录
		WithLogFilename("zap-web.log"), // 日志文件名，默认zap.log
		WithStdout(true),               // 一般生产环境，建议不输出到stdout
		WithJsonFormat(true),           // json格式化
		WithAddCaller(true),            // 打印行号
		WithCallerSkip(2),              // 如果基于这个Logger包，再包装一次，这个skip = 2,以此类推
		WithEnableColor(false),         // 日志是否染色，默认不染色
		WithLogLevel(zap.DebugLevel),   // 设置日志打印最低级别,如果不设置默认为info级别
		WithMaxAge(3),                  // 最大保存3天
		WithMaxSize(20),                // 每个日志文件最大20MB
		WithCompress(false),            // 日志不压缩
		WithStdout(false),              // 关闭终端输出
		// WithHostname("myapp.com"),      // 设置hostname
	}

	// 生成默认的日志句柄对象
	Default(opts...)

	// 模拟请求id
	reqId := RndUUIDMd5()
	ctx := context.Background()
	ctx = context.WithValue(ctx, XRequestID, reqId)

	Debug(ctx, "hello daheige", map[string]interface{}{
		"a": 1,
	})

	Info(ctx, "hello", map[string]interface{}{
		"a": 1,
		"b": 2,
		"c": "world",
	})

	Error(ctx, "exec error", zap.Any("details", map[string]interface{}{
		"name": "zap",
		"age":  30,
	}))

	Warn(ctx, "run warning", "key", 1234, "abc", "hello world")

	// DPanic调试信息
	DPanic(ctx, "panic debug", "abc", 23456, "message", "hello world")
}
