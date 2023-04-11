package wrapper

import (
	"sync"
)

var _ Wrapper = (*wrapWgImpl)(nil)

// wrapWgImpl sync.WaitGroup wrap impl
type wrapWgImpl struct {
	sync.WaitGroup
	recoveryFunc func()
}

// NewWgWrapper create wrapper entity
func NewWgWrapper(opts ...Option) Wrapper {
	w := &wrapWgImpl{}
	option := &Options{}
	for _, o := range opts {
		o(option)
	}

	w.recoveryFunc = option.RecoveryFunc
	if w.recoveryFunc == nil {
		w.recoveryFunc = defaultRecovery
	}

	return w
}

// Wrap fn func in goroutine to run
func (w *wrapWgImpl) Wrap(fn func()) {
	w.Add(1)
	go func() {
		defer w.Done()
		fn()
	}()
}

// WrapWithRecover exec func with recover
func (w *wrapWgImpl) WrapWithRecover(fn func()) {
	w.Add(1)
	go func() {
		defer w.recoveryFunc()
		defer w.Done()
		fn()
	}()
}
