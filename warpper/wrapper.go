package wrapper

import "log"

// Wrapper wrap goroutine to run
type Wrapper interface {
	// Wrap exec func in goroutine without recover catch
	Wrap(fn func())

	// WrapWithRecover safely execute func in goroutine
	WrapWithRecover(fn func())

	// Wait this func wait all goroutine finish
	Wait()
}

var wrapperMap = map[WrapType]constructor{
	WgWrapper: NewWgWrapper,
	ChWrapper: NewChanWrapper,
}

// New create wrapper interface
func New(wrapType WrapType, opts ...Option) Wrapper {
	if w, ok := wrapperMap[wrapType]; ok {
		return w(opts...)
	}

	panic("wrapper type not exists")
}

type constructor func(opts ...Option) Wrapper

// WrapType wrap type
type WrapType int

const (
	// WgWrapper waitGroup wrapper
	WgWrapper WrapType = iota
	// ChWrapper chan wrapper
	ChWrapper
)

// Register register wrapper
func Register(wrapType WrapType, c constructor) {
	_, ok := wrapperMap[wrapType]
	if ok {
		panic("registered wrapper already exists")
	}

	wrapperMap[wrapType] = c
}

// defaultRecovery default recover func.
func defaultRecovery() {
	if e := recover(); e != nil {
		log.Printf("wrapper exec recover:%v\n", e)
	}
}
