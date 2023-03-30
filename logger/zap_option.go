package logger

import (
	"go.uber.org/zap/zapcore"
)

// Option option for zapLogWriter
type Option func(z *zapLogWriter)

// apply add option for zapLogWriter
func (z *zapLogWriter) apply(opts ...Option) {
	for _, o := range opts {
		o(z)
	}
}

// WithMaxAge 日志保留时间，单位day
func WithMaxAge(d int) Option {
	return func(z *zapLogWriter) {
		z.maxAge = d
	}
}

// WithMaxSize 日志大小，单位MB
func WithMaxSize(size int) Option {
	return func(z *zapLogWriter) {
		z.maxSize = size
	}
}

// WithCompress 日志是否压缩
func WithCompress(b bool) Option {
	return func(z *zapLogWriter) {
		z.compress = b
	}
}

// WithAddCaller 是否输出文件名和行号
func WithAddCaller(b bool) Option {
	return func(z *zapLogWriter) {
		z.addCaller = b
	}
}

/*
* WithCallerSkip 设置callerSkip
addCaller = true,并且 callerSkip > 0 会设置zap.AddCallerSkip
zap源码包中logger.go#260 check func
check must always be called directly by a method in the Logger interface
(e.g., Check, Info, Fatal).
const callerSkipOffset = 2
这里的callerSkipOffset默认是2
如果基于这个Logger包，再包装一次，这个 skip = 2,以此类推
否则 skip=1
*/
func WithCallerSkip(skip int) Option {
	return func(z *zapLogWriter) {
		z.callerSkip = skip
	}
}

// WithLogLevel 日志级别，设置日志打印最低级别，如果不设置默认为info级别
// zap.InfoLevel is the default logging priority.
func WithLogLevel(level zapcore.Level) Option {
	return func(z *zapLogWriter) {
		z.logLevel = level
	}
}

// WriteToFile 设置日志是否写入文件中
func WriteToFile(b bool) Option {
	return func(z *zapLogWriter) {
		z.logWriteToFile = b
	}
}

// WithLogFilename 日志文件名
func WithLogFilename(filename string) Option {
	return func(z *zapLogWriter) {
		z.logFilename = filename
	}
}

// WithLogDir 日志目录
func WithLogDir(dir string) Option {
	return func(z *zapLogWriter) {
		z.logDir = dir
	}
}

// WithJsonFormat 是否json格式化
func WithJsonFormat(b bool) Option {
	return func(z *zapLogWriter) {
		z.jsonFormat = b
	}
}

// WithStdout 是否输出到终端
func WithStdout(b bool) Option {
	return func(z *zapLogWriter) {
		z.stdout = b
	}
}

// WithEnableColor 是否日志染色
func WithEnableColor(b bool) Option {
	return func(z *zapLogWriter) {
		z.enableColor = b
	}
}

// WithHostname 自定义hostname
func WithHostname(hostname string) Option {
	return func(z *zapLogWriter) {
		z.hostname = hostname
	}
}
