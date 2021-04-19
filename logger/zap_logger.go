// logger 基于zap日志库，进行封装的logger库
// 支持日志自动切割
package logger

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime/debug"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// defaultHostName default hostname.
var defaultHostName, _ = os.Hostname()

// tmFmtWithMS 时间格式化
const tmFmtWithMS = "2006-01-02 15:04:05.999"

// zapLogWriter zap log entry data.
type zapLogWriter struct {
	maxAge    int  // 日志保留天数
	maxSize   int  // 日志大小，单位为MB
	compress  bool // 是否日志压缩，默认日志不压缩
	addCaller bool // 是否添加行号，默认不添加文件行号

	// addCaller = true,并且 callerSkip > 0 会设置zap.AddCallerSkip
	callerSkip int

	logLevel       zapcore.Level // zap日志级别
	logWriteToFile bool          // 日志是否写入文件中
	logFilename    string        // 日志文件名，不包含路径，比如go-zap.log
	logDir         string        // 日志存放的目录
	jsonFormat     bool          // 是否json格式化
	stdout         bool          // 是否输出到终端

	// 日志是否染色
	// For example, InfoLevel is serialized to "info" and colored blue.
	enableColor bool

	// hostname host
	hostname string

	// zap底层Logger接口
	fLogger *zap.Logger
}

// New 创建一个Logger interface.
func New(opts ...Option) Logger {
	z := defaultZapLogEntry()

	z.apply(opts...)

	core, err := z.initCore()
	if err != nil {
		log.Fatalln("init zap core error: ", err)
	}

	// 当 addCaller = true 并且 callerSkip > 0 才会记录文件名和行号
	if z.addCaller && z.callerSkip > 0 {
		z.fLogger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(z.callerSkip))
	} else {
		z.fLogger = zap.New(core)
	}

	return z
}

// NewLogSugar zap log sugar语法糖
// 支持Debug,Info,Error,Panic,Warn,Fatal等方法
// 返回一个*zap.SugaredLogger
func NewLogSugar(opts ...Option) *zap.SugaredLogger {
	z := defaultZapLogEntry()

	z.apply(opts...)

	core, err := z.initCore()
	if err != nil {
		log.Fatalln("init zap core error: ", err)
	}

	// 当 addCaller = true 并且 callerSkip > 0 才会记录文件名和行号
	if z.addCaller && z.callerSkip > 0 {
		z.fLogger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(z.callerSkip))
	} else {
		z.fLogger = zap.New(core)
	}

	return z.fLogger.Sugar()
}

// defaultZapLogEntry create default zapLogWriter
func defaultZapLogEntry() *zapLogWriter {
	z := &zapLogWriter{
		maxAge:      7,
		maxSize:     512,
		compress:    false,
		logLevel:    zapcore.InfoLevel,
		logFilename: filepath.Base(os.Args[0]), // 默认程序运行时名称
		logDir:      os.TempDir(),
		stdout:      true, // 默认日志输出到stdout终端
		jsonFormat:  true,
		hostname:    defaultHostName,
	}

	return z
}

// Debug debug log.
func (z *zapLogWriter) Debug(ctx context.Context, msg string, fields ...interface{}) {
	z.fLogger.Debug(msg, z.parseFields(ctx, fields)...)
}

// Info info log.
func (z *zapLogWriter) Info(ctx context.Context, msg string, fields ...interface{}) {
	z.fLogger.Info(msg, z.parseFields(ctx, fields)...)
}

// Error error log.
func (z *zapLogWriter) Error(ctx context.Context, msg string, fields ...interface{}) {
	z.fLogger.Error(msg, z.parseFields(ctx, fields)...)
}

// Warn warn log.
func (z *zapLogWriter) Warn(ctx context.Context, msg string, fields ...interface{}) {
	z.fLogger.Warn(msg, z.parseFields(ctx, fields)...)
}

// DPanic dPanic log.
func (z *zapLogWriter) DPanic(ctx context.Context, msg string, fields ...interface{}) {
	z.fLogger.DPanic(msg, z.parseFields(ctx, fields)...)
}

// Recover 用来捕获程序运行出现的panic信息，并记录到日志中
// 这个panic信息，将采用 zap.DPanic 方法进行记录,程序继续运行，不退出
func (z *zapLogWriter) Recover(ctx context.Context, msg string, fields ...interface{}) {
	if err := recover(); err != nil {
		if len(fields) == 0 {
			fields = make([]interface{}, 0, 2)
		}

		fields = append(fields, Fullstack.String(), string(debug.Stack()))
		z.DPanic(ctx, msg, fields...)
	}
}

// Panic panic log.
func (z *zapLogWriter) Panic(ctx context.Context, msg string, fields ...interface{}) {
	z.fLogger.Panic(msg, z.parseFields(ctx, fields)...)
}

// Fatal fatal log.
func (z *zapLogWriter) Fatal(ctx context.Context, msg string, fields ...interface{}) {
	z.fLogger.Fatal(msg, z.parseFields(ctx, fields)...)
}

// parseFields 解析map[string]interface{}中的字段到zap.Field
func (z *zapLogWriter) parseFields(ctx context.Context, args []interface{}) []zap.Field {
	fLen := len(args)
	// 这里默认申请 len(args) + 20个容量，防止fields append过程中触发动态grow操作
	fields := make([]zap.Field, 0, fLen+20)
	for i := 0; i < fLen; {
		// This is a strongly-typed field. Consume it and move on.
		if f, ok := args[i].(zap.Field); ok {
			fields = append(fields, f)
			i++
			continue
		}

		// current args[i] is map
		if m, ok := args[i].(map[string]interface{}); ok {
			for k, val := range m {
				fields = append(fields, zap.Any(k, val))
			}

			i++
			continue
		}

		// Make sure this element isn't a dangling key.
		if i == fLen-1 {
			break
		}

		// Consume this value and the next, treating them as a key-value pair. If the
		// key isn't a string, add this pair to the slice of invalid pairs.
		key, val := args[i], args[i+1]
		switch v := key.(type) {
		case string:
			fields = append(fields, zap.Any(v, val))
		case int, int32, int64, float32, float64:
			fields = append(fields, zap.Any(fmt.Sprintf("%v", v), val))
		}

		i += 2
	}

	if curTime := ctx.Value(LocalTime); curTime == nil {
		// add time_local 请求本地时间字段
		fields = append(fields, zap.String(LocalTime.String(), time.Now().Format(tmFmtWithMS)))
	}

	fields = append(fields, zap.String(CurHostname.String(), z.hostname))
	// request_id 可能是一个数字，但建议请求ip使用uuid字符串
	if reqID := ctx.Value(XRequestID); reqID != nil {
		fields = append(fields, zap.Any(XRequestID.String(), reqID))
	} else {
		fields = append(fields, zap.String(XRequestID.String(), RndUUIDMd5()))
	}

	// request ip 地址存在就记录
	if ip := ctx.Value(ReqClientIP); ip != nil {
		reqIP, _ := ip.(string)
		fields = append(fields, zap.String(ReqClientIP.String(), reqIP))
	}

	// request method 请求方法
	if reqMethod := ctx.Value(RequestMethod); reqMethod != nil {
		method, _ := reqMethod.(string)
		fields = append(fields, zap.String(RequestMethod.String(), method))
	}

	// request uri 请求资源地址
	if reqURI := ctx.Value(RequestURI); reqURI != nil {
		uri, _ := reqURI.(string)
		fields = append(fields, zap.String(RequestURI.String(), uri))
	}

	return fields
}

// initCore 初始化zap core
func (z *zapLogWriter) initCore() (zapcore.Core, error) {
	// encoder config
	encoderConf := zapcore.EncoderConfig{
		TimeKey:        "time_local", // 本地时间字段
		LevelKey:       "level",
		MessageKey:     "msg",
		CallerKey:      "caller_line",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeTime:     zapcore.ISO8601TimeEncoder, // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder, // 全路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}

	// 是否染色
	if z.enableColor {
		encoderConf.EncodeLevel = zapcore.LowercaseColorLevelEncoder
	} else {
		encoderConf.EncodeLevel = zapcore.LowercaseLevelEncoder // 小写编码器
	}

	opts := make([]zapcore.WriteSyncer, 0, 2)
	if z.logWriteToFile {
		if z.logFilename == "" {
			z.logFilename = filepath.Base(os.Args[0])
		}

		if z.logDir == "" {
			z.logFilename = filepath.Join(os.TempDir(), z.logFilename) // 默认日志文件名称
		} else {
			if !z.checkPathExist(z.logDir) {
				if err := os.MkdirAll(z.logDir, 0755); err != nil {
					return nil, err
				}
			}

			z.logFilename = filepath.Join(z.logDir, z.logFilename)
		}

		// 日志最低级别设置
		syncWriter := zapcore.AddSync(&lumberjack.Logger{
			Filename:  z.logFilename, // ⽇志⽂件路径
			MaxSize:   z.maxSize,     // 单位为MB,默认为512MB
			MaxAge:    z.maxAge,      // 文件最多保存多少天
			LocalTime: true,          // 采用本地时间
			Compress:  z.compress,    // 是否压缩日志
		})

		opts = append(opts, syncWriter)
	}

	if z.stdout {
		opts = append(opts, zapcore.AddSync(os.Stdout))
	}

	// 创建一个混合WriteSyncer
	writerSyncer := zapcore.NewMultiWriteSyncer(opts...)

	// json格式化日志
	if z.jsonFormat {
		return zapcore.NewCore(zapcore.NewJSONEncoder(encoderConf), writerSyncer, z.logLevel), nil
	}

	return zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConf), writerSyncer, z.logLevel), nil
}

// checkPathExist check file or path exist
func (z *zapLogWriter) checkPathExist(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}

	if os.IsNotExist(err) {
		return false
	}

	return false
}
