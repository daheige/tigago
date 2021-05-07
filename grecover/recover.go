// Package grecover catch exec panic.
package grecover

import (
	"runtime/debug"
)

// Logger log interface
type Logger interface {
	Println(args ...interface{})
}

// LoggerFunc is a bridge between Logger and any third party logger
type LoggerFunc func(msg ...interface{})

// Println implements Logger interface.
func (f LoggerFunc) Println(args ...interface{}) { f(args...) }

var (
	// dummy logger writes nothing
	dummyLogger = LoggerFunc(func(...interface{}) {})

	// LogEntry log entry
	LogEntry Logger = dummyLogger
)

// CheckPanic check panic when exit
func CheckPanic() {
	if err := recover(); err != nil {
		LogEntry.Println("panic error: ", err)
		LogEntry.Println("full stack: ", string(CatchStack()))
	}
}

// CatchStack 捕获指定stack信息,一般在处理panic/recover中处理
// 返回完整的堆栈信息和函数调用信息
func CatchStack() []byte {
	return debug.Stack()
}
