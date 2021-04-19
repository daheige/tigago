package main

import (
	"context"
	"log"
	"time"

	"go.uber.org/zap"

	"github.com/daheige/tigago/logger"
)

// LogEntry 日志接口对象
var LogEntry logger.Logger

func main() {
	// 对于log option 下面的可以根据实际情况使用
	LogEntry = logger.New(
		logger.WithLogDir("./logs"),         // 日志目录
		logger.WithLogFilename("zap.log"),   // 日志文件名，默认zap.log
		logger.WithStdout(true),             // 一般生产环境，建议不输出到stdout
		logger.WithJsonFormat(true),         // json格式化
		logger.WithAddCaller(true),          // 打印行号
		logger.WithCallerSkip(1),            // 如果基于这个Logger包，再包装一次，这个skip = 2,以此类推
		logger.WithEnableColor(false),       // 日志是否染色，默认不染色
		logger.WithLogLevel(zap.DebugLevel), // 设置日志打印最低级别,如果不设置默认为info级别
		logger.WithMaxAge(3),                // 最大保存3天
		logger.WithMaxSize(20),              // 每个日志文件最大20MB
		logger.WithCompress(false),          // 日志不压缩
		// logger.WithHostname("myapp.com"),    // 设置hostname
		logger.WriteToFile(true), // 日志写入文件中
		logger.WithStdout(false), // 关闭日志写入终端
	)

	// 模拟请求id
	reqId := logger.RndUUIDMd5()
	ctx := context.Background()
	ctx = context.WithValue(ctx, logger.XRequestID, reqId)
	LogEntry.Info(ctx, "hello", map[string]interface{}{
		"a": 1,
		"b": 12,
	})

	LogEntry.Error(ctx, "exec error", zap.Any("details", map[string]interface{}{
		"name": "zap",
		"age":  30,
	}))

	LogEntry.Warn(ctx, "run warning", "key", 1234, "abc", "hello world")

	// DPanic调试信息
	LogEntry.DPanic(ctx, "panic debug", "abc", 23456, "message", "hello world")

	// 在协程中抛出了panic，然后用Recover进行捕获panic信息
	go func() {
		defer LogEntry.Recover(ctx, "exec panic")

		log.Println("hello world")

		// 主动抛出panic
		panic("abc")
	}()

	log.Println(123)

	time.Sleep(1 * time.Second)
	log.Println("server will exit")
}
