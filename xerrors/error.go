/*
*Package xerrors
自定义错误类型，一般用在api/微服务等业务逻辑中，处理错误
支持是否输出堆栈信息，可以把stack信息记录到日志文件中，方便定位问题
*/
package xerrors

import "runtime/debug"

// ErrImpl impl error interface
type ErrImpl struct {
	msg   string
	code  int
	frame []byte // 错误堆栈信息
}

// New 创建一个error
func New(text string, code int, isStack ...bool) error {
	var b bool
	if len(isStack) > 0 && isStack[0] {
		b = true
	}

	return MakeError(text, code, b)
}

// MakeError 创建一个error
func MakeError(text string, code int, isStack bool) *ErrImpl {
	err := &ErrImpl{
		msg:  text,
		code: code,
	}

	if isStack {
		err.frame = debug.Stack()
	}

	return err
}

// Error 实现了error interface{} Error方法
func (e *ErrImpl) Error() string {
	return e.msg
}

// Code 返回code
func (e *ErrImpl) Code() int {
	return e.code
}

// Stack 打印完整的错误堆栈信息
func (e *ErrImpl) Stack() []byte {
	return e.frame
}
